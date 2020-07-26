package events

import (
	"context"
	"sync"
)

// Event contains the data for the event and the topic
// it was published on.
type Event struct {
	Data  interface{}
	Topic string
}

// NewEventBus accepts a context and returns an EventBus for use
// in the application.
func NewEventBus(ctx context.Context) *EventBus {
	return &EventBus{
		subscribers: make(map[string][]chan<- Event),
		ctx:         ctx,
	}
}

// DataChannel is a channel of Events
type DataChannel chan Event

// DataChannelSlice is a slice of DataChannels
type DataChannelSlice []DataChannel

// EventBus stores the information about subscribers interested for // a particular topic
type EventBus struct {
	subscribers map[string][]chan<- Event
	rm          sync.RWMutex
	ctx         context.Context
}

// CreateSubscriber simplifies the prcess of setting up a subscriber.
// If you would like you can bypass this function entirely and take full controll
// over the subscriber creation process if that is something you need. The only thing required
// when calling subscribe is that you have an event channel you will listen to.
func (eb *EventBus) CreateSubscriber(listener func(topic string, data interface{})) chan Event {
	eventChannel := make(chan Event)

	go func(ctx context.Context, eventChannel <-chan Event) {
		for {
			select {
			case data := <-eventChannel:
				listener(data.Topic, data.Data)
			case <-ctx.Done():
				return
			}
		}
	}(eb.ctx, eventChannel)

	return eventChannel
}

// Subscribe takes a topic and the channel you will listen on.
func (eb *EventBus) Subscribe(topic string, ch chan<- Event) {
	eb.rm.Lock()

	if prev, found := eb.subscribers[topic]; found {
		eb.subscribers[topic] = append(prev, ch)
	} else {
		eb.subscribers[topic] = append([]chan<- Event{}, ch)
	}

	eb.rm.Unlock()
}

// Publish takes a topic and the data you would like to publish. This data is
// then dispersed to all subscribers of that topic.
func (eb *EventBus) Publish(topic string, data interface{}) {
	eb.rm.RLock()

	if chans, found := eb.subscribers[topic]; found {

		channels := append([]chan<- Event{}, chans...)

		go func(data Event, dataChannelSlices []chan<- Event) {
			for _, ch := range dataChannelSlices {
				ch <- data
			}
		}(Event{Data: data, Topic: topic}, channels)
	}

	eb.rm.RUnlock()
}
