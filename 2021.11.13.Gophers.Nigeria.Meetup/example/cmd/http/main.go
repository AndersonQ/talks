package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"time"

	"github.com/AndersonQ/talks/2021.11.13/example"
	"github.com/AndersonQ/talks/2021.11.13/example/config"
	"github.com/AndersonQ/talks/2021.11.13/example/events"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	guuid "github.com/google/uuid"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const serviceName = "http-2nd-example"

func main() {
	// catch the signals as soon as possible
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt) // a.k.a ctrl+C

	cfg := config.Parse(serviceName)

	// OpenTelemetry (OTel) tracer for service A.
	example.SetUp(cfg)

	kp := events.NewProducer(cfg)
	http.Handle("/",
		otelhttp.NewHandler(newHandler("2nd-handler", kp, cfg.Topic, cfg.Logger), "2nd-handler"))

	http.Handle("/final",
		otelhttp.NewHandler(newFinalHandler("final-handler", cfg.Logger), "final-handler"))
	cfg.Logger.Info().Msgf("starting server on port %s", cfg.Port)

	go func() {
		err := http.ListenAndServe(":"+cfg.Port, nil)
		cfg.Logger.Info().Msgf("http server stopped: %v", err)
	}()

	<-signalChan
	cfg.Logger.Info().Msg("Goodbye cruel world!")
	os.Exit(0)
}

func newFinalHandler(name string, logger zerolog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.With().Str("handler", name).Logger()

		// The headers will be sent as part of the response body to show the
		// headers OpenTelemetry uses.
		headers, _ := json.Marshal(r.Header)

		time.Sleep(time.Duration(5*rand.Intn(5)) * time.Millisecond)

		l.Debug().RawJSON("inbound_headers", headers).Msg("headers")
		l.Info().Msgf("/ %s - 200 OK", r.Method)
		w.Header().Set("content-type", "application/json")
		_, _ = w.Write([]byte(
			fmt.Sprintf(`{"inbound_headers":%q}`,
				headers)))
	}
}

var count int64

func newHandler(name string, p *kafka.Producer, topic string, logger zerolog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		sp := trace.SpanFromContext(ctx)
		sp.SetAttributes(
			attribute.String("some_key", "some_value"),
			attribute.Int64("some_int", count))
		traceID := sp.SpanContext().TraceID().String()

		l := logger.With().
			Str("handler", name).
			Str("trace_id", traceID).Logger()

		defer func() { atomic.AddInt64(&count, 1) }()

		uuid, _ := guuid.NewUUID()
		evKey := uuid.String()

		ev := events.Event{
			Key:     evKey,
			Payload: `{"name":"GoGotOlder", "age":12}`,
			Headers: map[string]string{}}
		otel.GetTextMapPropagator().Inject(ctx, ev.Headers)

		bs, _ := json.Marshal(ev.Headers)

		logger.Debug().RawJSON("event_headers", bs)
		if err := produce(p, topic, ev); err != nil {
			l.Err(err).Msg("handler failed: could not produce events")

			sp.RecordError(err, trace.WithStackTrace(true))
			sp.SetStatus(codes.Error, err.Error())

			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprint(w, "Internal server error")
			return
		}

		// simulates a failure by flipping a coin.
		if rand.Int()%2 == 0 {
			err := errors.New(http.StatusText(http.StatusTeapot))
			l.Err(err).Msg("handler failed: bad luck")

			sp.RecordError(err, trace.WithStackTrace(true))
			sp.SetStatus(codes.Error, err.Error())

			w.WriteHeader(http.StatusTeapot)
			_, _ = fmt.Fprintf(w, "I'm a tea pot\ntrace_id: %s", traceID)
			return
		}

		// The headers will be sent as part of the response body to show the
		// headers OpenTelemetry uses.
		headers, _ := json.Marshal(r.Header)

		time.Sleep(time.Duration(5*rand.Intn(5)) * time.Millisecond)

		l.Info().
			RawJSON("inbound_headers", headers).
			Msgf("/ %s - 200 OK", r.Method)
		w.Header().Set("content-type", "application/json")
		_, _ = w.Write([]byte(
			fmt.Sprintf(`{"trace_id":"%s","inbound_headers":%q}`,
				traceID,
				headers)))
	}
}

func produce(producer *kafka.Producer, topic string, e events.Event) error {
	return producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(e.Key),
		Value:          []byte(e.Payload),
		Headers:        e.Headers.KafkaHeader(),
	}, nil)
}
