package cache

import (
	"context"
	"errors"
	"time"

	"github.com/noo8xl/anvil-common/exceptions"
	"github.com/noo8xl/anvil-gateway/config"
	"github.com/redis/go-redis/v9"
)

type CacheService struct {
	timeout time.Duration
}

func InitCacheService() *CacheService {
	return &CacheService{
		timeout: time.Millisecond * 300000, // 5 minutes
	}
}

func (c *CacheService) connectClient(svcName string) (*redis.Client, error) {

	var client *redis.Client
	var opts *redis.Options
	var err error

	switch svcName {
	case "offers":
		opts = config.GetOffersRedisConfig()
	case "profile":
		opts = config.GetProfileRedisConfig()
	case "orders":
		opts = config.GetOrdersRedisConfig()
	case "reviews":
		opts = config.GetReviewsRedisConfig()
	case "blog":
		opts = config.GetBlogRedisConfig()
	case "promo":
		opts = config.GetPromoRedisConfig()
	case "2fa":
		opts = config.Get2FARedisConfig()
	case "notifications":
		opts = config.GetNotificationsRedisConfig()
	// case "payments":
	// 	opts = config.GetPaymentsRedisConfig()
	default:
		return nil, exceptions.HandleAnException(errors.New("gateway: invalid service name to get redis config"))
	}

	client = redis.NewClient(opts)
	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		return nil, exceptions.HandleAnException(err)
	}
	return client, nil
}
