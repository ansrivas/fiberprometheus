package main

import (
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
)

func main() {
	app := fiber.New()
	prometheus := fiberprometheus.New("my-service-name")
	basicAuthMiddleware := basicauth.New(basicauth.Config{
		Users: map[string]string{
			"john": "doe",
		},
	})
	app.Use(prometheus.Middleware)
	app.Get("/metrics", basicAuthMiddleware, prometheus.Handler())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World")
	})

	app.Listen(":3000")
}
