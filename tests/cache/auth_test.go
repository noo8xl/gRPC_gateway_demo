package cache_test

import (
	"testing"

	"github.com/noo8xl/anvil-gateway/cache"
)

var (
	email = "test@test.com"
	code  = "123456"
)

func init() {
	svc = cache.InitCacheService()
}

func TestSet2FACode(t *testing.T) {
	svc.Set2FACode(email, code)
}

func TestGet2FACode(t *testing.T) {
	c, err := svc.Get2FACode(email)
	if err != nil {
		t.Errorf("Error getting 2FA code: %v", err)
	}
	if c != code {
		t.Errorf("2FA code is not correct: %v", code)
	}
}

func TestClear2FACode(t *testing.T) {
	svc.Clear2FACode(email)
}
