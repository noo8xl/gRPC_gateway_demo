package config

import (
	"os"

	"github.com/redis/go-redis/v9"
)

// ################## -> setup
// https://help.ivanti.com/ht/help/en_US/ISM/2023/InstallAndDeploy/Content/Install_Deploy_guide/On-premise-Redis-Setup.htm

func GetProfileRedisConfig() *redis.Options {
	var opts *redis.Options
	env := os.Getenv("GO_ENV")

	switch env {
	case "development":
		opts = &redis.Options{
			Addr: "127.0.0.1" + ":" + "6379",
			DB:   0,
		}
	case "production":
		opts = &redis.Options{
			Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
			Username: os.Getenv("REDIS_USERNAME"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       0,
		}
	case "test":
		opts = &redis.Options{
			Addr: "127.0.0.1" + ":" + "6379",
			DB:   0,
		}
	default:
		return nil
	}

	return opts
}

func GetBlogRedisConfig() *redis.Options {
	var opts *redis.Options
	env := os.Getenv("GO_ENV")

	switch env {
	case "development":
		opts = &redis.Options{
			Addr: "127.0.0.1" + ":" + "6379",
			DB:   1,
		}
	case "production":
		opts = &redis.Options{
			Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
			Username: os.Getenv("REDIS_USERNAME"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       1,
		}
	case "test":
		opts = &redis.Options{
			Addr: "127.0.0.1" + ":" + "6379",
			DB:   1,
		}
	default:
		return nil
	}

	return opts
}

func GetOffersRedisConfig() *redis.Options {
	var opts *redis.Options
	env := os.Getenv("GO_ENV")

	switch env {
	case "development":
		opts = &redis.Options{
			Addr: "127.0.0.1" + ":" + "6379",
			DB:   2,
		}
	case "production":
		opts = &redis.Options{
			Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
			Username: os.Getenv("REDIS_USERNAME"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       2,
		}
	case "test":
		opts = &redis.Options{
			Addr: "127.0.0.1" + ":" + "6379",
			DB:   2,
		}
	default:
		return nil
	}

	return opts
}

func GetOrdersRedisConfig() *redis.Options {
	var opts *redis.Options
	env := os.Getenv("GO_ENV")

	switch env {
	case "development":
		opts = &redis.Options{
			Addr: "127.0.0.1" + ":" + "6379",
			DB:   3,
		}
	case "production":
		opts = &redis.Options{
			Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
			Username: os.Getenv("REDIS_USERNAME"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       3,
		}
	case "test":
		opts = &redis.Options{
			Addr: "127.0.0.1" + ":" + "6379",
			DB:   3,
		}
	default:
		return nil
	}

	return opts
}

func GetReviewsRedisConfig() *redis.Options {
	var opts *redis.Options
	env := os.Getenv("GO_ENV")

	switch env {
	case "development":
		opts = &redis.Options{
			Addr: "127.0.0.1" + ":" + "6379",
			DB:   4,
		}
	case "test":
		opts = &redis.Options{
			Addr: "127.0.0.1" + ":" + "6379",
			DB:   4,
		}
	case "production":
		opts = &redis.Options{
			Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
			Username: os.Getenv("REDIS_USERNAME"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       4,
		}
	default:
		return nil
	}

	return opts
}

func GetPromoRedisConfig() *redis.Options {
	var opts *redis.Options
	env := os.Getenv("GO_ENV")

	switch env {
	case "development":
		opts = &redis.Options{
			Addr: "127.0.0.1" + ":" + "6379",
			DB:   5,
		}
	case "test":
		opts = &redis.Options{
			Addr: "127.0.0.1" + ":" + "6379",
			DB:   5,
		}
	case "production":
		opts = &redis.Options{
			Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
			Username: os.Getenv("REDIS_USERNAME"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       5,
		}
	default:
		return nil
	}

	return opts
}

func Get2FARedisConfig() *redis.Options {

	var opts *redis.Options
	env := os.Getenv("GO_ENV")

	switch env {
	case "development":
		opts = &redis.Options{
			Addr: "127.0.0.1" + ":" + "6379",
			DB:   6,
		}
	case "test":
		opts = &redis.Options{
			Addr: "127.0.0.1" + ":" + "6379",
			DB:   6,
		}
	case "production":
		opts = &redis.Options{
			Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
			Username: os.Getenv("REDIS_USERNAME"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       6,
		}
	}

	return opts
}

func GetNotificationsRedisConfig() *redis.Options {
	var opts *redis.Options
	env := os.Getenv("GO_ENV")

	switch env {
	case "development":
		opts = &redis.Options{
			Addr: "127.0.0.1" + ":" + "6379",
			DB:   7,
		}
	case "production":
		opts = &redis.Options{
			Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
			Username: os.Getenv("REDIS_USERNAME"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       7,
		}
	case "test":
		opts = &redis.Options{
			Addr: "127.0.0.1" + ":" + "6379",
			DB:   7,
		}
	default:
		return nil
	}

	return opts
}
