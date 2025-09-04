package cache

import (
	"context"
	"encoding/json"
	"fmt"

	pb "github.com/noo8xl/anvil-api/main/reviews"
	"github.com/noo8xl/anvil-common/exceptions"
)

func (s *CacheService) ClearReviewCommentsList(reviewId uint64) error {
	client, err := s.connectClient("reviews")
	if err != nil {
		return exceptions.HandleAnException(err)
	}
	defer client.Close()

	if err = client.Del(context.Background(), fmt.Sprintf("reviews_comments:%d", reviewId)).Err(); err != nil {
		return exceptions.HandleAnException(err)
	}

	return nil
}

func (s *CacheService) SetReviewCommentsList(reviewId uint64, dto *pb.ReviewCommentsList) error {
	client, err := s.connectClient("reviews")
	if err != nil {
		return exceptions.HandleAnException(err)
	}
	defer client.Close()

	payload, err := json.Marshal(dto)
	if err != nil {
		return exceptions.HandleAnException(err)
	}

	if err = client.Set(context.Background(), fmt.Sprintf("reviews_comments:%d", reviewId), payload, s.timeout).Err(); err != nil {
		return exceptions.HandleAnException(err)
	}

	return nil
}

func (s *CacheService) GetReviewCommentsList(reviewId uint64) (*pb.ReviewCommentsList, error) {
	client, err := s.connectClient("reviews")
	if err != nil {
		return nil, exceptions.HandleAnException(err)
	}
	defer client.Close()

	result, err := client.Get(context.Background(), fmt.Sprintf("reviews_comments:%d", reviewId)).Result()
	if err != nil {
		return nil, exceptions.HandleAnException(err)
	}

	var dto pb.ReviewCommentsList
	err = json.Unmarshal([]byte(result), &dto)
	if err != nil {
		return nil, exceptions.HandleAnException(err)
	}

	return &dto, nil
}

func (s *CacheService) ClearReviewDetails(reviewId uint64) error {
	client, err := s.connectClient("reviews")
	if err != nil {
		return exceptions.HandleAnException(err)
	}
	defer client.Close()

	if err = client.Del(context.Background(), fmt.Sprintf("reviews:%d", reviewId)).Err(); err != nil {
		return exceptions.HandleAnException(err)
	}

	return nil
}

func (s *CacheService) SetReviewDetails(reviewId uint64, dto *pb.ReviewResponse) error {
	client, err := s.connectClient("reviews")
	if err != nil {
		return exceptions.HandleAnException(err)
	}
	defer client.Close()

	payload, err := json.Marshal(dto)
	if err != nil {
		return exceptions.HandleAnException(err)
	}

	if err = client.Set(context.Background(), fmt.Sprintf("reviews:%d", reviewId), payload, s.timeout).Err(); err != nil {
		return exceptions.HandleAnException(err)
	}

	return nil
}

func (s *CacheService) GetReviewDetails(reviewId uint64) (*pb.ReviewResponse, error) {
	client, err := s.connectClient("reviews")
	if err != nil {
		return nil, exceptions.HandleAnException(err)
	}
	defer client.Close()

	result, err := client.Get(context.Background(), fmt.Sprintf("reviews:%d", reviewId)).Result()
	if err != nil {
		return nil, exceptions.HandleAnException(err)
	}

	var dto pb.ReviewResponse
	err = json.Unmarshal([]byte(result), &dto)
	if err != nil {
		return nil, exceptions.HandleAnException(err)
	}

	return &dto, nil
}
