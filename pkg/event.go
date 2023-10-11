package pkg

import (
	"context"
	"time"
)

// Event is the main interface where the library turns around
// the service expects an Event interface and you can implement your own events
type Event interface {
	// Tag returns the event name
	Tag() string
	// CreatedAt returns the event creation time
	CreatedAt() time.Time
	// Payload returns the event payload
	Payload() any
	// Ctx returns event context
	Ctx() *context.Context
}

// The base event implements the Event interface adding some basic fields
type BaseEvent struct {
	tag      string
	createdAt time.Time
	payload   any
	ctx       *context.Context
}

func NewBaseEvent(name string) BaseEvent {
	return BaseEvent {
		tag: name,
		createdAt: time.Now(),
	}
}

func (b *BaseEvent) Tag() string {
	return b.tag
}

func (b *BaseEvent) CreatedAt() time.Time {
	return b.createdAt
}

func (b *BaseEvent) Payload() any {
	return b.payload
}

func (b *BaseEvent) Ctx() *context.Context {
	return b.ctx
}

func (b *BaseEvent) WithPayload(payload any) *BaseEvent {
	b.payload = payload
	return b
}
