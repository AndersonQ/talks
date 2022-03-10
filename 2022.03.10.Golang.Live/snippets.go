package main

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/blacklane/go-libs/x/events"
)

// start_event_type OMIT
type Header map[string]string

type Event struct {
	Headers Header
	Key     []byte
	Payload []byte
}

type Handler interface {
	Handle(context.Context, Event) error
}

// end_event_type OMIT

// start_event_helpers OMIT
type HandlerFunc func(ctx context.Context, e Event) error

func (h HandlerFunc) Handle(ctx context.Context, e Event) error {
	return h(ctx, e)
}

type Middleware func(Handler) Handler

// end_event_helpers OMIT

// start_handler_builder OMIT

type HandlerBuilder struct {
	middleware  []Middleware
	rawHandlers []Handler
}

func (hb *HandlerBuilder) UseMiddleware(m ...Middleware) {
	hb.middleware = append(hb.middleware, m...)
}

func (hb *HandlerBuilder) AddHandler(h Handler) {
	hb.rawHandlers = append(hb.rawHandlers, h)
}

func (hb HandlerBuilder) Build() []Handler {
	var handlers []Handler
	for _, rh := range hb.rawHandlers {
		h := rh
		for i := len(hb.middleware) - 1; i >= 0; i-- {
			h = hb.middleware[i](h)
		}
		handlers = append(handlers, h)
	}
	return handlers
}

// end_handler_builder OMIT

// mock code, not real code, just for the compiler stop complaining
type kafkaConsumer struct {
	wg       sync.WaitGroup
	handlers []Handler
	consumer consumer
	done     chan struct{}
}

type consumer int

func (c consumer) Poll(int) interface{} {
	return nil
}

func messageToEvent(msg interface{}) Event {
	return Event{}
}

func (c *kafkaConsumer) running() bool {
	return true
}

func (c *consumer) Shutdown() error {}
func (c *consumer) Run()            {}

// end mock code

// start_run_loop OMIT
func (c *kafkaConsumer) Run(timeout time.Duration) {
	go func() { // HL_loop_init
		for c.running() { // HL_loop_init
			msg := c.consumer.Poll(int(timeout.Milliseconds())) // HL_get_message
			for _, h := range c.handlers {                      // HL_get_message
				if msg == nil {
					// when c.consumer.Poll(timeoutMs) times out, it returns nil.
					continue
				}
				c.wg.Add(1)          // HL_shutdown
				go func(h Handler) { // HL_handle
					defer c.wg.Done() // HL_shutdown
					e := messageToEvent(msg)

					// Errors are ignored, a middleware or the handler should handle them
					_ = h.Handle(context.Background(), e) // HL_handle
				}(h)
			}
		} // HL_shutdown

		c.wg.Wait() // HL_shutdown
		close(c.done)
	}()
}

// end_run_loop OMIT

// start_consumer OMIT
func consume() {
	// catch the signals as soon as possible
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt) // a.k.a ctrl+C

	conf := &kafka.ConfigMap{
		"group.id":           "KafkaGroupID",
		"bootstrap.servers":  "localhost:9092",
		"session.timeout.ms": 6000,
		"auto.offset.reset":  "earliest",
	}

	// start_consumer_handler OMIT
	topic := "my-topic"
	c, err := events.NewKafkaConsumer(
		events.NewKafkaConsumerConfig(conf),
		[]string{topic},
		events.HandlerFunc( // HL
			func(ctx context.Context, e events.Event) error { // HL
				log.Printf("consumed event: %s", e.Payload) // HL
				return nil                                  // HL
			}))
	if err != nil {
		panic(err)
	}

	log.Printf("starting to consume events from %s, press CTRL+C to exit", topic)
	c.Run(time.Second) // HL_run

	<-shutdown                                              // HL_shutdown
	if err = c.Shutdown(context.Background()); err != nil { // HL_shutdown
		log.Printf("Shutdown error: %v", err)
		os.Exit(1)
	}
	log.Printf("Shutdown successfully")
	// end_consumer_handler OMIT

}

// start_consumer_interface OMIT
type Consumer interface {
	Run(timeout time.Duration)
	Shutdown(ctx context.Context) error
}

// end_consumer_interface OMIT

// start_producer_interface OMIT
type Producer interface {
	// Send sends an event to the given topic
	// Deprecated. use SendCtx instead
	Send(event Event, topic string) error // HL
	// SendCtx send an event to the given topic.
	// It also adds the OTel propagation headers and the X-Tracking-Id if not set
	// already.
	SendCtx(ctx context.Context, eventName string, event Event, topic string) error // HL
	// SendWithTrackingID adds the tracking ID to the event's headers and sends
	// it to the given topic
	// Deprecated. use SendCtx instead
	SendWithTrackingID(trackingID string, event Event, topic string) error
	// HandleEvents starts to listen to the producer events channel
	HandleEvents() error
	// Shutdown gracefully shuts down the producer, it respects the context timeout.
	Shutdown(ctx context.Context) error // HL
}

// end_producer_interface OMIT
