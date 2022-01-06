package main

import (
	"log"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	prometheus := fiberprometheus.New("my-service-name")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	app.Get("/", func(c *fiber.Ctx) error {
		// 503 Service Unavailable
		return fiber.ErrServiceUnavailable
	})

	app.Get("/not-ok", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusServiceUnavailable).SendString("Unavailable")
	})

	app.Get("/ok", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("available")
	})

	log.Fatal(app.Listen(":3000"))

}
