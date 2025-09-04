package cache

import (
	"context"
	"fmt"

	profilePb "github.com/noo8xl/anvil-api/main/profile"
	"github.com/noo8xl/anvil-common/exceptions"
	"google.golang.org/protobuf/proto"
)

// customer profile
func (c *CacheService) ClearCustomerProfile(customerId uint64) error {
	client, err := c.connectClient("profile")
	if err != nil {
		return exceptions.HandleAnException(err)
	}
	defer client.Close()

	status := client.Del(context.Background(), fmt.Sprintf("customer:%d", customerId))
	if status.Err() != nil {
		return exceptions.HandleAnException(status.Err())
	}

	return nil
}

func (c *CacheService) SetCustomerProfile(customerId uint64, customer *profilePb.CustomerResponse) error {
	client, err := c.connectClient("profile")
	if err != nil {
		return exceptions.HandleAnException(err)
	}
	defer client.Close()

	bytes, err := proto.Marshal(customer.ProtoReflect().Interface())
	if err != nil {
		return exceptions.HandleAnException(err)
	}

	if status := client.Set(context.Background(), fmt.Sprintf("customer:%d", customerId), bytes, c.timeout); status.Err() != nil {
		return exceptions.HandleAnException(status.Err())
	}

	return nil
}

func (c *CacheService) GetCustomerProfile(customerId uint64) (*profilePb.CustomerResponse, error) {
	client, err := c.connectClient("profile")
	if err != nil {
		return nil, exceptions.HandleAnException(err)
	}
	defer client.Close()
	customer := client.Get(context.Background(), fmt.Sprintf("customer:%d", customerId))
	if customer.Err() != nil {
		if customer.Err().Error() == "redis: nil" {
			return nil, nil
		}
		return nil, exceptions.HandleAnException(customer.Err())
	}
	if customer.Val() == "" {
		return nil, nil
	}

	customerPb := &profilePb.CustomerResponse{}
	if err := proto.Unmarshal([]byte(customer.Val()), customerPb.ProtoReflect().Interface()); err != nil {
		return nil, exceptions.HandleAnException(err)
	}
	return customerPb, nil
}

// public profile
func (c *CacheService) ClearPublicProfile(customerId uint64) error {
	client, err := c.connectClient("profile")
	if err != nil {
		return exceptions.HandleAnException(err)
	}
	defer client.Close()

	status := client.Del(context.Background(), fmt.Sprintf("public_profile:%d", customerId))
	if status.Err() != nil {
		return exceptions.HandleAnException(status.Err())
	}

	return nil
}

func (c *CacheService) SetPublicProfile(customerId uint64, profile *profilePb.PublicProfileResponse) error {
	client, err := c.connectClient("profile")
	if err != nil {
		return exceptions.HandleAnException(err)
	}
	defer client.Close()

	bytes, err := proto.Marshal(profile.ProtoReflect().Interface())
	if err != nil {
		return exceptions.HandleAnException(err)
	}

	if status := client.Set(context.Background(), fmt.Sprintf("public_profile:%d", customerId), bytes, c.timeout); status.Err() != nil {
		return exceptions.HandleAnException(status.Err())
	}

	return nil
}

func (c *CacheService) GetPublicProfile(customerId uint64) (*profilePb.PublicProfileResponse, error) {
	client, err := c.connectClient("profile")
	if err != nil {
		return nil, exceptions.HandleAnException(err)
	}
	defer client.Close()

	profile := client.Get(context.Background(), fmt.Sprintf("public_profile:%d", customerId))
	if profile.Err() != nil {
		if profile.Err().Error() == "redis: nil" {
			return nil, nil
		}
		return nil, exceptions.HandleAnException(profile.Err())
	}
	if profile.Val() == "" {
		return nil, nil
	}

	profilePb := &profilePb.PublicProfileResponse{}
	if err := proto.Unmarshal([]byte(profile.Val()), profilePb.ProtoReflect().Interface()); err != nil {
		return nil, exceptions.HandleAnException(err)
	}
	return profilePb, nil
}
