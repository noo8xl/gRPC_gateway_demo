package config_test

import (
	"testing"

	"github.com/noo8xl/anvil-gateway/config"
)

func TestGetProfileRedisConfig(t *testing.T) {
	redisConfig := config.GetProfileRedisConfig()

	if redisConfig == nil {
		t.Errorf("error: redisConfig is nil, expected not nil options struct")
	} else {
		if redisConfig.DB != 0 {
			t.Errorf("error: redisConfig.DB is %d, expected 0", redisConfig.DB)
		}
	}
}

func TestGetBlogRedisConfig(t *testing.T) {
	redisConfig := config.GetBlogRedisConfig()

	if redisConfig == nil {
		t.Errorf("error: redisConfig is nil, expected not nil options struct")
	} else {
		if redisConfig.DB != 1 {
			t.Errorf("error: redisConfig.DB is %d, expected 1", redisConfig.DB)
		}
	}
}

func TestGetOffersRedisConfig(t *testing.T) {
	redisConfig := config.GetOffersRedisConfig()

	if redisConfig == nil {
		t.Errorf("error: redisConfig is nil, expected not nil options struct")
	} else {
		if redisConfig.DB != 2 {
			t.Errorf("error: redisConfig.DB is %d, expected 2", redisConfig.DB)
		}
	}
}

func TestGetOrdersRedisConfig(t *testing.T) {
	redisConfig := config.GetOrdersRedisConfig()

	if redisConfig == nil {
		t.Errorf("error: redisConfig is nil, expected not nil options struct")
	} else {
		if redisConfig.DB != 3 {
			t.Errorf("error: redisConfig.DB is %d, expected 3", redisConfig.DB)
		}
	}
}

func TestGetReviewsRedisConfig(t *testing.T) {
	redisConfig := config.GetReviewsRedisConfig()

	if redisConfig == nil {
		t.Errorf("error: redisConfig is nil, expected not nil options struct")
	} else {
		if redisConfig.DB != 4 {
			t.Errorf("error: redisConfig.DB is %d, expected 4", redisConfig.DB)
		}
	}
}

func TestGetPromoRedisConfig(t *testing.T) {
	redisConfig := config.GetPromoRedisConfig()

	if redisConfig == nil {
		t.Errorf("error: redisConfig is nil, expected not nil options struct")
	} else {
		if redisConfig.DB != 5 {
			t.Errorf("error: redisConfig.DB is %d, expected 5", redisConfig.DB)
		}
	}
}

func TestGet2FARedisConfig(t *testing.T) {
	redisConfig := config.Get2FARedisConfig()

	if redisConfig == nil {
		t.Errorf("error: redisConfig is nil, expected not nil options struct")
	} else {
		if redisConfig.DB != 6 {
			t.Errorf("error: redisConfig.DB is %d, expected 6", redisConfig.DB)
		}
	}
}

func TestGetNotificationsRedisConfig(t *testing.T) {
	redisConfig := config.GetNotificationsRedisConfig()

	if redisConfig == nil {
		t.Errorf("error: redisConfig is nil, expected not nil options struct")
	} else {
		if redisConfig.DB != 7 {
			t.Errorf("error: redisConfig.DB is %d, expected 7", redisConfig.DB)
		}
	}
}
