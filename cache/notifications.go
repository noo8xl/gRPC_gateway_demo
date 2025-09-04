package cache

import (
	"context"
	"encoding/json"
	"fmt"

	notificationPb "github.com/noo8xl/anvil-api/main/notifications"
	"github.com/noo8xl/anvil-common/exceptions"
)

// ###########################################
// used as a temp storage for notifications
// and renew by timeout or
// if customer got a new notification
// ###########################################

func (s *CacheService) ClearNotifications(notificationId uint64) error {
	client, err := s.connectClient("notifications")
	if err != nil {
		return exceptions.HandleAnException(err)
	}
	defer client.Close()

	if err = client.Del(context.Background(), fmt.Sprintf("notifications:%d", notificationId)).Err(); err != nil {
		return exceptions.HandleAnException(err)
	}

	return nil
}

func (s *CacheService) SetNotificationsList(list *notificationPb.GetNotificationsListResponse) error {
	client, err := s.connectClient("notifications")
	if err != nil {
		return exceptions.HandleAnException(err)
	}
	defer client.Close()

	payload, err := json.Marshal(list)
	if err != nil {
		return exceptions.HandleAnException(err)
	}

	if err = client.Set(context.Background(), fmt.Sprintf("notifications:%d", list.List[0].CustomerId), payload, s.timeout).Err(); err != nil {
		return exceptions.HandleAnException(err)
	}

	return nil
}

func (s *CacheService) GetNotificationsList(customerId uint64) (*notificationPb.GetNotificationsListResponse, error) {
	client, err := s.connectClient("notifications")
	if err != nil {
		return nil, exceptions.HandleAnException(err)
	}
	defer client.Close()

	result, err := client.Get(context.Background(), fmt.Sprintf("notifications:%d", customerId)).Result()
	if err != nil {
		return nil, exceptions.HandleAnException(err)
	}

	var list notificationPb.GetNotificationsListResponse
	if err := json.Unmarshal([]byte(result), &list); err != nil {
		return nil, exceptions.HandleAnException(err)
	}

	return &list, nil

}
