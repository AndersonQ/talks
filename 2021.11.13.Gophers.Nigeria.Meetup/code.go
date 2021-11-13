package main

import (
	"context"
	"log"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

func main() {
	// start_create_span1 OMIT
	tracer := otel.Tracer("ex.com/my/application")
	ctx, span := tracer.Start(context.Background(), "my span name") // HL
	defer span.End()
	// end_create_span1 OMIT

	// start-span-from-ctx OMIT
	span := trace.SpanFromContext(ctx)
	// end-span-from-ctx OMIT

	// start-span-attribute OMIT
	fooKey := attribute.Key("foo")            // HL
	bar42KeyValue := attribute.Int("bar", 42) // HL

	ctx, span := tracer.Start(
		ctx,
		"my span name",
		trace.WithAttributes( // HL
			fooKey.String("fooValue"),
			bar42KeyValue,
			attribute.Bool("admin", false)))
	defer span.End()
	// end-span-attribute OMIT

	// start-attribute-1 OMIT

	span.SetAttributes(
		fooKey.String("fooValue"),
		bar42KeyValue,
		attribute.Bool("admin", false))
	// end-attribute-1 OMIT

	// start-spevent OMIT
	span.AddEvent("some-event", // HL
		trace.WithStackTrace(true),
		trace.WithAttributes(attribute.String("event-attr-key", "event-ettr-value")))
	// end-spevent OMIT

	// start-propagate-http OMIT
	r, _ := http.NewRequest(http.MethodGet, "https://ex.com", nil)

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(r.Header)) // HL
	// end-propagate-http OMIT

}

type Handler struct{}

// start-receive-http OMIT
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header)) // HL

	ctx, span := otel.Tracer("ex.com/my/app/http/handler").
		Start(ctx, "handler name", trace.WithAttributes( /* ... */ ))
	defer span.End()

	// ...
	return
}

// end-receive-http OMIT

type EventsHandler struct{}

const AttrKeyEventName = attribute.Key("events.name")

// start-receive-events OMIT
func (h *EventsHandler) Handle(ctx context.Context, e Event) error {
	tr := otel.Tracer("ex.com/my/app/events/handler")

	ctx = otel.GetTextMapPropagator().Extract(ctx, e.Headers) // HL

	ctx, sp := tr.Start(
		ctx,
		"events name",
		trace.WithSpanKind(trace.SpanKindConsumer),
		trace.WithAttributes(AttrKeyEventName.String("events name")),
	)
	defer sp.End()

	// ...
	return nil
}

// end-receive-events OMIT

// start-events OMIT
type Header map[string]string // HL
type Event struct {
	Headers Header
	Key     []byte
	Payload []byte
}

// end-events OMIT

// start-events-propagation OMIT
// Ensure Handler implements OTel propagation.TextMapCarrier
var _ = propagation.TextMapCarrier(Header{})

func (h Header) Get(key string) string {
	return h[key]
}

func (h Header) Set(key, value string) {
	h[key] = value
}

func (h Header) Keys() []string {
	keys := make([]string, len(h))
	for k := range h {
		keys = append(keys, k)
	}

	return keys
}

// end-events-propagation OMIT

func export() {

	// start-resource OMIT
	res, err := resource.New(context.TODO(), // HL
		resource.WithAttributes( // HL
			semconv.ServiceNameKey.String("my awesome service"),
			semconv.ServiceVersionKey.String("develop"),
			semconv.DeploymentEnvironmentKey.String("local"),
			attribute.String("resource_attribute", "added as resource"),
		),
	)
	if err != nil {
		log.Fatalf("failed to create otel sdk/resource: %v", err)
	}
	// end-resource OMIT

	// start-client OMIT
	otlpClient := otlptracegrpc.NewClient( // HL
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint("localhost:55680"))
	// end-client OMIT

	// start-exporter OMIT

	otlpExporter, err := otlptrace.New(context.TODO(), otlpClient) // HL
	if err != nil {
		log.Fatalf("failed to create OTel exporter, disabling OTel: %v", err)
		return
	}
	// end-exporter OMIT

	// start-stdoutexporter OMIT
	stdoutExporter, err := stdouttrace.New()
	if err != nil {
		log.Fatalf("failed to create stdouttrace exporter: %v", err)
	}
	// end-stdoutexporter OMIT

	// start-traceprovider OMIT
	tracerProvider := trace.NewTracerProvider( // HL_provider
		trace.WithSampler(trace.AlwaysSample()), // HL_Sampler
		trace.WithResource(res),                 // HL_resource
		trace.WithSyncer(otlpExporter),          // HL_exporter
	)
	// end-traceprovider OMIT

	// start-traceprovider-2 OMIT
	tracerProvider := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithResource(res),
		trace.WithBatcher(otlpExporter), // HL
	)
	// end-traceprovider-2 OMIT

	// start-debug-exporter OMIT
	if debug {
		tracerProvider.RegisterSpanProcessor(
			trace.NewSimpleSpanProcessor(stdoutExporter))
	}
	// end-debug-exporter OMIT

	// start-register OMIT
	otel.SetTracerProvider(tracerProvider) // HL
	otel.SetTextMapPropagator(             // HL
		propagation.NewCompositeTextMapPropagator(
			propagation.Baggage{},
			propagation.TraceContext{},
		),
	)
	// end-register OMIT

	// start-getglobals OMIT
	otel.GetTracerProvider()
	otel.GetTextMapPropagator()
	// end-getglobals OMIT

}
