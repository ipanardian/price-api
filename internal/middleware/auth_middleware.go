package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ipanardian/price-api/internal/cache"
	"github.com/ipanardian/price-api/internal/constant"
	dtoV1 "github.com/ipanardian/price-api/internal/dto/v1"
	"github.com/ipanardian/price-api/internal/helpers"
	"github.com/spf13/viper"
)

func AuthApiMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		apiKey := c.Get("X-API-KEY")
		apiSignature := c.Get("X-API-SIGNATURE")

		if apiKey == "" || apiSignature == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(dtoV1.ResponseWrapper{
				Status:        constant.RequestFailure,
				StatusCode:    "401",
				StatusNumber:  constant.StatusUnauthorized,
				StatusMessage: "Missing API Key or Signature",
				Timestamp:     time.Now().Unix(),
			})
		}

		_aesKey := viper.GetString("AES_KEY")
		_apiKey := viper.GetString("API_KEY")
		_apiSecretEn := viper.GetString("API_SECRET")

		if apiKey != _apiKey {
			return c.Status(fiber.StatusUnauthorized).JSON(dtoV1.ResponseWrapper{
				Status:        constant.RequestFailure,
				StatusCode:    "401",
				StatusNumber:  constant.StatusUnauthorized,
				StatusMessage: "Invalid API Key",
				Timestamp:     time.Now().Unix(),
			})
		}

		var _apiSecretDe string
		_apiSecretDeMem, f := cache.MemGet(fmt.Sprintf("scrt.%s", apiKey))
		if !f {
			var err error
			_apiSecretDe, err = helpers.DecryptAES256CBC(_aesKey, _apiSecretEn)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(dtoV1.ResponseWrapper{
					Status:        constant.RequestFailure,
					StatusCode:    "500",
					StatusNumber:  constant.DecryptAESError,
					StatusMessage: "Sorry, this is not working properly. We know about this mistake and are working to fix it.",
					Timestamp:     time.Now().Unix(),
				})
			}

			cache.MemSet(fmt.Sprintf("scrt.%s", apiKey), _apiSecretDe, 60*time.Minute)
		} else {
			_apiSecretDe = _apiSecretDeMem.(string)
		}

		signature := helpers.GenerateSignature(_apiSecretDe)
		if apiSignature != signature {
			return c.Status(fiber.StatusUnauthorized).JSON(dtoV1.ResponseWrapper{
				Status:        constant.RequestFailure,
				StatusCode:    "401",
				StatusNumber:  constant.StatusUnauthorized,
				StatusMessage: "Invalid Signature",
				Timestamp:     time.Now().Unix(),
			})
		}
		c.Locals("x-api-key", apiKey)

		return c.Next()
	}
}
