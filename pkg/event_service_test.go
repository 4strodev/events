package events

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestEventServiceEmit checks if all event listeners are reacting to the event
// the test adds 10 event listeners to the service. This listeners will increment
// a counter. Then after the emit the counter is 10 will mean that all listeners
// reacted to the event.
func TestEventServiceEmit(t *testing.T) {
	waitGroup := sync.WaitGroup{}
	eventService := NewEventService()
	var activatedListeners int
	var resultsChannel = make(chan struct{})
	eventService.Start()

	go func() {
		for {
			select {
			case <-resultsChannel:
				activatedListeners++
				waitGroup.Done()
			}
		}
	}()

	for i := 0; i < 10; i++ {
		waitGroup.Add(1)
		eventService.Listen("*", func(event Event) error {
			resultsChannel <- struct{}{}
			return nil
		})
	}

	event := NewBaseEvent("event")
	eventService.Emit(&event)
	waitGroup.Wait()
	assert.Equal(t, 10, activatedListeners)
}

// TestEventServiceListenerDestroyer checks if the returned function
// to destroy event listeners works. The test counts how many listeners
// reacted to the event. The test starts with 10 and then removes 1.
// The test chekcs if the first emit count 10 reactions. The second
// run should count 9 reactions.
func TestEventServiceListenerDestroyer(t *testing.T) {
	waitGroup := sync.WaitGroup{}
	eventService := NewEventService()
	var activatedListeners int
	var resultsChannel = make(chan struct{})
	eventService.Start()

	go func() {
		for {
			select {
			case <-resultsChannel:
				activatedListeners++
				waitGroup.Done()
			}
		}
	}()

	waitGroup.Add(1)
	destroyer := eventService.Listen("*", func(event Event) error {
		resultsChannel <- struct{}{}
		return nil
	})

	for i := 0; i < 9; i++ {
		waitGroup.Add(1)
		eventService.Listen("*", func(event Event) error {
			resultsChannel <- struct{}{}
			return nil
		})
	}

	event := NewBaseEvent("event")

	eventService.Emit(&event)
	waitGroup.Wait()
	assert.Equal(t, 10, activatedListeners)

	destroyer()
	activatedListeners = 0
	waitGroup.Add(9)
	eventService.Emit(&event)
	waitGroup.Wait()
	assert.Equal(t, 9, activatedListeners)
}

// TestEventServiceStop checks if the service after sopping it do not admit more events
func TestEventServiceStop(t *testing.T) {
	//waitGroup := sync.WaitGroup{}
	eventService := NewEventService()
	eventService.Start()
	eventService.Stop()
}

