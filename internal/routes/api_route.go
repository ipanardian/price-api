package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	_ "github.com/ipanardian/price-api/docs"
	dtoV1 "github.com/ipanardian/price-api/internal/dto/v1"
	"github.com/ipanardian/price-api/internal/handler"
)

func SetupApiRouter(app *fiber.App, handler *handler.Handler) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Price API Service v0.1.0")
	})

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("Pong!")
	})

	app.Get("/docs/*", swagger.New(swagger.Config{
		Title: "Price API Documentation",
	}))

	api := app.Group("v1")
	api.Get("/price", handler.Api.GetPrice)

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(dtoV1.ResponseWrapper{
			Status:        0,
			StatusCode:    "404",
			StatusNumber:  "404",
			StatusMessage: "404 Not Found",
			Timestamp:     time.Now().Unix(),
		})
	})
}
