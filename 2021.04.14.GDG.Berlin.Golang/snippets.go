package main

import (
	"context"
	"sync"
	"time"
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

type kafkaConsumer struct {
	wg       sync.WaitGroup
	handlers []Handler
	consumer consummer
	done     chan struct{}
}

type consummer int

func (c consummer) Poll(int) interface{} {
	return nil
}

// start_run_loop OMIT
func (c *kafkaConsumer) Run(timeout time.Duration) {
	go func() {
		for c.running() {

			msg := c.consumer.Poll(int(timeout.Milliseconds()))
			for _, h := range c.handlers {
				c.wg.Add(1)
				go func(h Handler) {
					defer c.wg.Done()
					e := messageToEvent(msg)

					// Errors are ignored, a middleware or the handler should handle them
					_ = h.Handle(context.Background(), e)
				}(h)
			}
		}

		c.wg.Wait()
		close(c.done)
	}()
}

// end_run_loop OMIT

func messageToEvent(msg interface{}) Event {
	return Event{}
}

func (c *kafkaConsumer) running() bool {
	return true
}
