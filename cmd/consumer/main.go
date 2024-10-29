package main

import (
	"time"

	"github.com/ipanardian/price-api/internal/cache"
	_ "github.com/ipanardian/price-api/internal/config"
	"github.com/ipanardian/price-api/internal/injector"
	"github.com/ipanardian/price-api/internal/logger"
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
	hermesService, err := injector.InitializeHermesService()
	if err != nil {
		logger.Log.Error("Init hermes service error: ", zap.Error(err))
		return
	}

	hermesService.Run()
	hermesService.HealthCheck()

	select {}
}
