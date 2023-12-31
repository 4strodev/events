# Events
This is a simple event service that allows you to emit events, listen events using glob regex to match event names. Destroy event listeners. It is context aware and allows you to use it for graceful shutdown.

## Installation
```sh
go get github.com/4strodev/events@latest
```

## Usage
Loggin example
```go
package main

import (
    "log"
    "sync"
    "context"
    "github.com/4strodev/events/pkg"
)

func main() {
    // Create new event serivce with default configuration
    eventService := events.NewEventService()
    // The waitgroup will be used only in this example
    // In real applications normally the main goroutine is locked by a server, GUI, etc.
    waitGroup := sync.WaitGroup{}
    
    // You decide when to start the service
    // Note! You can add listeners before start the service.
    // The start method simply starts to process events from event queue
    eventService.Start()

    waitGroup.Add(1)
    // Here we will implement a simple loggin system using events
    eventService.Listen("*", func(event events.Event) error {
        defer waitGroup.Done()
        select {
        // Note that the Event interface provide a context to allow you to end task if the service is stopped
        // this is recomended for those tasks that has external connections:
        // message brokers, databases, external APIs, etc.
        case <-event.Ctx().Done():
            return context.Cause(event.Ctx())
        default:
            log.Println(event.Tag())
            log.Println(event.CreatedAt())
            log.Println(event.Payload())
            return  nil
        }
    })

    event := events.NewBaseEvent("event")
    eventService.Emit(event.WithPayload("A simple payload"))

    waitGroup.Wait()

    // Also you decide when to stop listening events
    // Note! When you stop the event service you cannot start it again. This method exist for graceful shutdowns
    eventService.Stop()
}
```

## More information
[Go doc](https://pkg.go.dev/github.com/4strodev/events)

