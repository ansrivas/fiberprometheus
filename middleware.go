package fiberprometheus

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

// FiberPrometheus ...
type FiberPrometheus struct {
	requestsTotal   *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
	requestInFlight *prometheus.GaugeVec
	defaultURL      string
}

// New creates a new instance of FiberPrometheus middleware
func New(servicename string) *FiberPrometheus {
	counter := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name:        "http_requests_total",
			Help:        "Count all http requests by status code, method and path.",
			ConstLabels: prometheus.Labels{"service": servicename},
		},
		[]string{"status_code", "method", "path"},
	)
	histogram := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:        "http_request_duration_seconds",
		Help:        "Duration of all HTTP requests by status code, method and path.",
		ConstLabels: prometheus.Labels{"service": servicename},
	},
		[]string{"status_code", "method", "path"},
	)

	gauge := promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name:        "http_requests_in_progress_total",
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

func (ps *FiberPrometheus) handler(c *fiber.Ctx) {
	p := fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler())
	p(c.Fasthttp)
}

// RegisterAt will register the prometheus handler at a given URL
func (ps *FiberPrometheus) RegisterAt(app *fiber.App, url string) {
	ps.defaultURL = url
	app.Get(ps.defaultURL, ps.handler)
}

// Middleware is the actual default middleware implementation
func (ps *FiberPrometheus) Middleware(ctx *fiber.Ctx) {

	start := time.Now()
	method := string(ctx.Fasthttp.Method())
	path := string(ctx.Fasthttp.Path())

	if path == ps.defaultURL {
		ctx.Next()
		return
	}

	ps.requestInFlight.WithLabelValues(method, path).Inc()
	ctx.Next()
	ps.requestInFlight.WithLabelValues(method, path).Dec()

	statusCode := strconv.Itoa(ctx.Fasthttp.Response.StatusCode())
	ps.requestsTotal.WithLabelValues(statusCode, method, path).
		Inc()

	elapsed := float64(time.Since(start).Nanoseconds()) / 1000000000
	ps.requestDuration.WithLabelValues(statusCode, method, path).
		Observe(elapsed)

}
