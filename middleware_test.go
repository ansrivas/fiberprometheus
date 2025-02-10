package fiberprometheus

import (
	"context"
	"fmt"
	"io"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

var Tracer = otel.Tracer("FOCUZ")

func otelTracingInit(t *testing.T) {
	// Add trace resource attributes
	res, err := resource.New(
		context.Background(),
		resource.WithTelemetrySDK(),
		resource.WithOS(),
		resource.WithHost(),
		resource.WithAttributes(attribute.String("service.name", "fiber")),
	)
	if err != nil {
		t.Errorf("cant create otlp resource: %v", err)
		t.Fail()
	}

	// Create stdout exporter
	traceExporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		t.Errorf("cant create otlp exporter: %v", err)
		t.Fail()
	}

	// Create OTEL trace provider
	tp := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter),
		trace.WithResource(res),
	)

	os.Setenv("OTEL_TRACES_EXPORTER", "otlp")
	os.Setenv("OTEL_TRACES_SAMPLER", "always_on")

	// Set OTLP Provider
	otel.SetTracerProvider(tp)

	// SetTextMapPropagator configures the OpenTelemetry text map propagator
	// using a composite of TraceContext and Baggage propagators.
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
}

func tracingMiddleware(c *fiber.Ctx) error {
	// Create a new context with cancellation capability from Fiber context
	ctx, cancel := context.WithCancel(c.UserContext())

	// Start a new span with attributes for tracing the current request
	_, span := Tracer.Start(ctx, c.Route().Name)

	fmt.Println("spanid", span.SpanContext().SpanID())

	// Ensure the span is ended and context is cancelled when the request completes
	defer span.End()
	defer cancel()

	// Continue with the next middleware/handler
	return c.Next()
}

func TestMiddlewareWithExamplar(t *testing.T) {
	t.Parallel()

	otelTracingInit(t)

	app := fiber.New()

	prometheus := New("test-service")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(tracingMiddleware)
	app.Use(prometheus.Middleware)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World")
	})

	req := httptest.NewRequest("GET", "/", nil)
	resp, _ := app.Test(req, -1)
	if resp.StatusCode != 200 {
		t.Fail()
	}

	time.Sleep(15 * time.Second)

	req = httptest.NewRequest("GET", "/metrics", nil)
	req.Header.Set("Accept", "application/openmetrics-text")
	resp, _ = app.Test(req, -1)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	got := string(body)

	t.Logf("debug: %v", got)

	// Check Metrics Response
	want := `http_requests_total{method="GET",path="/",service="test-service",status_code="200"} 1`
	if !strings.Contains(got, want) {
		t.Errorf("got %s; want %s", got, want)
	}

	want = `http_request_duration_seconds_count{method="GET",path="/",service="test-service",status_code="200"} 1`
	if !strings.Contains(got, want) {
		t.Errorf("got %s; want %s", got, want)
	}

	want = `http_requests_in_progress_total{method="GET",service="test-service"} 0`
	if !strings.Contains(got, want) {
		t.Errorf("got %s; want %s", got, want)
	}
}
