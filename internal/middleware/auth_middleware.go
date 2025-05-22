package middleware

import (
	"encoding/json"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ipanardian/price-api/internal/cache"
	"github.com/ipanardian/price-api/internal/constant"
	dtoV1 "github.com/ipanardian/price-api/internal/dto/v1"
	"github.com/ipanardian/price-api/internal/helpers"
	"github.com/ipanardian/price-api/internal/logger"
	"github.com/ipanardian/price-api/internal/model/frame"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// AuthApiMiddleware checks if the request has a valid API key, if not, it returns a 401 status code.
// The API key is encrypted using AES-256-CBC and stored in a file named "key.json".
// The decrypted API key is cached in memory for 1 hour to reduce the number of decryption operations.
func AuthApiMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		apiKey := c.Get("X-API-KEY")
		_aesKey := viper.GetString("AES_KEY")

		if apiKey == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(dtoV1.ResponseWrapper{
				Status:        constant.RequestFailure,
				StatusCode:    "401",
				StatusNumber:  constant.StatusUnauthorized,
				StatusMessage: "Missing API Key",
				Timestamp:     time.Now().Unix(),
			})
		}

		var apiKeyList []frame.APIKey
		apiKeyListMem, f := cache.MemGet("api_keys")
		if !f {
			apiKeys, err := os.ReadFile("key/key.json")
			if err != nil {
				logger.Log.Error("Read file error: ", zap.Error(err))
				return returnStatusInternalServerError(c)
			}

			if err := json.Unmarshal(apiKeys, &apiKeyList); err != nil {
				logger.Log.Error("Unmarshal error: ", zap.Error(err))
				return returnStatusInternalServerError(c)
			}

			apiKeysDecrypted := make([]frame.APIKey, len(apiKeyList))
			for i, v := range apiKeyList {
				key, err := helpers.DecryptAES256CBC(_aesKey, v.Key)
				if err != nil {
					logger.Log.Error("Decrypt error: ", zap.Error(err), zap.String("key", v.Key))
					return returnStatusInternalServerError(c)
				}

				apiKeysDecrypted[i] = frame.APIKey{
					Name: v.Name,
					Key:  key,
				}
			}

			apiKeyList = apiKeysDecrypted
			cache.MemSet("api_keys", apiKeyList, 60*time.Minute)
		} else {
			apiKeyList = apiKeyListMem.([]frame.APIKey)
		}

		for _, v := range apiKeyList {
			if v.Key == apiKey {
				c.Locals("x-api-key", v.Key)
				return c.Next()
			}
		}

		return c.Status(fiber.StatusUnauthorized).JSON(dtoV1.ResponseWrapper{
			Status:        constant.RequestFailure,
			StatusCode:    "401",
			StatusNumber:  constant.StatusUnauthorized,
			StatusMessage: "Invalid API Key",
			Timestamp:     time.Now().Unix(),
		})
	}
}

func returnStatusInternalServerError(c *fiber.Ctx) error {
	return c.Status(fiber.StatusInternalServerError).JSON(dtoV1.ResponseWrapper{
		Status:        constant.RequestFailure,
		StatusCode:    "500",
		StatusNumber:  constant.InternalError,
		StatusMessage: "Sorry, this is not working properly. We know about this mistake and are working to fix it.",
		Timestamp:     time.Now().Unix(),
	})
}
