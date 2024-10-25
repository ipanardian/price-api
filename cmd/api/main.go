// @version 1.0
// @securityDefinitions.apikey ClientIdAuth
// @in header
// @name X-Client-ID
// @securityDefinitions.apikey ClientSignatureAuth
// @in header
// @name Authorization
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ipanardian/price-api/internal/cache"
	"github.com/ipanardian/price-api/internal/config"
	"github.com/ipanardian/price-api/internal/constant"
	"github.com/ipanardian/price-api/internal/handler"
	"github.com/ipanardian/price-api/internal/injector"
	"github.com/ipanardian/price-api/internal/logger"
	"github.com/ipanardian/price-api/internal/middleware"
	"github.com/ipanardian/price-api/internal/routes"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func init() {
	cache.InitRedis(cache.CacheConfig{
		Host:     viper.GetString("REDIS_HOST"),
		Port:     viper.GetString("REDIS_PORT"),
		Password: viper.GetString("REDIS_PASS"),
		DB:       viper.GetInt("REDIS_DB"),
	})

	cache.InitMemory(cache.MemoryCacheConfig{
		DefaultExpiration: time.Minute * 1,
		CleanupInterval:   time.Minute * 10,
	})
}

func main() {
	app := fiber.New(config.FiberConfig())
	port := viper.GetString(constant.REST_PORT)

	middleware.InitApiMiddleware(app)

	apiService, err := injector.InitializeApiService()
	if err != nil {
		logger.Log.Error("Init api service error: ", zap.Error(err))
		return
	}

	handler := handler.NewHandler(apiService)
	routes.SetupApiRouter(app, handler)

	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}
