package cache

import (
	"context"
	"encoding/json"
	"fmt"

	pb "github.com/noo8xl/anvil-api/main/orders"
	"github.com/noo8xl/anvil-common/exceptions"
)

func (s *CacheService) ClearOrderDetails(orderId uint64) error {
	client, err := s.connectClient("orders")
	if err != nil {
		return exceptions.HandleAnException(err)
	}
	defer client.Close()

	if err = client.Del(context.Background(), fmt.Sprintf("orders:%d", orderId)).Err(); err != nil {
		return exceptions.HandleAnException(err)
	}

	return nil
}

func (s *CacheService) SetOrderDetails(orderId uint64, dto *pb.Order) error {
	client, err := s.connectClient("orders")
	if err != nil {
		return exceptions.HandleAnException(err)
	}
	defer client.Close()

	payload, err := json.Marshal(dto)
	if err != nil {
		return exceptions.HandleAnException(err)
	}

	if err = client.Set(context.Background(), fmt.Sprintf("orders:%d", orderId), payload, s.timeout).Err(); err != nil {
		return exceptions.HandleAnException(err)
	}

	return nil
}

func (s *CacheService) GetOrderDetails(orderId uint64) (*pb.Order, error) {
	client, err := s.connectClient("orders")
	if err != nil {
		return nil, exceptions.HandleAnException(err)
	}
	defer client.Close()

	result, err := client.Get(context.Background(), fmt.Sprintf("orders:%d", orderId)).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return nil, nil
		}
		return nil, exceptions.HandleAnException(err)
	}

	var dto pb.Order
	err = json.Unmarshal([]byte(result), &dto)
	if err != nil {
		return nil, exceptions.HandleAnException(err)
	}

	return &dto, nil
}

func (s *CacheService) ClearOrdersList(orderId uint64) error {
	client, err := s.connectClient("orders")
	if err != nil {
		return exceptions.HandleAnException(err)
	}
	defer client.Close()

	if err = client.Del(context.Background(), fmt.Sprintf("orders_list:%d", orderId)).Err(); err != nil {
		return exceptions.HandleAnException(err)
	}

	return nil
}

func (s *CacheService) SetOrdersList(dto *pb.OrdersList) error {
	client, err := s.connectClient("orders")
	if err != nil {
		return exceptions.HandleAnException(err)
	}
	defer client.Close()

	payload, err := json.Marshal(dto)
	if err != nil {
		return exceptions.HandleAnException(err)
	}

	for _, order := range dto.OrdersList {
		if err = client.Set(context.Background(), fmt.Sprintf("orders_list:%d", order.OrderId), payload, s.timeout).Err(); err != nil {
			return exceptions.HandleAnException(err)
		}
	}

	return nil
}

func (s *CacheService) GetOrdersList(orderId uint64) (*pb.OrdersList, error) {
	client, err := s.connectClient("orders")
	if err != nil {
		return nil, exceptions.HandleAnException(err)
	}
	defer client.Close()

	result, err := client.Get(context.Background(), fmt.Sprintf("orders_list:%d", orderId)).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return nil, nil
		}
		return nil, exceptions.HandleAnException(err)
	}

	var dto pb.OrdersList
	err = json.Unmarshal([]byte(result), &dto)
	if err != nil {
		return nil, exceptions.HandleAnException(err)
	}

	return &dto, nil
}

func (s *CacheService) ClearFilteredOrdersList(customerId uint64) error {
	client, err := s.connectClient("orders")
	if err != nil {
		return exceptions.HandleAnException(err)
	}
	defer client.Close()

	if err = client.Del(context.Background(), fmt.Sprintf("filtered_orders_list:%d", customerId)).Err(); err != nil {
		return exceptions.HandleAnException(err)
	}

	return nil
}

func (s *CacheService) SetFilteredOrdersList(customerId uint64, dto *pb.OrdersList) error {
	client, err := s.connectClient("orders")
	if err != nil {
		return exceptions.HandleAnException(err)
	}
	defer client.Close()

	payload, err := json.Marshal(dto)
	if err != nil {
		return exceptions.HandleAnException(err)
	}

	if err = client.Set(context.Background(), fmt.Sprintf("filtered_orders_list:%d", customerId), payload, s.timeout).Err(); err != nil {
		return exceptions.HandleAnException(err)
	}

	return nil
}

func (s *CacheService) GetFilteredOrdersList(customerId uint64) (*pb.OrdersList, error) {
	client, err := s.connectClient("orders")
	if err != nil {
		return nil, exceptions.HandleAnException(err)
	}
	defer client.Close()

	result, err := client.Get(context.Background(), fmt.Sprintf("filtered_orders_list:%d", customerId)).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return nil, nil
		}
		return nil, exceptions.HandleAnException(err)
	}

	var dto pb.OrdersList
	err = json.Unmarshal([]byte(result), &dto)
	if err != nil {
		return nil, exceptions.HandleAnException(err)
	}

	return &dto, nil
}

// ############################## orders compliance area

func (s *CacheService) ClearComplianceRequestsList(customerId uint64) error {
	client, err := s.connectClient("orders")
	if err != nil {
		return exceptions.HandleAnException(err)
	}
	defer client.Close()

	if err = client.Del(context.Background(), fmt.Sprintf("compliance_requests_list:%d", customerId)).Err(); err != nil {
		return exceptions.HandleAnException(err)
	}

	return nil
}

func (s *CacheService) SetComplianceRequestsList(customerId uint64, dto *pb.ComplianceRequestsList) error {
	client, err := s.connectClient("orders")
	if err != nil {
		return exceptions.HandleAnException(err)
	}
	defer client.Close()

	payload, err := json.Marshal(dto)
	if err != nil {
		return exceptions.HandleAnException(err)
	}

	if err = client.Set(context.Background(), fmt.Sprintf("compliance_requests_list:%d", customerId), payload, s.timeout).Err(); err != nil {
		return exceptions.HandleAnException(err)
	}

	return nil
}

func (s *CacheService) GetComplianceRequestsList(customerId uint64) (*pb.ComplianceRequestsList, error) {
	client, err := s.connectClient("orders")
	if err != nil {
		return nil, exceptions.HandleAnException(err)
	}
	defer client.Close()

	result, err := client.Get(context.Background(), fmt.Sprintf("compliance_requests_list:%d", customerId)).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return nil, nil
		}
		return nil, exceptions.HandleAnException(err)
	}

	var dto pb.ComplianceRequestsList
	err = json.Unmarshal([]byte(result), &dto)
	if err != nil {
		return nil, exceptions.HandleAnException(err)
	}

	return &dto, nil
}
