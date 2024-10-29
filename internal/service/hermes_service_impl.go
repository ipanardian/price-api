package service

import (
	"context"
	"encoding/json"
	"fmt"
	dbg "runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/ipanardian/price-api/internal/cache"
	"github.com/ipanardian/price-api/internal/helpers"
	"github.com/ipanardian/price-api/internal/logger"
	"github.com/ipanardian/price-api/internal/model/frame"
	"github.com/ipanardian/price-api/internal/notification"
	"github.com/recws-org/recws"
	"github.com/spf13/viper"
)

type HermesServiceImpl struct {
	ws                 *recws.RecConn
	prices             map[string]*frame.PriceHermes
	pricesMx           sync.RWMutex
	pricePools         chan *frame.PriceHermes
	hermesCh           chan frame.HermesPriceFeed
	hermesMx           sync.Mutex
	hermesSubCh        chan []string
	hermesUnsubCh      chan []string
	hermesPriceIds     []string
	hermesIsSubscribed bool
	subsMx             sync.Mutex
}

func NewHermesService() HermesService {
	return &HermesServiceImpl{
		pricePools:         make(chan *frame.PriceHermes, 100),
		prices:             make(map[string]*frame.PriceHermes),
		pricesMx:           sync.RWMutex{},
		hermesCh:           make(chan frame.HermesPriceFeed),
		hermesMx:           sync.Mutex{},
		hermesSubCh:        make(chan []string),
		hermesUnsubCh:      make(chan []string),
		hermesPriceIds:     []string{},
		hermesIsSubscribed: false,
		subsMx:             sync.Mutex{},
	}
}

func (b *HermesServiceImpl) Connect() (err error) {
	ws := &recws.RecConn{
		KeepAliveTimeout: 0,
	}

	ws.Dial("wss://hermes.pyth.network/ws", nil)

	if !ws.IsConnected() {
		logger.Log.Sugar().Error("hermes connect error!")
		<-time.After(3 * time.Second)
		b.Connect()
		return
	}

	b.ws = ws
	logger.Log.Sugar().Info("Connected to Hermes")

	go func() {
		for {
			func() {
				defer func() {
					if r := recover(); r != nil {
						logger.Log.Sugar().Errorf("Hermes Pool Reader Error: %s\n%s", r, string(dbg.Stack()))
						ws.CloseAndReconnect()

						b.subsMx.Lock()
						b.hermesIsSubscribed = false
						b.subsMx.Unlock()

						time.Sleep(3 * time.Second)
					}
				}()

				if ws == nil {
					return
				}

				ws.Conn.SetReadDeadline(time.Now().Add(time.Second * 10))
				_, message, e := ws.ReadMessage()
				if e != nil {
					logger.Log.Sugar().Errorf("Hermes read error: %v", e)
					ws.CloseAndReconnect()

					b.subsMx.Lock()
					b.hermesIsSubscribed = false
					b.subsMx.Unlock()

					time.Sleep(3 * time.Second)
					return
				}

				ticker := &frame.HermesResponse{}
				e = json.Unmarshal(message, ticker)
				if e != nil {
					logger.Log.Sugar().Errorf("Hermes unmarshal error: %v", e)
					return
				}

				if strings.EqualFold(ticker.Status, "success") {
					return
				}

				if strings.EqualFold(ticker.Status, "error") {
					logger.Log.Sugar().Errorf("Hermes response error: %v", ticker.Error)
					return
				}

				if strings.EqualFold(ticker.Type, "price_update") {
					b.hermesMx.Lock()
					b.hermesCh <- ticker.PriceFeed
					b.hermesMx.Unlock()
					return
				}

				logger.Log.Sugar().Errorf("Hermes unhandled message: %s", string(message))
			}()
		}
	}()

	go func() {
		for {
			func() {
				ch := <-b.hermesSubCh
				b.hermesMx.Lock()
				defer func() {
					if r := recover(); r != nil {
						logger.Log.Sugar().Errorf("Hermes Pool Error: %s\n%s", r, string(dbg.Stack()))
					}
					b.hermesMx.Unlock()
				}()

				if len(ch) > 40 {
					logger.Log.Sugar().Error("Maximum symbols subscription is reached")
					return
				}

				if ws == nil {
					logger.Log.Sugar().Errorf("Hermes error: ws client is missing")
					return
				}

				ids := make([]string, len(ch))
				for i, item := range ch {
					ids[i] = fmt.Sprintf(`"%s"`, item)
				}

				var idStr string
				if len(ids) > 1 {
					idStr = strings.Join(ids, ",")
				} else {
					idStr = ids[0]
				}

				args := fmt.Sprintf(`{"ids":[%s],"type":"subscribe","binary":true}`, idStr)
				e := b.ws.WriteMessage(websocket.TextMessage, []byte(args))
				if e != nil {
					logger.Log.Sugar().Errorf("subscribe error: %s", err)
					b.ws.CloseAndReconnect()
					b.subsMx.Lock()
					b.hermesIsSubscribed = false
					b.subsMx.Unlock()
				}
				b.subsMx.Lock()
				b.hermesIsSubscribed = true
				b.subsMx.Unlock()
			}()
		}
	}()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Log.Sugar().Errorf("Hermes subscription watcher error: %s\n%s", r, string(dbg.Stack()))
			}
		}()

		if b.ws == nil {
			return
		}

		for {
			func() {
				defer func() {
					b.subsMx.Unlock()
				}()
				time.Sleep(10 * time.Second)

				b.subsMx.Lock()
				if b.ws.IsConnected() && !b.hermesIsSubscribed && len(b.hermesPriceIds) > 0 {
					logger.Log.Sugar().Infoln("Hermes resubscribing")
					e := b.Subscribe(b.hermesPriceIds)
					if e != nil {
						b.hermesIsSubscribed = false
						return
					}

				} else if !b.ws.IsConnected() {
					logger.Log.Sugar().Warnln("Hermes disconnected")
					b.ws.CloseAndReconnect()
					b.hermesIsSubscribed = false
				}
			}()
		}
	}()

	time.Sleep(5 * time.Second)

	return
}

