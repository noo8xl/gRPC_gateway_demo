package cache_test

import (
	"testing"
	"time"

	ordersPb "github.com/noo8xl/anvil-api/main/orders"
)

var (
	orderId uint64 = 1

	orders = &ordersPb.OrdersList{
		OrdersList: []*ordersPb.OrderShortCard{
			{
				OrderId:     1,
				Title:       "Order 1",
				Body:        "Order 1 description",
				Status:      "ACTIVE",
				Price:       100000,
				RoundAmount: 2,
				CreatedAt:   time.Now().Format(time.RFC3339),
			},
			{
				OrderId:   2,
				Title:     "Order 2",
				Body:      "Order 2 description",
				Status:    "ACTIVE",
				Price:     200000,
				CreatedAt: time.Now().Format(time.RFC3339),
			},
		},
	}

	order = &ordersPb.Order{
		OrderBasics: &ordersPb.OrderBasics{
			CustomerId:  customerId,
			ApplicantId: 4,
			OfferId:     8,
			OrderId:     uint64(orderId),
		},
		OrderDetails: &ordersPb.OrderDetails{
			Title:     "Order 1",
			Body:      "Order 1 description",
			Status:    "ACTIVE",
			Price:     100000,
			CreatedAt: time.Now().Format(time.RFC3339),
			UpdatedAt: time.Now().Format(time.RFC3339),
		},
	}

	complianceRequests = &ordersPb.ComplianceRequestsList{
		ComplianceRequestsList: []*ordersPb.ComplianceRequest{
			{
				Id: 1,
				OrderBasics: &ordersPb.OrderBasics{
					CustomerId:  customerId,
					ApplicantId: 4,
					OfferId:     8,
					OrderId:     uint64(orderId),
				},
				Status:    "ACTIVE",
				Round:     1,
				Claim:     "Test claim",
				CreatedAt: time.Now().Format(time.RFC3339),
			},
			{
				Id: 2,
				OrderBasics: &ordersPb.OrderBasics{
					CustomerId:  customerId,
					ApplicantId: 4,
					OfferId:     8,
					OrderId:     uint64(orderId),
				},
				Status: "ACTIVE",
				Round:  2,
				Claim:  "Test claim 2",
			},
		},
	}
)

func TestSetOrdersList(t *testing.T) {
	if err := svc.SetOrdersList(orders); err != nil {
		t.Fatalf("TestSetOrdersList error: failed to set orders list: %v", err)
	}
}

func TestSetOrderDetails(t *testing.T) {
	if err := svc.SetOrderDetails(orderId, order); err != nil {
		t.Fatalf("TestSetOrderDetails error: failed to set order details: %v", err)
	}
}

func TestSetFilteredOrdersList(t *testing.T) {
	if err := svc.SetFilteredOrdersList(customerId, orders); err != nil {
		t.Fatalf("TestSetFilteredOrdersList error: failed to set filtered orders list: %v", err)
	}
}

func TestSetComplianceRequestsList(t *testing.T) {
	if err := svc.SetComplianceRequestsList(customerId, complianceRequests); err != nil {
		t.Fatalf("TestSetComplianceRequestsList error: failed to set compliance requests list: %v", err)
	}
}

func TestGetOrderDetails(t *testing.T) {
	order, err := svc.GetOrderDetails(orderId)
	if err != nil {
		t.Fatalf("TestGetOrderDetails error: failed to get order details: %v", err)
	}

	if order == nil {
		t.Fatalf("TestGetOrderDetails error: order details is nil")
	}
}

func TestGetOrdersList(t *testing.T) {
	orders, err := svc.GetOrdersList(customerId)
	if err != nil {
		t.Fatalf("TestGetOrdersList error: failed to get orders list: %v", err)
	}

	if orders != nil {
		if len(orders.OrdersList) == 0 {
			t.Fatalf("TestGetOrdersList error: orders list is empty")
		}
	}
}

func TestGetFilteredOrdersList(t *testing.T) {
	orders, err := svc.GetFilteredOrdersList(customerId)
	if err != nil {
		t.Fatalf("TestGetFilteredOrdersList error: failed to get filtered orders list: %v", err)
	}

	if orders == nil {
		t.Fatalf("TestGetFilteredOrdersList error: filtered orders list is nil")
	}
}

func TestGetComplianceRequestsList(t *testing.T) {
	list, err := svc.GetComplianceRequestsList(customerId)
	if err != nil {
		t.Fatalf("TestGetComplianceRequestsList error: failed to get compliance requests list: %v", err)
	}

	if list == nil {
		t.Fatalf("TestGetComplianceRequestsList error: compliance requests list is nil")
	}
}

func TestClearComplianceRequestsList(t *testing.T) {
	if err := svc.ClearComplianceRequestsList(customerId); err != nil {
		t.Fatalf("TestClearComplianceRequestsList error: failed to clear compliance requests list: %v", err)
	}
}

func TestClearOrdersList(t *testing.T) {
	if err := svc.ClearOrdersList(customerId); err != nil {
		t.Fatalf("TestClearOrdersList error: failed to clear orders list: %v", err)
	}
}

func TestClearFilteredOrdersList(t *testing.T) {
	if err := svc.ClearFilteredOrdersList(customerId); err != nil {
		t.Fatalf("TestClearFilteredOrdersList error: failed to clear filtered orders list: %v", err)
	}
}

func TestClearOrderDetails(t *testing.T) {
	if err := svc.ClearOrderDetails(orderId); err != nil {
		t.Fatalf("TestClearOrderDetails error: failed to clear order details: %v", err)
	}
}
