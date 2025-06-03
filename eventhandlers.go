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
	case func(*Session, *IncomingMessage):
		return messageCreateEventHandler(v)
	}

	return nil
}

// messageCreateEventHandler is an event handler for MessageCreate events.
type messageCreateEventHandler func(*Session, *IncomingMessage)

// Type returns the event type for MessageCreate events.
func (eh messageCreateEventHandler) Type() string {
	return messageReceiveEventType
}

// New returns a new instance of MessageCreate.
func (eh messageCreateEventHandler) New() interface{} {
	return &IncomingMessage{}
}

// Handle is the handler for MessageCreate events.
func (eh messageCreateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*IncomingMessage); ok {
		eh(s, t)
	}
}

func init() {
	registerInterfaceProvider(messageCreateEventHandler(nil))
}
