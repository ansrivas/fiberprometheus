//
// Copyright (c) 2020 Ankur Srivastava
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
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestMiddleware(t *testing.T) {
	app := fiber.New()

	prometheus := New("test-service")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World")
	})
	req := httptest.NewRequest("GET", "/", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != 200 {
		t.Fail()
	}

	req = httptest.NewRequest("GET", "/metrics", nil)
	resp, _ = app.Test(req)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	got := string(body)
	want := `http_requests_total{method="GET",path="/",service="test-service",status_code="200"} 1`
	if !strings.Contains(got, want) {
		t.Errorf("got %s; want %s", got, want)
	}

	want = `http_request_duration_seconds_count{method="GET",path="/",service="test-service",status_code="200"} 1`
	if !strings.Contains(got, want) {
		t.Errorf("got %s; want %s", got, want)
	}

	want = `http_requests_in_progress_total{method="GET",path="/",service="test-service"} 0`
	if !strings.Contains(got, want) {
		t.Errorf("got %s; want %s", got, want)

	}
}

func TestMiddlewareWithServiceName(t *testing.T) {
	app := fiber.New()

	prometheus := NewWith("test-service", "my_service", "http")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World")
	})
	req := httptest.NewRequest("GET", "/", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != 200 {
		t.Fail()
	}

	req = httptest.NewRequest("GET", "/metrics", nil)
	resp, _ = app.Test(req)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	got := string(body)
	want := `my_service_http_requests_total{method="GET",path="/",service="test-service",status_code="200"} 1`
	if !strings.Contains(got, want) {
		t.Errorf("got %s; want %s", got, want)
	}

	want = `my_service_http_request_duration_seconds_count{method="GET",path="/",service="test-service",status_code="200"} 1`
	if !strings.Contains(got, want) {
		t.Errorf("got %s; want %s", got, want)
	}

	want = `my_service_http_requests_in_progress_total{method="GET",path="/",service="test-service"} 0`
	if !strings.Contains(got, want) {
		t.Errorf("got %s; want %s", got, want)
	}
}
