package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"time"

	"github.com/AndersonQ/talks/2021.11.13/example"
	"github.com/AndersonQ/talks/2021.11.13/example/config"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const serviceName = "http-root-example"

func main() {
	// catch the signals as soon as possible
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt) // a.k.a ctrl+C

	cfg := config.Parse(serviceName)

	// OpenTelemetry (OTel) tracer for service.
	example.SetUp(cfg)

	http.HandleFunc("/", newHandler(cfg.Logger))
	cfg.Logger.Info().Msgf("starting server on port %s", cfg.Port)

	go func() {
		err := http.ListenAndServe(":"+cfg.Port, nil)
		cfg.Logger.Info().Msgf("http server stopped: %v", err)
	}()

	<-signalChan
	cfg.Logger.Info().Msg("Goodbye cruel world!")
	os.Exit(0)

}

var count int64

func newHandler(logger zerolog.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, sp := otel.Tracer(serviceName).Start(r.Context(), "root-handler")
		defer sp.End()

		traceID := sp.SpanContext().TraceID().String()
		l := logger.With().Str("trace_id", traceID).Logger()

		sp.SetAttributes(
			attribute.String("some_key", "some_value"),
			attribute.Int64("some_int", count))

		defer func() { atomic.AddInt64(&count, 1) }()

		time.Sleep(time.Duration(5*rand.Intn(5)) * time.Millisecond)

		url := "http://localhost:1618"
		req, _ := http.NewRequest(http.MethodGet, url, nil)

		// Inject OTel headers to propagate the t
		otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

		// The headers will be sent as part of the response body to show the
		// headers OpenTelemetry uses.
		headers, _ := json.Marshal(req.Header)

		_, err := http.DefaultClient.Do(req)
		if err != nil {
			l.Err(err).
				RawJSON("outbound_headers", headers).
				Msgf("handler failed: could not call %s", url)

			sp.RecordError(err, trace.WithStackTrace(true))
			sp.SetStatus(codes.Error, err.Error())

			w.WriteHeader(http.StatusFailedDependency)
			_, _ = fmt.Fprintf(w, "failed to call %s\ntrace_id: %s\noutbound_headers: %s",
				url,
				traceID,
				headers)
			return
		}

		l.Info().RawJSON("outbound_headers", headers).Msgf("/ %s - 200 OK", r.Method)
		w.Header().Set("content-type", "application/json")
		_, _ = w.Write([]byte(
			fmt.Sprintf(`{"trace_id":"%s","outbound_headers":%q}`,
				traceID,
				headers)))
	}
}