func (b *HermesServiceImpl) SetPriceFeedIds(priceIdsStr string) error {
	priceIds := strings.Split(priceIdsStr, ",")
	b.hermesPriceIds = priceIds
	return nil
}

func (b *HermesServiceImpl) Subscribe(ids []string) error {
	b.hermesSubCh <- ids
	return nil
}

func (b *HermesServiceImpl) Listen() {
	go func() {
		ctx, done := context.WithCancel(context.Background())
		defer func() {
			if r := recover(); r != nil {
				logger.Log.Sugar().Errorf("Hermes Pool Error: %s\n%s", r, string(dbg.Stack()))
			}

			done()
		}()

		e := b.Subscribe(b.hermesPriceIds)
		if e != nil {
			logger.Log.Sugar().Errorf("Hermes subscribe: %v", e)
			return
		}

		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-b.hermesCh:
				go func(event frame.HermesPriceFeed) {
					defer func() {
						if r := recover(); r != nil {
							logger.Log.Sugar().Errorf("Hermes Pool Socket Error: %s\n%s", r, string(dbg.Stack()))
						}
					}()

					priceFeed, e := event.GetPriceNoOlderThan(60)
					if e != nil {
						return
					}

					prc := &frame.PriceHermes{
						ID:          event.ID,
						Price:       priceFeed.Price,
						Expo:        priceFeed.Expo,
						Conf:        priceFeed.Conf,
						PublishTime: priceFeed.PublishTime,
					}
					b.Send(prc)

				}(msg)
			}
		}
	}()
}

func (b *HermesServiceImpl) Send(prc *frame.PriceHermes) {
	b.pricePools <- prc
}

func (b *HermesServiceImpl) Sync() {
	go func() {
		func() {
			defer func() {
				if r := recover(); r != nil {
					logger.Log.Sugar().Errorln("syncPrice", string(dbg.Stack()))
				}
			}()

			for prc := range b.pricePools {
				b.pricesMx.Lock()
				b.prices[prc.ID] = prc
				b.pricesMx.Unlock()
			}
		}()
	}()

	//Push to redis
	go func() {
		for {
			func() {
				ctx, done := context.WithCancel(context.Background())
				defer func() {
					if r := recover(); r != nil {
						logger.Log.Sugar().Errorln("pushToRedis", string(dbg.Stack()))
					}
					done()
				}()
				pipeline := cache.Client().Pipeline()
				b.pricesMx.RLock()
				for i, prc := range b.prices {
					jsonStr, _ := json.Marshal(prc)
					pipeline.Set(ctx, fmt.Sprintf("price:%s", i), jsonStr, 1*time.Minute)
				}
				b.pricesMx.RUnlock()
				_, e := pipeline.Exec(ctx)
				if e != nil {
					logger.Log.Sugar().Errorln("redis exec price", e)
				}
			}()
			time.Sleep(1000 * time.Millisecond)
		}
	}()
}

func (b *HermesServiceImpl) Run() {
	logger.Log.Sugar().Info("Price service started")

	priceIdsStr := viper.GetString("PRICE_FEED_IDS")
	if priceIdsStr == "" {
		logger.Log.Sugar().Error("Please provide price IDs")
		return
	}

	b.SetPriceFeedIds(priceIdsStr)
	b.Connect()
	b.Listen()
	b.Sync()
}

func (b *HermesServiceImpl) HealthCheck() {
	go func() {
		logger.Log.Sugar().Info("Price healthcheck started")

		for {
			func() {
				time.Sleep(180 * time.Second)

				priceIdsStr := viper.GetString("PRICE_FEED_IDS")
				if priceIdsStr == "" {
					return
				}

				priceIds := strings.Split(priceIdsStr, ",")
				for _, id := range priceIds {
					id = helpers.RemoveLeading0xIfExists(id)
					price, err := cache.Get[frame.PriceHermes](context.Background(), fmt.Sprintf("price:%s", id))
					if err != nil {
						notification.SendPriceAlert(id, "Price not found in redis. Please check!")
						continue
					}

					if !price.Price.IsPositive() {
						notification.SendPriceAlert(id, fmt.Sprintf("Invalid price: %s", price.Price.String()))
					}

					if helpers.IsLastUpdateExpired(price.PublishTime, 60) {
						notification.SendPriceAlert(id, fmt.Sprintf("Price not updated. Last update: %s", time.Unix(price.PublishTime, 0).Format(time.RFC822)))
					}
				}
			}()
		}
	}()
}
