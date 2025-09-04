package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/noo8xl/anvil-api/main/offers"
	"github.com/noo8xl/anvil-common/exceptions"
)

// SetOfferDetails -> set offer detailed data
func (s *CacheService) SetOfferDetails(offerId uint64, dto *offers.Offer) error {
	client, err := s.connectClient("offers")
	if err != nil {
		return exceptions.HandleAnException(err)
	}
	defer client.Close()

	payload, err := json.Marshal(dto)
	if err != nil {
		return exceptions.HandleAnException(err)
	}

	if err = client.Set(context.Background(), fmt.Sprintf("offers:%d", dto.OfferId), payload, s.timeout).Err(); err != nil {
		return exceptions.HandleAnException(err)
	}

	return nil
}

// GetOfferDetails -> get offer detailed data
func (s *CacheService) GetOfferDetails(offerId uint64) (*offers.Offer, error) {
	client, err := s.connectClient("offers")
	if err != nil {
		return nil, exceptions.HandleAnException(err)
	}
	defer client.Close()

	result, err := client.Get(context.Background(), fmt.Sprintf("offers:%d", offerId)).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return nil, nil
		}
		return nil, exceptions.HandleAnException(err)
	}

	var dto offers.Offer
	err = json.Unmarshal([]byte(result), &dto)
	if err != nil {
		return nil, exceptions.HandleAnException(err)
	}

	return &dto, nil
}

func (s *CacheService) ClearOfferDetails(offerId uint64) error {
	client, err := s.connectClient("offers")
	if err != nil {
		return exceptions.HandleAnException(err)
	}
	defer client.Close()

	if err = client.Del(context.Background(), fmt.Sprintf("offers:%d", offerId)).Err(); err != nil {
		return exceptions.HandleAnException(err)
	}

	return nil
}

func (s *CacheService) ClearApplicantsList(offerId uint64) error {
	client, err := s.connectClient("offers")
	if err != nil {
		return exceptions.HandleAnException(err)
	}
	defer client.Close()

	if err = client.Del(context.Background(), fmt.Sprintf("applicants:%d", offerId)).Err(); err != nil {
		return exceptions.HandleAnException(err)
	}

	return nil
}

// SetApplicantsList -> set applicants list
func (s *CacheService) SetApplicantsList(offerId uint64, dto *offers.ApplicantsList) error {
	client, err := s.connectClient("offers")
	if err != nil {
		return exceptions.HandleAnException(err)
	}
	defer client.Close()

	payload, err := json.Marshal(dto)
	if err != nil {
		return exceptions.HandleAnException(err)
	}

	if err = client.Set(context.Background(), fmt.Sprintf("applicants:%d", offerId), payload, s.timeout).Err(); err != nil {
		return exceptions.HandleAnException(err)
	}

	return nil
}

// GetApplicantsList -> get applicants list
func (s *CacheService) GetApplicantsList(offerId uint64) (*offers.ApplicantsList, error) {
	client, err := s.connectClient("offers")
	if err != nil {
		return nil, exceptions.HandleAnException(err)
	}
	defer client.Close()

	result, err := client.Get(context.Background(), fmt.Sprintf("applicants:%d", offerId)).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return nil, nil
		}
		return nil, exceptions.HandleAnException(err)
	}

	var dto offers.ApplicantsList
	err = json.Unmarshal([]byte(result), &dto)
	if err != nil {
		return nil, exceptions.HandleAnException(err)
	}

	return &dto, nil
}
