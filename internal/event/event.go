package event

import (
	"sync"
)

// Event represents an event that can be dispatched and handled
type Event struct {
	Name string
	Data map[string]interface{}
}

// Handler is a function that handles an event
type Handler func(event Event)

// EventDispatcher manages event subscriptions and dispatching
type EventDispatcher struct {
	handlers map[string][]Handler
	mu       sync.RWMutex
}

// NewEventDispatcher creates a new event dispatcher
func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handlers: make(map[string][]Handler),
	}
}

// Subscribe registers a handler for an event
func (d *EventDispatcher) Subscribe(eventName string, handler Handler) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, exists := d.handlers[eventName]; !exists {
		d.handlers[eventName] = make([]Handler, 0)
	}

	d.handlers[eventName] = append(d.handlers[eventName], handler)
}

// Dispatch sends an event to all registered handlers
func (d *EventDispatcher) Dispatch(event Event) {
	d.mu.RLock()
	handlers, exists := d.handlers[event.Name]
	d.mu.RUnlock()

	if !exists {
		return
	}

	// Call all handlers
	for _, handler := range handlers {
		go handler(event)
	}
}

// DispatchSync sends an event to all registered handlers synchronously
func (d *EventDispatcher) DispatchSync(event Event) {
	d.mu.RLock()
	handlers, exists := d.handlers[event.Name]
	d.mu.RUnlock()

	if !exists {
		return
	}

	// Call all handlers synchronously
	for _, handler := range handlers {
		handler(event)
	}
}

// Unsubscribe removes a handler for an event
// Note: This is a simplified implementation that removes all handlers for an event
func (d *EventDispatcher) Unsubscribe(eventName string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	delete(d.handlers, eventName)
}

// HasSubscribers checks if an event has any subscribers
func (d *EventDispatcher) HasSubscribers(eventName string) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()

	handlers, exists := d.handlers[eventName]
	return exists && len(handlers) > 0
}
