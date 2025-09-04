package cache_test

import (
	"testing"
	"time"

	notificationPb "github.com/noo8xl/anvil-api/main/notifications"
)

var (
	customerId uint64 = 4

	notificationsList = &notificationPb.GetNotificationsListResponse{
		List: []*notificationPb.Notification{
			{
				Id:         1,
				CustomerId: 4,
				Area:       "profile",
				Title:      "test title",
				Body:       "test body",
				CreatedAt:  time.Now().Format(time.RFC3339),
			},
			{
				Id:         2,
				CustomerId: 4,
				Area:       "orders",
				Title:      "test title 2",
				Body:       "test body 2",
				CreatedAt:  time.Now().Format(time.RFC3339),
			},
		},
	}
)

func TestSetNotificationsList(t *testing.T) {

	err := svc.SetNotificationsList(notificationsList)
	if err != nil {
		t.Errorf("error setting notifications list: Expected nil, got %v", err)
	}

}

func TestGetNotificationsList(t *testing.T) {

	list, err := svc.GetNotificationsList(customerId)
	if err != nil {
		t.Errorf("error getting notifications list: Expected nil, got %v", err)
	}

	if list == nil {
		t.Errorf("error getting notifications list: Expected a list of notifications, got nil")
	}

}

func TestClearNotifications(t *testing.T) {

	err := svc.ClearNotifications(customerId)
	if err != nil {
		t.Errorf("error clearing notifications: Expected nil, got %v", err)
	}

}
