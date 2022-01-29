//
// Copyright (c) 2021-present Ankur Srivastava and Contributors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

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

// Register FiberPrometheus Collectors to
// a custom Registry
func (fp *FiberPrometheus) MustRegister(registry *prometheus.Registry) {
	registry.MustRegister(
		fp.requestsTotal,
		fp.requestDuration,
		fp.requestInFlight,
	)
}

func create(servicename, namespace, subsystem string, labels map[string]string) *FiberPrometheus {
	constLabels := make(prometheus.Labels)
	if servicename != "" {
		constLabels["service"] = servicename
	}
	for label, value := range labels {
		constLabels[label] = value
	}

	counter := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name:        prometheus.BuildFQName(namespace, subsystem, "requests_total"),
			Help:        "Count all http requests by status code, method and path.",
			ConstLabels: constLabels,
		},
		[]string{"status_code", "method", "path"},
	)
	histogram := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:        prometheus.BuildFQName(namespace, subsystem, "request_duration_seconds"),
		Help:        "Duration of all HTTP requests by status code, method and path.",
		ConstLabels: constLabels,
		Buckets: []float64{
			0.000000001, // 1ns
			0.000000002,
			0.000000005,
			0.00000001, // 10ns
			0.00000002,
			0.00000005,
			0.0000001, // 100ns
			0.0000002,
			0.0000005,
			0.000001, // 1µs
			0.000002,
			0.000005,
			0.00001, // 10µs
			0.00002,
			0.00005,
			0.0001, // 100µs
			0.0002,
			0.0005,
			0.001, // 1ms
			0.002,
			0.005,
			0.01, // 10ms
			0.02,
			0.05,
			0.1, // 100 ms
			0.2,
			0.5,
			1.0, // 1s
			2.0,
			5.0,
			10.0, // 10s
			15.0,
			20.0,
			30.0,
		},
	},
		[]string{"status_code", "method", "path"},
	)

	gauge := promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name:        prometheus.BuildFQName(namespace, subsystem, "requests_in_progress_total"),
		Help:        "All the requests in progress",
		ConstLabels: constLabels,
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
	return create(servicename, "http", "", nil)
}

// NewWith creates a new instance of FiberPrometheus middleware but with an ability
// to pass namespace and a custom subsystem
// Here servicename is created as a constant-label for the metrics
// Namespace, subsystem get prefixed to the metrics.
//
// For e.g namespace = "my_app", subsyste = "http" then then metrics would be
// `my_app_http_requests_total{...,service= "servicename"}`
func NewWith(servicename, namespace, subsystem string) *FiberPrometheus {
	return create(servicename, namespace, subsystem, nil)
}

// NewWithLabels creates a new instance of FiberPrometheus middleware but with an ability
// to pass namespace and a custom subsystem
// Here labels are created as a constant-labels for the metrics
// Namespace, subsystem get prefixed to the metrics.
//
// For e.g namespace = "my_app", subsystem = "http" and lables = map[string]string{"key1": "value1", "key2":"value2"}
// then then metrics would become
// `my_app_http_requests_total{...,key1= "value1", key2= "value2" }``
func NewWithLabels(labels map[string]string, namespace, subsystem string) *FiberPrometheus {
	return create("", namespace, subsystem, labels)
}

// RegisterAt will register the prometheus handler at a given URL
func (ps *FiberPrometheus) RegisterAt(app *fiber.App, url string) {
	ps.defaultURL = url
	app.Get(ps.defaultURL, adaptor.HTTPHandler(promhttp.Handler()))
}

// Middleware is the actual default middleware implementation
func (ps *FiberPrometheus) Middleware(ctx *fiber.Ctx) error {

	start := time.Now()
	method := ctx.Route().Method
	path := ctx.Route().Path

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
