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
	RateLimitMaxRequest = 50
	RateLimitExpiration = 1 * time.Second
)
