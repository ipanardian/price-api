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
	"github.com/spf13/viper"
)

func InitApiMiddleware(app *fiber.App) {
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${queryParams} | ${userAgent} | ${error}\n",
		TimeFormat: "02-Jan-2006 15:04:05",
		TimeZone:   "UTC",
		CustomTags: map[string]logger.LogFunc{
			"userAgent": func(output logger.Buffer, c *fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
				reqHeaders := c.GetReqHeaders()
				return output.WriteString(reqHeaders["User-Agent"][0])
			},
		},
	}))

	app.Use(cors.New())

	app.Use(limiter.New(limiter.Config{
		Max:        viper.GetInt(constant.RateLimitMaxRequest),
		Expiration: viper.GetDuration(constant.RateLimitExpiration),
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
