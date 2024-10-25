package config

import (
	"os"
	"strings"

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
