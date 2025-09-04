package cache

import (
	"context"
	"encoding/json"
	"fmt"

	blogPb "github.com/noo8xl/anvil-api/main/blog"
	"github.com/noo8xl/anvil-common/exceptions"
)

// ClearBlog -> clear a customer blog
func (s *CacheService) ClearBlog(customerId uint64) error {
	client, err := s.connectClient("blog")
	if err != nil {
		return exceptions.HandleAnException(err)
	}
	defer client.Close()

	if err = client.Del(context.Background(), fmt.Sprintf("blog:%d", customerId)).Err(); err != nil {
		return exceptions.HandleAnException(err)
	}

	return nil
}

// SetBlog -> set a customer blog ( a list of blog items)
func (s *CacheService) SetBlog(customerId uint64, dto *blogPb.Blog) error {
	client, err := s.connectClient("blog")
	if err != nil {
		return exceptions.HandleAnException(err)
	}
	defer client.Close()

	payload, err := json.Marshal(dto)
	if err != nil {
		return exceptions.HandleAnException(err)
	}
	if err = client.Set(context.Background(), fmt.Sprintf("blog:%d", customerId), payload, s.timeout).Err(); err != nil {
		return exceptions.HandleAnException(err)
	}

	return nil
}

// GetBlog -> get a customer blog ( a list of blog items)
func (s *CacheService) GetBlog(customerId uint64) (*blogPb.Blog, error) {
	client, err := s.connectClient("blog")
	if err != nil {
		return nil, exceptions.HandleAnException(err)
	}
	defer client.Close()

	result, err := client.Get(context.Background(), fmt.Sprintf("blog:%d", customerId)).Result()
	if err != nil {
		return nil, exceptions.HandleAnException(err)
	}

	var dto blogPb.Blog
	err = json.Unmarshal([]byte(result), &dto)
	if err != nil {
		return nil, exceptions.HandleAnException(err)
	}

	return &dto, nil
}
