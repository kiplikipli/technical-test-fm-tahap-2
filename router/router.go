package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kiplikipli/technical-test-fm-tahap-2/middleware"
)

func Initalize(router *fiber.App) {
	router.Use(middleware.Security)

	router.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).SendString("Hello, World!")
	})

	router.Use(middleware.Json)
}
