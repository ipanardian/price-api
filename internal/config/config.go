package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	dtoV1 "github.com/ipanardian/price-api/internal/dto/v1"
	"github.com/ipanardian/price-api/internal/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func init() {
	if !strings.HasSuffix(os.Args[0], ".test") {
		viper.SetConfigFile(".env")
		viper.AddConfigPath(".")
		err := viper.ReadInConfig()
		if err != nil {
			logger.Log.Error("Failed to read environment", zap.Error(err))
			return
		}
	} else {
		viper.SetConfigFile("../.test.env")
		viper.AddConfigPath(".")
		_ = viper.ReadInConfig()
	}

	readRuntimeEnv()
}

func readRuntimeEnv() {
	viper.SetDefault("STATUS_SERVICE", 92)
}

func FiberConfig() fiber.Config {
	return fiber.Config{
		// Override default error handler
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Status code defaults to 500
			code := fiber.StatusInternalServerError

			// Retrieve the custom status code if it's a *fiber.Error
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			logger.Log.Error("Internal Server Error", zap.Error(err))

			// Send custom error page
			err = c.Status(code).JSON(dtoV1.ResponseWrapper{
				Status:        0,
				StatusCode:    fmt.Sprintf("%d", code),
				StatusNumber:  "50000",
				StatusMessage: "Sorry, this is not working properly. We know about this mistake and are working to fix it.",
				Timestamp:     time.Now().Unix(),
			})
			if err != nil {
				// In case the SendFile fails
				return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
			}

			// Return from handler
			return nil
		},
	}
}
