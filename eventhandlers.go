package Milky_go_sdk

// See events.go
// Following are all the event types.
// Event type values are used to match the events returned
const (
	messageReceiveEventType = "message_receive"
)

func handlerForInterface(handler interface{}) EventHandler {
	switch v := handler.(type) {
	case func(*Session, interface{}):
		return interfaceEventHandler(v)
	case func(*Session, *ReceiveMessage):
		return messageReceiveEventHandler(v)
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

func init() {
	registerInterfaceProvider(messageReceiveEventHandler(nil))
}
