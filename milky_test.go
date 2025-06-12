package Milky_go_sdk

import (
	"fmt"
	"os"
	"strconv"
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
	fmt.Println("Gateway WS:", os.Getenv("TEST_WS_GATEWAY"))
	session, err := New(os.Getenv("TEST_WS_GATEWAY"), os.Getenv("TEST_REST_GATEWAY"), logger)
	if err != nil {
		m.Fatalf("Failed to create session: %v", err)
	}
	session.AddHandler(func(session2 *Session, m *ReceiveMessage) {
		if m == nil {
			return
		}
		if m.MessageScene == "group" {
			fmt.Printf("Received group message: GroupId %d, MessageSeq %d, ", m.Group.GroupId, m.MessageSeq)
		} else if m.MessageScene == "friend" {
			fmt.Printf("Received friend message: SenderName %s, MessageSeq %d, ", m.Friend.Nickname, m.MessageSeq)
		}
		fmt.Printf("Received message: Sender %d", m.SenderId)
		if m.Segments != nil {
			for _, segment := range m.Segments {
				fmt.Printf(" Segment: %v", segment)
			}
		}
		fmt.Println()
	})
	err = session.Open()
	if err != nil {
		m.Fatalf("Failed to open session: %v", err)
	}
	info, err := session.GetLoginInfo()
	if err != nil {
		m.Fatalf("Failed to get login info: %v", err)
	}
	_ = logger.Log(LevelInfo, "Login info: UserId %d, Nickname %s", info.UIN, info.Nickname)
	text := TextElement{
		Text: "Hello, this is a test message from MilkyGo SDK!",
	}
	var elements []IMessageElement
	elements = append(elements, &text)
	//// get from env
	targetUserID := os.Getenv("TEST_TARGET_USER_ID")
	targetUserIDInt, err := strconv.ParseInt(targetUserID, 10, 64)
	if err != nil {
		m.Fatalf("Invalid TEST_TARGET_USER_ID: %v", err)
	}
	message, err := session.SendPrivateMessage(targetUserIDInt, &elements)
	if err != nil {
		m.Fatalf("Failed to send private message: %v", err)
	}
	_ = logger.Log(LevelInfo, "Sent message: MessageId %d, Time %d", message.MessageSeq, message.Time)
	//targetGroupID := os.Getenv("TEST_TARGET_GROUP_ID")
	//targetGroupIDInt, err := strconv.ParseInt(targetGroupID, 10, 64)
	//if err != nil {
	//	m.Fatalf("Invalid TEST_TARGET_GROUP_ID: %v", err)
	//}
	//message, err := session.SendGroupMessage(targetGroupIDInt, &elements)
	//if err != nil {
	//	m.Fatalf("Failed to send group message: %v", err)
	//}
	//_ = logger.Log(LevelInfo, "Sent message: MessageId %d, Time %d", message.MessageSeq, message.Time)
	select {}
}
