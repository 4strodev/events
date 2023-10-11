package events

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/gobwas/glob"
	"github.com/satori/uuid"
)

type EventListener func(Event) error
type ErrorHandler func(error) error
type ListenerDestroyer func()

type eventListenerCollection map[string]EventListener

type EventService struct {
	listeners      map[string]eventListenerCollection
	eventQueue     chan Event
	listenersMutex sync.Mutex
	closedMutex    sync.Mutex
	config         EventServiceConfig
	ctx            context.Context
	cancelFunc     context.CancelCauseFunc
	closed         bool
}

// Start turn up the event listener engine.
// At this moment listeners will start reacting to the emitted events
func (e *EventService) Start() {
	go func() {
		for event := range e.eventQueue {
			for listenerTag := range e.listeners {
				g := glob.MustCompile(listenerTag)
				if !g.Match(event.Tag()) {
					continue
				}

				listeners := e.getListeners(listenerTag)
				for _, h := range listeners {
					go func(event Event, listener EventListener) {
						err := listener(event)
						err = e.config.ErrorHandler(err)
						if err != nil {
							log.Fatal(err)
						}
					}(event, h)
				}
			}
		}
	}()
}

// Register a new event listener. Remember tag accepts glob regex.
// Ej.
//
// - listen all events: "*"
//
// - listen all events under a namespace: "users.*"
//
// - listen all created events: "*.created"
//
// For more information see https://github.com/gobwas/glob
func (e *EventService) Listen(tag string, listener EventListener) ListenerDestroyer {
	if e.IsClosed() {
		panic("event service is closed")
	}
	id := e.putListener(tag, listener)
	return func() {
		e.listenersMutex.Lock()
		delete(e.listeners[tag], id)
		e.listenersMutex.Unlock()
	}
}

// Emit will add a provided event to the event queue and all listeners that match with event name
// will react to this event
func (e *EventService) Emit(event Event) {
	if e.IsClosed() {
		panic("event service is closed")
	}
	e.eventQueue <- event.WithCtx(e.ctx)
}

// Stop ends event emitting and listening. At this moment every emit or listen will panic
func (e *EventService) Stop() {
	// First check if event service is closed
	if e.IsClosed() {
		return
	}

	// Close event service
	close(e.eventQueue)
	e.closedMutex.Lock()
	e.closed = true
	e.closedMutex.Unlock()

	// Cancel current listeners
	e.cancelFunc(errors.New("Event service closed"))
}

// getListeners return all listeners registered with a tag
func (e *EventService) getListeners(tag string) eventListenerCollection {
	e.listenersMutex.Lock()
	listeners := e.listeners[tag]
	e.listenersMutex.Unlock()
	return listeners
}

// putListener will add a new listener and returns the random listener id
func (e *EventService) putListener(tag string, listener EventListener) string {
	e.listenersMutex.Lock()
	listeners := e.listeners[tag]
	// Create a random UUID and assing listener
	id := uuid.NewV4().String()
	if listeners == nil {
		listeners = make(eventListenerCollection)
	}
	listeners[id] = listener
	e.listeners[tag] = listeners
	e.listenersMutex.Unlock()

	return id
}

// IsClosed cheks if event service is cosed
func (e *EventService) IsClosed() bool {
	e.closedMutex.Lock()
	closed := e.closed
	e.closedMutex.Unlock()

	return closed
}

// Returns a new event service with default config or custom one
func NewEventService(config ...EventServiceConfig) *EventService {
	var serviceConfig EventServiceConfig

	if len(config) > 0 {
		serviceConfig = config[0]
	} else {
		serviceConfig = GetDefaultConfig()
	}
	ctx, cancel := context.WithCancelCause(context.Background())
	eventService := EventService{
		listeners:      make(map[string]eventListenerCollection),
		eventQueue:     make(chan Event),
		listenersMutex: sync.Mutex{},
		config:         serviceConfig,
		ctx:            ctx,
		cancelFunc:     cancel,
		closed:         false,
	}
	return &eventService
}
