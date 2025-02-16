package event

import (
	"reflect"
	"sync"
)

type EventBus interface {
	Publish(event interface{})
	Subscribe(eventType interface{}, handler func(event interface{}))
	Unsubscribe(eventType interface{}, handler func(event interface{}))
	Stop()
}

type eventBus struct {
	handlers map[string][]func(event interface{})
	mu       sync.RWMutex
	stopChan chan struct{}
	stopOnce sync.Once
}

func NewEventBus() EventBus {
	return &eventBus{
		handlers: make(map[string][]func(event interface{})),
		stopChan: make(chan struct{}),
	}
}

func (b *eventBus) Publish(event interface{}) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	eventType := reflect.TypeOf(event).String()
	if handlers, exists := b.handlers[eventType]; exists {
		for _, handler := range handlers {
			go handler(event) // Non-blocking event handling
		}
	}
}

func (b *eventBus) Subscribe(eventType interface{}, handler func(event interface{})) {
	b.mu.Lock()
	defer b.mu.Unlock()

	typeStr := reflect.TypeOf(eventType).String()
	b.handlers[typeStr] = append(b.handlers[typeStr], handler)
}

func (b *eventBus) Unsubscribe(eventType interface{}, handler func(event interface{})) {
	b.mu.Lock()
	defer b.mu.Unlock()

	typeStr := reflect.TypeOf(eventType).String()
	if handlers, exists := b.handlers[typeStr]; exists {
		for i, h := range handlers {
			if reflect.ValueOf(h).Pointer() == reflect.ValueOf(handler).Pointer() {
				b.handlers[typeStr] = append(handlers[:i], handlers[i+1:]...)
				break
			}
		}
	}
}

// Event definitions
type Event interface {
	EventType() string
}

type PollCreatedEvent struct {
	Poll interface{}
}

func (e PollCreatedEvent) EventType() string {
	return "poll.created"
}

type VoteRecordedEvent struct {
	Poll interface{}
	Vote interface{}
}

func (e VoteRecordedEvent) EventType() string {
	return "vote.recorded"
}

func (b *eventBus) Stop() {
	b.stopOnce.Do(func() {
		close(b.stopChan)
		b.mu.Lock()
		defer b.mu.Unlock()
		b.handlers = make(map[string][]func(event interface{}))
	})
}
