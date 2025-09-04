package cache_test

import (
	"testing"
	"time"

	blogPb "github.com/noo8xl/anvil-api/main/blog"
)

var (
	blog = &blogPb.Blog{
		Blog: []*blogPb.BlogItem{
			{
				Id:          1,
				Title:       "Test Post",
				Description: "Test Description",
				Likes:       1,
				Dislikes:    4,
				CreatedAt:   time.Now().Format(time.RFC3339),
				UpdatedAt:   time.Now().Format(time.RFC3339),
				Tags:        []string{"test", "test2"},
				ImageList:   []string{"test.com", "test2.com"},
			},
			{
				Id:          2,
				Title:       "Test Post 2",
				Description: "Test Description 2",
				Likes:       3,
				Dislikes:    2,
				CreatedAt:   time.Now().Format(time.RFC3339),
				UpdatedAt:   time.Now().Format(time.RFC3339),
				Tags:        []string{"test", "test2"},
				ImageList:   []string{"test.com", "test2.com"},
			},
		},
	}
)

func TestSetBlog(t *testing.T) {
	err := svc.SetBlog(customerId, blog)
	if err != nil {
		t.Errorf("error setting blog: Expected nil, got %v", err)
	}
}

func TestGetBlog(t *testing.T) {
	blog, err := svc.GetBlog(customerId)
	if err != nil {
		t.Errorf("error getting blog: Expected nil, got %v", err)
	}
	if blog != nil {
		if len(blog.Blog) == 0 {
			t.Errorf("error getting blog: blog is empty")
		}
	}
}

func TestClearBlog(t *testing.T) {
	err := svc.ClearBlog(customerId)
	if err != nil {

		t.Errorf("error cleaning blog: Expected nil, got %v", err)
	}
}
