package cache

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ipanardian/price-api/internal/logger"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type CacheConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

const DefaultTimeout = 10000 * time.Millisecond

var config CacheConfig
var redisClient *redis.Client

func InitRedis(conf CacheConfig) {
	redisClient = redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("%s:%s", conf.Host, conf.Port),
		Password:        conf.Password, // no password set
		DB:              conf.DB,       // use default DB
		MaxRetries:      3,
		MaxRetryBackoff: 500 * time.Millisecond,
		MinRetryBackoff: 1 * time.Second,
		DialTimeout:     1 * time.Second,
		ReadTimeout:     60 * time.Minute,
		WriteTimeout:    60 * time.Minute,
		PoolSize:        1000,
		MinIdleConns:    30,
		MaxIdleConns:    1000,
		MaxActiveConns:  10000,
		ConnMaxIdleTime: 3 * time.Second,
	})

	config = conf

	go func() {
		for {
			func() {
				ctx, done := context.WithCancel(context.Background())
				defer done()
				_, err := redisClient.WithTimeout(DefaultTimeout).Ping(ctx).Result()
				if err != nil {
					logger.Log.Error("redis ping error", zap.Error(err))
					time.Sleep(1 * time.Minute)
					return
				}

				hostname, _ := os.Hostname()

				e := SetWithExpiration(ctx, fmt.Sprintf("connection:%s:%s", "price-api", hostname), "connected", 1*time.Minute)
				if e != nil {
					logger.Log.Error("redis set ping error", zap.Error(err))
				}
			}()
			time.Sleep(1 * time.Minute)
		}
	}()
}

func Client() *redis.Client {
	return redisClient
}

func Config() CacheConfig {
	return config
}

func Get[T any](ctx context.Context, key string) (T, error) {
	result := new(T)
	e := redisClient.WithTimeout(DefaultTimeout).Get(ctx, key).Scan(result)
	return *result, e
}

func GetDel[T any](ctx context.Context, key string) (T, error) {
	result := new(T)
	e := redisClient.WithTimeout(DefaultTimeout).GetDel(ctx, key).Scan(result)
	return *result, e
}

func Gets[T any](ctx context.Context, pattern string) ([]T, error) {
	result := []T{}
	iter := redisClient.WithTimeout(DefaultTimeout).Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		p, e := Get[T](ctx, iter.Val())
		if e != nil {
			return nil, e
		}
		result = append(result, p)
	}

	if e := iter.Err(); e != nil {
		return nil, e
	}

	return result, nil
}

func Set(ctx context.Context, key string, value any) error {
	return redisClient.WithTimeout(DefaultTimeout).Set(ctx, key, value, 0).Err()
}

func SetWithExpiration(ctx context.Context, key string, value any, expiration time.Duration) error {
	return redisClient.WithTimeout(DefaultTimeout).Set(ctx, key, value, expiration).Err()
}

func MSet(ctx context.Context, collection any) error {
	return redisClient.WithTimeout(DefaultTimeout).MSet(ctx, collection).Err()
}

func MSetWithExpiration(ctx context.Context, collection map[string]any, expiration time.Duration) error {
	return redisClient.WithTimeout(DefaultTimeout).MSet(ctx, collection, expiration).Err()
}

func Exists(ctx context.Context, key string) bool {
	return redisClient.WithTimeout(DefaultTimeout).Exists(ctx, key).Err() == nil
}

func Keys(ctx context.Context, pattern string) ([]string, error) {
	return redisClient.WithTimeout(DefaultTimeout).Keys(ctx, pattern).Result()
}

func LPush(ctx context.Context, key string, value any) error {
	return redisClient.WithTimeout(DefaultTimeout).LPush(ctx, key, value).Err()
}

func RPush(ctx context.Context, key string, value any) error {
	return redisClient.WithTimeout(DefaultTimeout).RPush(ctx, key, value).Err()
}

func LPop[T any](ctx context.Context, key string) (T, error) {
	result := new(T)
	e := redisClient.WithTimeout(DefaultTimeout).LPop(ctx, key).Scan(result)
	return *result, e
}

func RPop[T any](ctx context.Context, key string) (T, error) {
	result := new(T)
	e := redisClient.WithTimeout(DefaultTimeout).RPop(ctx, key).Scan(result)
	return *result, e
}

func Remove(ctx context.Context, key string) error {
	return redisClient.WithTimeout(DefaultTimeout).Del(ctx, key).Err()
}

func Flush(ctx context.Context) error {
	return redisClient.WithTimeout(DefaultTimeout).FlushDB(ctx).Err()
}
