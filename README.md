# fiberprometheus

Prometheus middleware for gofiber.

**Note: Requires Go 1.21 and above**

![Release](https://img.shields.io/github/release/ansrivas/fiberprometheus.svg)
[![Discord](https://img.shields.io/badge/discord-join%20channel-7289DA)](https://gofiber.io/discord)
![Test](https://github.com/ansrivas/fiberprometheus/workflows/Test/badge.svg)
![Security](https://github.com/ansrivas/fiberprometheus/workflows/Security/badge.svg)
![Linter](https://github.com/ansrivas/fiberprometheus/workflows/Linter/badge.svg)

Following metrics are available by default:

```
http_requests_total
http_request_duration_seconds
http_requests_in_progress_total
http_cache_results
```

### Install v2

```
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
  app.Use(prometheus.Middleware)

  app.Get("/", func(c *fiber.Ctx) error {
    return c.SendString("Hello World")
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

### Grafana Board

- https://grafana.com/grafana/dashboards/14331
