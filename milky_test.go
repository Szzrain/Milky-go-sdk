package Milky_go_sdk

import (
	"fmt"
	"testing"
)

type TestLogger struct {
	Level Level
}

func (l *TestLogger) Log(level Level, keyvals ...interface{}) error {
	if len(keyvals) > 0 {
		if msg, ok := keyvals[0].(string); ok {
			keyvals = keyvals[1:]
			fmt.Println(fmt.Sprintf(msg, keyvals...))
		}
	}
	return nil
}

func TestMilky(m *testing.T) {
	logger := &TestLogger{Level: LevelDebug}
	session, err := New("ws://127.0.0.1:3039/event", logger)
	if err != nil {
		m.Fatalf("Failed to create session: %v", err)
	}
	session.AddHandler(func(session2 *Session, m *IncomingMessage) {
		if m == nil {
			return
		}
		fmt.Printf("Received message: Sender %d, Content: %v\n", m.SenderId, m.Segments)
	})
	err = session.Open()
	if err != nil {
		m.Fatalf("Failed to open session: %v", err)
	}
	select {}
}
