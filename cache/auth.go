package cache

import (
	"context"
	"errors"
	"fmt"

	"github.com/noo8xl/anvil-common/exceptions"
)

func (s *CacheService) Set2FACode(email, code string) {
	client, err := s.connectClient("2fa")
	if err != nil {
		exceptions.HandleAnException(err)
	}
	defer client.Close()

	if err = client.Set(context.Background(), fmt.Sprintf("2FA:%s", email), code, s.timeout).Err(); err != nil {
		exceptions.HandleAnException(err)
	}
}

func (s *CacheService) Get2FACode(email string) (string, error) {
	client, err := s.connectClient("2fa")
	if err != nil {
		return "", exceptions.HandleAnException(err)
	}
	defer client.Close()

	c, err := client.Get(context.Background(), fmt.Sprintf("2FA:%s", email)).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return "", errors.New("code not found")
		}
		return "", exceptions.HandleAnException(err)
	}

	if c == "" {
		return "", errors.New("invalid code")
	}

	return c, nil
}

func (s *CacheService) Clear2FACode(email string) {
	client, err := s.connectClient("2fa")
	if err != nil {
		exceptions.HandleAnException(err)
	}
	defer client.Close()

	if err = client.Del(context.Background(), fmt.Sprintf("2FA:%s", email)).Err(); err != nil {
		exceptions.HandleAnException(err)
	}
}
