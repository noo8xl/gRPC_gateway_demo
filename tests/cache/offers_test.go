package cache_test

import (
	"testing"
	"time"

	offersPb "github.com/noo8xl/anvil-api/main/offers"
)

var (
	offer = &offersPb.Offer{
		OfferId:     1,
		PostedBy:    4,
		AppliedBy:   []uint64{5},
		ApprovedFor: 5,
		OfferDetails: &offersPb.OfferDetails{
			Title:          "Test Offer",
			Price:          100,
			Description:    "Test Description",
			Status:         "ACTIVE",
			Currency:       "USD",
			CreatedAt:      time.Now().Format(time.RFC3339),
			UpdatedAt:      time.Now().Format(time.RFC3339),
			ExpiredAt:      time.Now().Format(time.RFC3339),
			TagsList:       []string{"Test Tag", "Test Tag 2"},
			CategoriesList: []string{"Test Category", "Test Category 2"},
		},
	}
)

func TestSetOfferDetails(t *testing.T) {
	err := svc.SetOfferDetails(offer.OfferId, offer)
	if err != nil {
		t.Errorf("Error setting offer details: %v", err)
	}
}

func TestGetOfferDetails(t *testing.T) {
	offer, err := svc.GetOfferDetails(offer.OfferId)
	if err != nil {
		t.Errorf("Error getting offer details: %v", err)
	}

	if offer != nil {
		if offer.OfferDetails == nil {
			t.Errorf("testGetOfferDetails error: Expected details for offer, got nil")
		}
	}
}

func TestGetOfferDetailsWithInvalidOfferId(t *testing.T) {
	offer, err := svc.GetOfferDetails(666)
	if err != nil {
		t.Errorf("Error getting offer details: %v", err)
	}

	if offer != nil {
		t.Errorf("testGetOfferDetails error: Expected nil for offer")
	}
}

func TestClearOfferDetails(t *testing.T) {
	err := svc.ClearOfferDetails(offer.OfferId)
	if err != nil {
		t.Errorf("Error clearing offer details: %v", err)
	}
}
