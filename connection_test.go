package Milky_go_sdk

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

type quietLogger struct{}

func (quietLogger) Infof(string, ...interface{})  {}
func (quietLogger) Errorf(string, ...interface{}) {}
func (quietLogger) Debugf(string, ...interface{}) {}
func (quietLogger) Warnf(string, ...interface{})  {}
func (quietLogger) Info(...interface{})           {}
func (quietLogger) Error(...interface{})          {}
func (quietLogger) Debug(...interface{})          {}
func (quietLogger) Warn(...interface{})           {}

func TestConnectionNotificationsAndShutdown(t *testing.T) {
	peers := make(chan *websocket.Conn, 8)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := (&websocket.Upgrader{}).Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		peers <- conn
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}))
	defer server.Close()
	session, _ := New("ws"+strings.TrimPrefix(server.URL, "http"), server.URL, "", quietLogger{})
	defer session.Shutdown()
	changes := make(chan bool, 16)
	session.OnConnectionChange = func(s *Session, connected bool) {
		_ = s.IsConnected() // 通知不得持有会话锁。
		changes <- connected
	}
	refused := atomic.Bool{}
	attempts := make(chan struct{}, 16)
	dialer := *websocket.DefaultDialer
	dialer.NetDialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		if refused.Load() {
			select {
			case attempts <- struct{}{}:
			default:
			}
			return nil, &net.OpError{Op: "dial", Net: network, Err: syscall.ECONNREFUSED}
		}
		return (&net.Dialer{}).DialContext(ctx, network, addr)
	}
	session.Dialer = &dialer
	expectChange := func(want bool) {
		t.Helper()
		select {
		case got := <-changes:
			if got != want {
				t.Fatalf("connected = %v, want %v", got, want)
			}
		case <-time.After(6 * time.Second):
			t.Fatal("connection notification timed out")
		}
	}
	if err := session.Open(); err != nil {
		t.Fatal(err)
	}
	expectChange(true)
	peer := <-peers
	refused.Store(true)
	_ = peer.Close()
	expectChange(false)
	select {
	case <-attempts:
	case <-time.After(6 * time.Second):
		t.Fatal("no reconnect attempt")
	}
	if session.IsConnected() {
		t.Fatal("refused connection still reports connected")
	}
	refused.Store(false)
	expectChange(true)
	if !session.IsConnected() {
		t.Fatal("reconnection did not restore connected state")
	}
	peer = <-peers
	refused.Store(true)
	_ = peer.Close()
	expectChange(false)
	started := time.Now()
	if err := session.Shutdown(); err != nil {
		t.Fatal(err)
	}
	if time.Since(started) > time.Second {
		t.Fatal("shutdown waited for reconnect backoff")
	}
	if err := session.Open(); err == nil {
		t.Fatal("shut down session reopened")
	}
}
