package main

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/AndersonQ/talks/2021.11.13/example"
	"github.com/AndersonQ/talks/2021.11.13/example/config"
	"github.com/AndersonQ/talks/2021.11.13/example/events"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const serviceName = "event-example"

// OpenTelemetry attribute keys.
const (
	AttrKeyEventKey = attribute.Key("event.key")
)

func main() {
	// catch the signals as soon as possible
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt) // a.k.a ctrl+C

	cfg := config.Parse(serviceName)

	example.SetUp(cfg)

	ch := events.NewConsumer(cfg).Run()

	go func() {
		for e := range ch {
			l := cfg.Logger.With().
				Str("handler", "event_handler").Logger()

			go handle(l.WithContext(context.Background()), e)
		}
	}()

	<-signalChan

	cfg.Logger.Info().Msg("Goodbye cruel world!")
	os.Exit(0)
}

func handle(ctx context.Context, e events.Event) {
	ctx = otel.GetTextMapPropagator().Extract(ctx, e.Headers)
	ctx, sp := otel.Tracer(serviceName).Start(ctx, "events-handler",
		trace.WithSpanKind(trace.SpanKindConsumer),
		trace.WithAttributes(
			AttrKeyEventKey.String(e.Key)))
	defer sp.End()

	traceID := sp.SpanContext().TraceID().String()
	l := zerolog.Ctx(ctx).With().Str("trace_id", traceID).Logger()

	time.Sleep(time.Duration(5*rand.Intn(5)) * time.Millisecond)

	url := "http://localhost:1618/final"
	req, _ := http.NewRequest(http.MethodGet, url, nil)

	// Inject OTel headers to propagate the t
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	// The headers will be sent as part of the response body to show the
	// headers OpenTelemetry uses.
	headers, _ := json.Marshal(req.Header)

	_, err := http.DefaultClient.Do(req)
	if err != nil {

		sp.RecordError(err, trace.WithStackTrace(true))
		sp.SetStatus(codes.Error, err.Error())

		l.Err(err).
			RawJSON("outbound_headers", headers).
			Msgf("handler failed: could not call %s", url)

		return
	}

	l.Info().
		RawJSON("event_payload", []byte(e.Payload)).
		Interface("event_headers", e.Headers).
		Str("event_key", e.Key).
		Msg("processed events")
}
