package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/ipanardian/price-api/internal/constant"
	dtoV1 "github.com/ipanardian/price-api/internal/dto/v1"
)

func InitApiMiddleware(app *fiber.App) {
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	app.Use(logger.New(logger.Config{
		Format:     "${time} ${status} - ${method} ${path} ${latency}\n",
		TimeFormat: "02-Jan-2006 15:04:05",
		TimeZone:   "UTC",
	}))

	app.Use(cors.New())

	app.Use(limiter.New(limiter.Config{
		Max:        constant.RateLimitMaxRequest,
		Expiration: constant.RateLimitExpiration,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(dtoV1.ResponseWrapper{
				Status:        constant.RequestFailure,
				StatusCode:    fmt.Sprintf("%d", fiber.StatusTooManyRequests),
				StatusNumber:  constant.TooManyRequest,
				StatusMessage: "Too many requests",
				Timestamp:     time.Now().Unix(),
			})
		},
	}))
}
