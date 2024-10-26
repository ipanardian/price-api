package logger

import (
	"context"
	"os"
	"sync"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ctxKey struct{}

var once sync.Once
var Log *zap.Logger

func init() {
	Log = Get()
}

func Get() *zap.Logger {
	once.Do(func() {
		stdout := zapcore.AddSync(os.Stdout)

		logLevel := int(zap.DebugLevel)
		if viper.GetString("LOG_LEVEL") == "debug" {
			logLevel = int(zap.DebugLevel)
		}

		level := zap.NewAtomicLevelAt(zapcore.Level(logLevel))
		layoutTime := "2006/01/20 15:04:05"

		productionCfg := zap.NewProductionEncoderConfig()
		productionCfg.EncodeTime = zapcore.TimeEncoderOfLayout(layoutTime)
		productionCfg.EncodeLevel = zapcore.CapitalLevelEncoder

		developmentCfg := zap.NewDevelopmentEncoderConfig()
		developmentCfg.EncodeTime = zapcore.TimeEncoderOfLayout(layoutTime)
		developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

		consoleEncoderDev := zapcore.NewConsoleEncoder(developmentCfg)
		consoleEncoderProd := zapcore.NewConsoleEncoder(productionCfg)

		core := zapcore.NewCore(consoleEncoderDev, stdout, level)
		if viper.GetString("ENV") != "development" {
			core = zapcore.NewTee(
				zapcore.NewCore(consoleEncoderProd, stdout, level),
			)
		}

		Log = zap.New(core, zap.AddCaller())
	})

	return Log
}

func FromCtx(ctx context.Context) *zap.Logger {
	if log, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok {
		return log
	} else if log := Log; log != nil {
		return log
	}

	return zap.NewNop()
}

func WithContext(ctx context.Context, l *zap.Logger) context.Context {
	if lp, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok {
		if lp == l {
			// Do not store same logger.
			return ctx
		}
	}

	return context.WithValue(ctx, ctxKey{}, l)
}
