package cache_test

import (
	"testing"
	"time"

	ordersPb "github.com/noo8xl/anvil-api/main/orders"
	profilePb "github.com/noo8xl/anvil-api/main/profile"
	reviewsPb "github.com/noo8xl/anvil-api/main/reviews"
)

var (
	profile = &profilePb.CustomerResponse{
		Base: &profilePb.Customer{
			CustomerId: customerId,
			Email:      "john.doe@example.com",
			Name:       "John Doe",
			Avatar:     "https://example.com/avatar.png",
		},
		Details: &profilePb.CustomerDetails{
			Role:        "CUSTOMER",
			CreatedAt:   time.Now().Format(time.RFC3339),
			TwoStepType: "OFF",
		},
		Params: &profilePb.CustomerParams{
			IsBanned:   false,
			IsVerified: true,
			IsPremium:  true,
			IsTwoFa:    false,
			IsKyc:      false,
		},
		Kyc: nil,
		Orders: &ordersPb.CustomerStats{
			CompletedProjects: 10,
		},
		Reviews: &reviewsPb.CustomerStats{
			Rank:          4.8,
			TotalProjects: 9,
		},
	}

	publicProfile = &profilePb.PublicProfileResponse{
		Customer: &profilePb.PublicCustomer{
			Id:          customerId,
			Email:       "john.doe@example.com",
			Name:        "John Doe",
			Avatar:      "https://example.com/avatar.png",
			Title:       "Software Engineer",
			Description: "I am a software engineer with a passion for building scalable and efficient systems.",
			Tags:        []string{"software", "engineering", "developer"},
			IsVerified:  false,
			IsPremium:   false,
		},
		Orders: &ordersPb.CustomerStats{
			CompletedProjects: 10,
		},
		Reviews: &reviewsPb.CustomerStats{
			Rank:          4.8,
			TotalProjects: 9,
		},
	}
)

func TestSetCustomerProfile(t *testing.T) {
	if err := svc.SetCustomerProfile(customerId, profile); err != nil {
		t.Fatalf("TestSetCustomerProfile error: failed to set customer profile: %v", err)
	}
}

func TestGetCustomerProfile(t *testing.T) {
	profile, err := svc.GetCustomerProfile(customerId)
	if err != nil {
		t.Fatalf("TestGetCustomerProfile error: failed to get customer profile: %v", err)
	}

	if profile == nil {
		t.Fatalf("TestGetCustomerProfile error: customer profile is nil")
	}

	if profile.Base.CustomerId != customerId {
		t.Fatalf("TestGetCustomerProfile error: customer profile id is not equal to customer id")
	}
}

func TestClearCustomerProfile(t *testing.T) {
	if err := svc.ClearCustomerProfile(customerId); err != nil {
		t.Fatalf("TestClearCustomerProfile error: failed to clear customer profile: %v", err)
	}
}

func TestSetPublicProfile(t *testing.T) {
	if err := svc.SetPublicProfile(customerId, publicProfile); err != nil {
		t.Fatalf("TestSetPublicProfile error: failed to set public profile: %v", err)
	}
}

func TestGetPublicProfile(t *testing.T) {
	pp, err := svc.GetPublicProfile(customerId)
	if err != nil {
		t.Fatalf("TestGetPublicProfile error: failed to get public profile: %v", err)
	}

	if pp == nil {
		t.Fatalf("TestGetPublicProfile error: public profile is nil")
	}

	if pp.Customer.Id != customerId {
		t.Fatalf("TestGetPublicProfile error: public profile id is not equal to customer id")
	}
}

func TestClearPublicProfile(t *testing.T) {
	if err := svc.ClearPublicProfile(customerId); err != nil {
		t.Fatalf("TestClearPublicProfile error: failed to clear public profile: %v", err)
	}
}
