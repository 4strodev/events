package events

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
	// WithCtx sets context for event.
	// This method will be used by the
	// events service when you emit an event.
	// This means that if you put a context this will be overrited.
	WithCtx(context.Context) Event
	// Ctx returns event context
	Ctx() context.Context
}

// The base event implements the Event interface adding some basic fields
type BaseEvent struct {
	tag       string
	createdAt time.Time
	payload   any
	ctx       context.Context
}

func NewBaseEvent(name string) BaseEvent {
	return BaseEvent{
		tag:       name,
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

func (b *BaseEvent) Ctx() context.Context {
	return b.ctx
}

func (b *BaseEvent) WithCtx(ctx context.Context) Event {
	b.ctx = ctx
	return b
}

func (b *BaseEvent) WithPayload(payload any) *BaseEvent {
	b.payload = payload
	return b
}
