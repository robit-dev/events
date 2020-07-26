# Robit Event Bus
Robit Event Bus is a simple channel driven event bus for golang projects.

## Motivation
Event buses are a powerful way of allowing different packages to communicate. At Robit we wanted something simple that required very little boilerplatee code to integrate in projects.


## Installation
All you need to get started is a golang project and then you can simply install with
`go get -u github.com/robit-dev/events`

## Tests
Tests can be run with `go test ./...`

## How to use?

```golang
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/robit-dev/events"
)

func testSubscriber(topic string, data interface{}) {
	// Event contains the TestEventData as well as the topic the event was published on
	fmt.Println(topic)         // "TEST_EVENT"
	fmt.Println(data.(string)) // "Test event data"
}

func main() {
	// Create a new event bus with the context
	eb := events.NewEventBus(context.TODO())

	// Create a subscriber that can be used to subscribe to an event (or multiple events)
	subscriber := eb.CreateSubscriber(testSubscriber)

	// Subscribe to test event
	eb.Subscribe("TEST_EVENT", subscriber)

	// Publish an event (an event can be anything)
	eb.Publish("TEST_EVENT", "Test event data")

	// Sleep to see the output
	time.Sleep(time.Second * 1)
}
```

## Contribute
Contributing is easy, raise an issue or submit a PR to the repo. If you are planning on working on something large it is recommended that you raise an issue first to make sure it is in line with the vision for the project.

## License
This code is released under the MIT License.

## Stack
<b>Built with</b> [Golang](https://golang.org/)

MIT © [Robit Development, LLC](https://robit.dev)