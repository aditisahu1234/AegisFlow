package runtime

import (
	"sync"
	"time"
)

type EventType string

const (
	EventStarting EventType = "starting"
	EventHealthy  EventType = "healthy"
	EventFailed   EventType = "failed"
	EventStopping EventType = "stopping"
	EventStopped  EventType = "stopped"
)

type Event struct {
	Time      time.Time
	Component string
	Type      EventType
	Err       error
}

type EventBus struct {
	mu sync.RWMutex

	subscribers []chan Event
}

func NewEventBus() *EventBus {

	return &EventBus{}
}

func (b *EventBus) Subscribe() <-chan Event {

	ch := make(chan Event, 100)

	b.mu.Lock()
	b.subscribers = append(
		b.subscribers,
		ch,
	)
	b.mu.Unlock()

	return ch
}

func (b *EventBus) Publish(e Event) {

	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, sub := range b.subscribers {

		select {

		case sub <- e:

		default:
		}
	}
}
