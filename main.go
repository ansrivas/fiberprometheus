package main

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

type PrometheusStruct struct {
	reqs       *prometheus.CounterVec
	latency    *prometheus.HistogramVec
	MetricsURL string
}

var (
	// DefaultBuckets prometheus buckets in seconds.
	DefaultBuckets = []float64{0.3, 1.2, 5.0}
)

const (
	reqsName    = "http_requests_total"
	latencyName = "http_request_duration_seconds"
)

func New(name string, buckets ...float64) *PrometheusStruct {
	p := PrometheusStruct{}
	p.reqs = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name:        reqsName,
			Help:        "How many HTTP requests processed, partitioned by status code, method and HTTP path.",
			ConstLabels: prometheus.Labels{"service": name},
		},
		[]string{"code", "method", "path"},
	)
	// prometheus.MustRegister(p.reqs)

	if len(buckets) == 0 {
		buckets = DefaultBuckets
	}

	p.latency = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:        latencyName,
		Help:        "How long it took to process the request, partitioned by status code, method and HTTP path.",
		ConstLabels: prometheus.Labels{"service": name},
		Buckets:     buckets,
	},
		[]string{"code", "method", "path"},
	)
	// prometheus.MustRegister(p.latency)

	p.MetricsURL = "/metrics"
	return &p
}

var (

// p = fasthttpadaptor.NewFastHTTPHandler(
// 	promhttp.HandlerFor(
// 		prometheus.DefaultGatherer,
// 		promhttp.HandlerOpts{
// 			// Opt into OpenMetrics to support exemplars.
// 			EnableOpenMetrics: true,
// 		},
// 	),
// )

// errCnt = prometheus.NewCounterVec(
// 	prometheus.CounterOpts{
// 		Name: "promhttp_metric_handler_errors_total",
// 		Help: "Total number of internal errors encountered by the promhttp metric handler.",
// 	},
// 	[]string{"status", "method", "path"},
// )
)

// func init() {
// 	prometheus.MustRegister(errCnt)
// }

func (ps *PrometheusStruct) prometheusHandler(c *fiber.Ctx) {
	p := fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler())
	p(c.Fasthttp)

}

func (ps *PrometheusStruct) RegisterHandler(app *fiber.App, url string) {
	ps.MetricsURL = url
	app.Get(ps.MetricsURL, ps.prometheusHandler)

}
func (ps *PrometheusStruct) Middleware(ctx *fiber.Ctx) {

	start := time.Now()
	me := string(ctx.Fasthttp.Method())
	path := string(ctx.Fasthttp.Path())

	if path == ps.MetricsURL {
		ctx.Next()
		return
	}

	ctx.Next()

	sc := ctx.Fasthttp.Response.StatusCode()
	statusCode := strconv.Itoa(sc)
	ps.reqs.WithLabelValues(statusCode, me, path).
		Inc()

	ps.latency.WithLabelValues(statusCode, me, path).
		Observe(float64(time.Since(start).Nanoseconds()) / 1000000000)
}

func main() {
	app := fiber.New()

	prometh := New("test-app")
	prometh.RegisterHandler(app, "/metrics")

	app.Use(prometh.Middleware)

	app.Get("/", func(c *fiber.Ctx) {
		c.Send("Hello, World!")
	})

	app.Get("/404", func(c *fiber.Ctx) {
		c.Status(404).Send("You just lost us an unpaid employee.")
	})

	app.Listen(3000)
}
