package Milky_go_sdk

// See events.go
// Following are all the event types.
// Event type values are used to match the events returned
const (
	messageReceiveEventType = "message_receive"
	friendRequestEventType  = "friend_request"
	botOfflineEventType     = "bot_offline"
)

func handlerForInterface(handler interface{}) EventHandler {
	switch v := handler.(type) {
	case func(*Session, interface{}):
		return interfaceEventHandler(v)
	case func(*Session, *ReceiveMessage):
		return messageReceiveEventHandler(v)
	case func(*Session, *FriendRequest):
		return friendRequestEventHandler(v)
	case func(*Session, *BotOffline):
		return botOfflineEventHandler(v)
	}

	return nil
}

// messageReceiveEventHandler is an event handler for ReceiveMessage events.
type messageReceiveEventHandler func(*Session, *ReceiveMessage)

// Type returns the event type for MessageCreate events.
func (eh messageReceiveEventHandler) Type() string {
	return messageReceiveEventType
}

// New returns a new instance of MessageCreate.
func (eh messageReceiveEventHandler) New() interface{} {
	return &ReceiveMessage{}
}

// Handle is the handler for MessageCreate events.
func (eh messageReceiveEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*ReceiveMessage); ok {
		eh(s, t)
	}
}

type friendRequestEventHandler func(*Session, *FriendRequest)

// Type returns the event type for FriendRequest events.
func (eh friendRequestEventHandler) Type() string {
	return friendRequestEventType
}

// New returns a new instance of FriendRequest.
func (eh friendRequestEventHandler) New() interface{} {
	return &FriendRequest{}
}

// Handle is the handler for FriendRequest events.
func (eh friendRequestEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*FriendRequest); ok {
		eh(s, t)
	}
}

type botOfflineEventHandler func(*Session, *BotOffline)

func (eh botOfflineEventHandler) Handle(session *Session, i interface{}) {
	if t, ok := i.(*BotOffline); ok {
		eh(session, t)
	}
}

func (eh botOfflineEventHandler) Type() string {
	return botOfflineEventType
}

func (eh botOfflineEventHandler) New() interface{} {
	return &BotOffline{}
}

func init() {
	registerInterfaceProvider(messageReceiveEventHandler(nil))
	registerInterfaceProvider(friendRequestEventHandler(nil))
	registerInterfaceProvider(botOfflineEventHandler(nil))
}
