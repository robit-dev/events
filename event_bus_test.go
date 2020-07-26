package events_test

import (
	"events"
	"log"
	"testing"

	"context"

	"github.com/stretchr/testify/assert"
)

// Constant test events
const (
	TEST_EVENT = "TEST_EVENT"
)

type TestEvent struct {
	Cancel context.CancelFunc
	Data   string
}

func testSubscriber(topic string, data interface{}) {
	log.Printf("%s: %s", topic, data.(string))
}

func TestPublishAndSubscribeToEvents(t *testing.T) {
	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a new event bus with the context
	eb := events.NewEventBus(ctx)

	// Create a subscriber that can be used to subscribe to an event (or multiple events)
	subscriber := eb.CreateSubscriber(testSubscriber)

	// Subscribe to test event
	eb.Subscribe(TEST_EVENT, subscriber)

	// Publish an event
	eb.Publish(TEST_EVENT, TestEvent{
		Cancel: cancel,
		Data:   "Test event data",
	})

	// Listen for the event
	event := <-subscriber

	// Assert data received from event matches expectations
	eventStruct := event.Data.(TestEvent)
	assert.Equal(t, TEST_EVENT, event.Topic)
	assert.Equal(t, "Test event data", eventStruct.Data)
}
