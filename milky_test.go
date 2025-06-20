package Milky_go_sdk

import (
	"fmt"
	"os"
	"testing"
)

type TestLogger struct {
}

func (l *TestLogger) Infof(format string, args ...interface{}) {
	fmt.Printf("INFO: "+format+"\n", args...)
}

func (l *TestLogger) Errorf(format string, args ...interface{}) {
	fmt.Printf("ERROR: "+format+"\n", args...)
}

func (l *TestLogger) Debugf(format string, args ...interface{}) {
	fmt.Printf("DEBUG: "+format+"\n", args...)
}

func (l *TestLogger) Warnf(format string, args ...interface{}) {
	fmt.Printf("WARN: "+format+"\n", args...)
}

func (l *TestLogger) Info(args ...interface{}) {
	fmt.Println("INFO:", fmt.Sprint(args...))
}

func (l *TestLogger) Error(args ...interface{}) {
	fmt.Println("ERROR:", fmt.Sprint(args...))
}

func (l *TestLogger) Debug(args ...interface{}) {
	fmt.Println("DEBUG:", fmt.Sprint(args...))
}

func (l *TestLogger) Warn(args ...interface{}) {
	fmt.Println("WARN:", fmt.Sprint(args...))
}

func TestMilky(m *testing.T) {
	logger := &TestLogger{}
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
			fmt.Printf("Received friend message: SenderName %s, MessageSeq %d, FriendCategoryID %d, FriendCategoryName %s ", m.Friend.Nickname, m.MessageSeq, m.Friend.Category.CategoryID, m.Friend.Category.CategoryName)
		}
		fmt.Printf("Received message: Sender %d", m.SenderId)
		if m.Segments != nil {
			for _, segment := range m.Segments {
				fmt.Printf(" Segment: %v", segment)
			}
		}
		fmt.Println()
	})
	session.AddHandler(func(session2 *Session, m *FriendRequest) {
		if m == nil {
			return
		}
		fmt.Printf("Received friend request: UserId %d, Comment %s, State: %s\n", m.InitiatorID, m.Comment, m.State)

	})
	err = session.Open()
	if err != nil {
		m.Fatalf("Failed to open session: %v", err)
	}
	info, err := session.GetLoginInfo()
	if err != nil {
		m.Fatalf("Failed to get login info: %v", err)
	}
	logger.Infof("Login info: UserId %d, Nickname %s", info.UIN, info.Nickname)
	text := TextElement{
		Text: "Hello, this is a test message from MilkyGo SDK!",
	}
	var elements []IMessageElement
	elements = append(elements, &text)
	//// get from env
	//targetUserID := os.Getenv("TEST_TARGET_USER_ID")
	//targetUserIDInt, err := strconv.ParseInt(targetUserID, 10, 64)
	//if err != nil {
	//	m.Fatalf("Invalid TEST_TARGET_USER_ID: %v", err)
	//}
	//message, err := session.SendPrivateMessage(targetUserIDInt, &elements)
	//if err != nil {
	//	m.Fatalf("Failed to send private message: %v", err)
	//}
	//_ = logger.Log(LevelInfo, "Sent message: MessageId %d, Time %d", message.MessageSeq, message.Time)
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
