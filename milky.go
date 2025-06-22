package Milky_go_sdk

import (
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

const version = "0.3.0"

func New(wsGateway string, restGateway string, token string, logger Logger) (s *Session, err error) {
	// Create an empty Session interface.
	s = &Session{
		ShouldReconnectOnError: true,
		MaxRestRetries:         3,
		Client:                 &http.Client{Timeout: 20 * time.Second},
		Dialer:                 websocket.DefaultDialer,
		UserAgent:              "MilkyGo (" + wsGateway + ", v" + version + ")",
		LastHeartbeatAck:       time.Now().UTC(),
		WSGateway:              wsGateway,
		RestGateway:            restGateway,
		Logger:                 logger,
		Token:                  token,
	}

	return
}
