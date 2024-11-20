**NOTE**: 🚨 We are currently migrating this middleware to the official gofiber/contrib repo, official link etc. will be posted soon.



# fiberprometheus

Prometheus middleware for [Fiber](https://github.com/gofiber/fiber)).

**Note: Requires Go 1.22 and above**

![Release](https://img.shields.io/github/release/ansrivas/fiberprometheus.svg)
[![Discord](https://img.shields.io/badge/discord-join%20channel-7289DA)](https://gofiber.io/discord)
![Test](https://github.com/ansrivas/fiberprometheus/workflows/Test/badge.svg)
![Security](https://github.com/ansrivas/fiberprometheus/workflows/Security/badge.svg)
![Linter](https://github.com/ansrivas/fiberprometheus/workflows/Linter/badge.svg)

Following metrics are available by default:

```text
http_requests_total
http_request_duration_seconds
http_requests_in_progress_total
```

### Install v2

```console
go get -u github.com/gofiber/fiber/v2
go get -u github.com/ansrivas/fiberprometheus/v2
```

### Example using v2

```go
package main

import (
  "github.com/ansrivas/fiberprometheus/v2"
  "github.com/gofiber/fiber/v2"
)

func main() {
  app := fiber.New()

  // This here will appear as a label, one can also use
  // fiberprometheus.NewWith(servicename, namespace, subsystem )
  // or
  // labels := map[string]string{"custom_label1":"custom_value1", "custom_label2":"custom_value2"}
  // fiberprometheus.NewWithLabels(labels, namespace, subsystem )
  prometheus := fiberprometheus.New("my-service-name")
  prometheus.RegisterAt(app, "/metrics")
  prometheus.SetSkipPaths([]string{"/ping"}) // Optional: Remove some paths from metrics
  app.Use(prometheus.Middleware)

  app.Get("/", func(c *fiber.Ctx) error {
    return c.SendString("Hello World")
  })

  app.Get("/ping", func(c *fiber.Ctx) error {
    return c.SendString("pong")
  })

  app.Post("/some", func(c *fiber.Ctx) error {
    return c.SendString("Welcome!")
  })

  app.Listen(":3000")
}
```

### Result

- Hit the default url at http://localhost:3000
- Navigate to http://localhost:3000/metrics

### Grafana Dashboard

- https://grafana.com/grafana/dashboards/14331
