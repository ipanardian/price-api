package constant

import "time"

const (
	REST_PORT     = "REST_PORT"
	AES_KEY       = "AES_KEY"
	HMAC_KEY      = "HMAC_KEY"
	DISCORD_TOKEN = "DISCORD_TOKEN"
)

const (
	ApplicationCache = "application.%s"
	PriceCache       = "price.%s"
)

const (
	RateLimitMaxRequest = "RATE_LIMIT_MAX_REQUEST"
	RateLimitExpiration = "RATE_LIMIT_EXPIRATION"
)

const (
	RedisPriceCacheTTL       = 5 * time.Minute
	RedisPriceUpdateInterval = 500 * time.Millisecond
)
