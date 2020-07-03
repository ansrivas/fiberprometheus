# FiberPrometheus

![Release](https://img.shields.io/github/release/gofiber/csrf.svg)
[![Discord](https://img.shields.io/badge/discord-join%20channel-7289DA)](https://gofiber.io/discord)
![Test](https://github.com/gofiber/csrf/workflows/Test/badge.svg)
![Security](https://github.com/gofiber/csrf/workflows/Security/badge.svg)
![Linter](https://github.com/gofiber/csrf/workflows/Linter/badge.svg)

### Install
```
go get -u github.com/gofiber/fiber
go get -u github.com/ansrivas/fiberprometheus
```
### Example
```go
package main

import (
  "github.com/gofiber/fiber"
  "github.com/ansrivas/fiberprometheus"
)

func main() {
  app := fiber.New()

  prometheus := fiberprometheus.New("my-service-name")
  prometheus.RegisterAt(app, "/metrics")
  app.Use(prometheus.Middleware)
  
  app.Get("/", func(c *fiber.Ctx) {
    c.Send("Hello World")
  })

  app.Post("/register", func(c *fiber.Ctx) {
    c.Send("Welcome!")
  })

  app.Listen(3000)
}
```
