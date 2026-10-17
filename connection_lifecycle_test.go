package Milky_go_sdk

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func newConnectionTestSession(t *testing.T) *Session {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := (&websocket.Upgrader{}).Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}))
	t.Cleanup(server.Close)
	session, err := New("ws"+strings.TrimPrefix(server.URL, "http"), server.URL, "", quietLogger{})
	if err != nil {
		t.Fatal(err)
	}
	session.ShouldReconnectOnError = false
	t.Cleanup(func() { _ = session.Shutdown() })
	return session
}

func TestConnectionShutdownCancelsDial(t *testing.T) {
	session, err := New("ws://127.0.0.1", "", "", quietLogger{})
	if err != nil {
		t.Fatal(err)
	}
	started := make(chan struct{})
	release := make(chan struct{})
	openDone := make(chan error, 1)
	shutdownDone := make(chan error, 1)
	dialer := *websocket.DefaultDialer
	dialer.NetDialContext = func(ctx context.Context, _, _ string) (net.Conn, error) {
		close(started)
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-release:
			return nil, errors.New("test dial released")
		}
	}
	session.Dialer = &dialer
	go func() { openDone <- session.Open() }()
	select {
	case <-started:
	case <-time.After(5 * time.Second):
		close(release)
		t.Fatal("dial did not start")
	}
	go func() { shutdownDone <- session.Shutdown() }()
	defer close(release)
	select {
	case err := <-shutdownDone:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("shutdown did not cancel the pending dial")
	}
	if err := <-openDone; err == nil {
		t.Fatal("canceled dial succeeded")
	}
}

func TestConnectionStaleHeartbeatKeepsNewConnection(t *testing.T) {
	session := newConnectionTestSession(t)
	if err := session.Open(); err != nil {
		t.Fatal(err)
	}
	session.RLock()
	previous := session.wsConn
	listening := session.listening
	session.RUnlock()
	if err := session.Close(); err != nil {
		t.Fatal(err)
	}
	if err := session.Open(); err != nil {
		t.Fatal(err)
	}
	// An old heartbeat may already be running when Close and Open replace its socket.
	session.heartbeat(previous, listening, time.Hour)
	if !session.IsConnected() {
		t.Fatal("an old heartbeat closed the replacement connection")
	}
}

func TestConnectionShutdownCancelsHandshake(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(started)
		<-release
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	t.Cleanup(server.Close)
	defer close(release)
	session, err := New("ws"+strings.TrimPrefix(server.URL, "http"), server.URL, "", quietLogger{})
	if err != nil {
		t.Fatal(err)
	}
	openDone := make(chan error, 1)
	shutdownDone := make(chan error, 1)
	go func() { openDone <- session.Open() }()
	select {
	case <-started:
	case <-time.After(5 * time.Second):
		t.Fatal("WebSocket handshake did not start")
	}
	go func() { shutdownDone <- session.Shutdown() }()
	select {
	case err := <-shutdownDone:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("shutdown did not cancel the pending WebSocket handshake")
	}
	if err := <-openDone; err == nil {
		t.Fatal("canceled handshake succeeded")
	}
}

func TestConnectionCloseAndShutdownNotifications(t *testing.T) {
	session := newConnectionTestSession(t)
	changes := make(chan bool, 8)
	session.OnConnectionChange = func(s *Session, connected bool) {
		_ = s.IsConnected()
		changes <- connected
	}
	expectChange := func(want bool) {
		t.Helper()
		select {
		case got := <-changes:
			if got != want {
				t.Fatalf("connection change = %v, want %v", got, want)
			}
		case <-time.After(5 * time.Second):
			t.Fatal("connection notification timed out")
		}
	}
	if err := session.Open(); err != nil {
		t.Fatal(err)
	}
	expectChange(true)
	if err := session.Close(); err != nil {
		t.Fatal(err)
	}
	expectChange(false)
	if err := session.Close(); err != nil {
		t.Fatal(err)
	}
	if err := session.Open(); err != nil {
		t.Fatal(err)
	}
	expectChange(true)
	if err := session.Shutdown(); err != nil {
		t.Fatal(err)
	}
	expectChange(false)
	if err := session.Shutdown(); err != nil {
		t.Fatal(err)
	}
	if err := session.Open(); err == nil {
		t.Fatal("shut down session reopened")
	}
	if session.IsConnected() {
		t.Fatal("shut down session reports connected")
	}
	select {
	case got := <-changes:
		t.Fatalf("duplicate connection notification: %v", got)
	default:
	}
}
