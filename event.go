package Milky_go_sdk

// EventHandler is an interface for events.
type EventHandler interface {
	// Type returns the type of event this handler belongs to.
	Type() string

	// Handle is called whenever an event of Type() happens.
	// It is the receivers responsibility to type assert that the interface
	// is the expected struct.
	Handle(*Session, interface{})
}

// EventInterfaceProvider is an interface for providing empty interfaces for events.
type EventInterfaceProvider interface {
	// Type is the type of event this handler belongs to.
	Type() string

	// New returns a new instance of the struct this event handler handles.
	// This is called once per event.
	// The struct is provided to all handlers of the same Type().
	New() interface{}
}

// eventHandlerInstance is a wrapper around an event handler, as functions
// cannot be compared directly.
type eventHandlerInstance struct {
	eventHandler EventHandler
}

// interfaceEventType is the event handler type for interface{} events.
const interfaceEventType = "__INTERFACE__"

// interfaceEventHandler is an event handler for interface{} events.
type interfaceEventHandler func(*Session, interface{})

// Type returns the event type for interface{} events.
func (eh interfaceEventHandler) Type() string {
	return interfaceEventType
}

// Handle is the handler for an interface{} event.
func (eh interfaceEventHandler) Handle(s *Session, i interface{}) {
	eh(s, i)
}

var registeredInterfaceProviders = map[string]EventInterfaceProvider{}

// registerInterfaceProvider registers a provider so that MilkyGo can
// access it's New() method.
func registerInterfaceProvider(eh EventInterfaceProvider) {
	if _, ok := registeredInterfaceProviders[eh.Type()]; ok {
		return
		// XXX:
		// if we should error here, we need to do something with it.
		// fmt.Errorf("event %s already registered", eh.Type())
	}
	registeredInterfaceProviders[eh.Type()] = eh
	return
}

// addEventHandler adds an event handler that will be fired anytime
// the WSAPI matching eventHandler.Type() fires.
func (s *Session) addEventHandler(eventHandler EventHandler) func() {
	s.handlersMu.Lock()
	defer s.handlersMu.Unlock()

	if s.handlers == nil {
		s.handlers = map[string][]*eventHandlerInstance{}
	}

	ehi := &eventHandlerInstance{eventHandler}
	s.handlers[eventHandler.Type()] = append(s.handlers[eventHandler.Type()], ehi)

	return func() {
		s.removeEventHandlerInstance(eventHandler.Type(), ehi)
	}
}

// addEventHandler adds an event handler that will be fired the next time
// the WSAPI matching eventHandler.Type() fires.
func (s *Session) addEventHandlerOnce(eventHandler EventHandler) func() {
	s.handlersMu.Lock()
	defer s.handlersMu.Unlock()

	if s.onceHandlers == nil {
		s.onceHandlers = map[string][]*eventHandlerInstance{}
	}

	ehi := &eventHandlerInstance{eventHandler}
	s.onceHandlers[eventHandler.Type()] = append(s.onceHandlers[eventHandler.Type()], ehi)

	return func() {
		s.removeEventHandlerInstance(eventHandler.Type(), ehi)
	}
}

// AddHandler allows you to add an event handler that will be fired anytime
// The return value of this method is a function, that when called will remove the
// event handler.
func (s *Session) AddHandler(handler interface{}) func() {
	eh := handlerForInterface(handler)

	if eh == nil {
		s.Logger.Errorf("Invalid handler type, handler will never be called")
		return func() {}
	}

	return s.addEventHandler(eh)
}

// AddHandlerOnce allows you to add an event handler that will be fired the next time
// the WSAPI event that matches the function fires.
// See AddHandler for more details.
func (s *Session) AddHandlerOnce(handler interface{}) func() {
	eh := handlerForInterface(handler)

	if eh == nil {
		s.Logger.Errorf("Invalid handler type, handler will never be called")
		return func() {}
	}

	return s.addEventHandlerOnce(eh)
}

// removeEventHandler instance removes an event handler instance.
func (s *Session) removeEventHandlerInstance(t string, ehi *eventHandlerInstance) {
	s.handlersMu.Lock()
	defer s.handlersMu.Unlock()

	handlers := s.handlers[t]
	for i := range handlers {
		if handlers[i] == ehi {
			s.handlers[t] = append(handlers[:i], handlers[i+1:]...)
		}
	}

	onceHandlers := s.onceHandlers[t]
	for i := range onceHandlers {
		if onceHandlers[i] == ehi {
			s.onceHandlers[t] = append(onceHandlers[:i], onceHandlers[i+1:]...)
		}
	}
}

// Handles calling permanent and once handlers for an event type.
func (s *Session) handle(t string, i interface{}) {
	for _, eh := range s.handlers[t] {
		if s.SyncEvents {
			eh.eventHandler.Handle(s, i)
		} else {
			go eh.eventHandler.Handle(s, i)
		}
	}

	if len(s.onceHandlers[t]) > 0 {
		for _, eh := range s.onceHandlers[t] {
			if s.SyncEvents {
				eh.eventHandler.Handle(s, i)
			} else {
				go eh.eventHandler.Handle(s, i)
			}
		}
		s.onceHandlers[t] = nil
	}
}

// Handles an event type by calling internal methods, firing handlers and firing the
// interface{} event.
func (s *Session) handleEvent(t string, i interface{}) {
	s.handlersMu.RLock()
	defer s.handlersMu.RUnlock()

	// All events are dispatched internally first.
	// s.onInterface(i)

	// Then they are dispatched to anyone handling interface{} events.
	s.handle(interfaceEventType, i)

	// Finally they are dispatched to any typed handlers.
	s.handle(t, i)
}
