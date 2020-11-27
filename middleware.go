package fiberprometheus

import (
	"strconv"
	"time"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// FiberPrometheus ...
type FiberPrometheus struct {
	requestsTotal   *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
	requestInFlight *prometheus.GaugeVec
	defaultURL      string
}

func create(servicename, namespace, subsystem string) *FiberPrometheus {
	counter := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name:        prometheus.BuildFQName(namespace, subsystem, "requests_total"),
			Help:        "Count all http requests by status code, method and path.",
			ConstLabels: prometheus.Labels{"service": servicename},
		},
		[]string{"status_code", "method", "path"},
	)
	histogram := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:        prometheus.BuildFQName(namespace, subsystem, "request_duration_seconds"),
		Help:        "Duration of all HTTP requests by status code, method and path.",
		ConstLabels: prometheus.Labels{"service": servicename},
	},
		[]string{"status_code", "method", "path"},
	)

	gauge := promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name:        prometheus.BuildFQName(namespace, subsystem, "requests_in_progress_total"),
		Help:        "All the requests in progress",
		ConstLabels: prometheus.Labels{"service": servicename},
	}, []string{"method", "path"})

	return &FiberPrometheus{
		requestsTotal:   counter,
		requestDuration: histogram,
		requestInFlight: gauge,
		defaultURL:      "/metrics",
	}
}

// New creates a new instance of FiberPrometheus middleware
// servicename is available as a const label
func New(servicename string) *FiberPrometheus {
	return create(servicename, "http", "")
}

// NewWith creates a new instance of FiberPrometheus middleware but with an ability
// to pass namespace and a custom subsystem
// Here servicename is created as a constant-label for the metrics
// Namespace, subsystem get prefixed to the metrics.
// For e.g namespace = "my_app", subsyste = "http" then then metrics would be
// my_app_http_requests_total{...,"service": servicename}
func NewWith(servicename, namespace, subsystem string) *FiberPrometheus {
	return create(servicename, namespace, subsystem)
}

// RegisterAt will register the prometheus handler at a given URL
func (ps *FiberPrometheus) RegisterAt(app *fiber.App, url string) {
	ps.defaultURL = url
	app.Get(ps.defaultURL, adaptor.HTTPHandler(promhttp.Handler()))
}

// Middleware is the actual default middleware implementation
func (ps *FiberPrometheus) Middleware(ctx *fiber.Ctx) error {

	start := time.Now()
	method := string(ctx.Method())
	path := string(ctx.Path())

	if path == ps.defaultURL {
		return ctx.Next()

	}

	ps.requestInFlight.WithLabelValues(method, path).Inc()
	defer func() {
		ps.requestInFlight.WithLabelValues(method, path).Dec()
	}()
	if err := ctx.Next(); err != nil {
		return err
	}

	statusCode := strconv.Itoa(ctx.Response().StatusCode())
	ps.requestsTotal.WithLabelValues(statusCode, method, path).
		Inc()

	elapsed := float64(time.Since(start).Nanoseconds()) / 1000000000
	ps.requestDuration.WithLabelValues(statusCode, method, path).
		Observe(elapsed)

	return nil
}
