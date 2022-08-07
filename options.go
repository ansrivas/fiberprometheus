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

import "github.com/prometheus/client_golang/prometheus"

type fiberPrometheusOptions struct {
	serviceName string
	namespace   string
	subsystem   string
	labels      map[string]string
	registry    prometheus.Registerer
}

// Option is used to set options when creating a new instance of FiberPrometheus.
//
// Example usage:
//
//     prom := New("my_service", WithNamespace("my_namespace"))
//
type Option func(*fiberPrometheusOptions)

// WithServiceName sets the service name as a constant label.
func WithServiceName(serviceName string) Option {
	return func(o *fiberPrometheusOptions) {
		o.serviceName = serviceName
	}
}

// WithNamespace will prefix the metrics with the given namespace.
func WithNamespace(namespace string) Option {
	return func(o *fiberPrometheusOptions) {
		o.namespace = namespace
	}
}

// WithSubsystem will prefix the metrics with the given subsystem.
func WithSubsystem(subsystem string) Option {
	return func(o *fiberPrometheusOptions) {
		o.subsystem = subsystem
	}
}

// WithLabels will set constant labels for the metrics.
func WithLabels(labels map[string]string) Option {
	return func(o *fiberPrometheusOptions) {
		o.labels = labels
	}
}

// WithRegistry will register the collector with the given registry.
func WithRegistry(registry prometheus.Registerer) Option {
	return func(o *fiberPrometheusOptions) {
		o.registry = registry
	}
}
