package cache_test

import (
	"github.com/noo8xl/anvil-gateway/cache"
)

var (
	svc *cache.CacheService
)

func init() {
	svc = cache.InitCacheService()
}
