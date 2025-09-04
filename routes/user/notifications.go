package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	authPb "github.com/noo8xl/anvil-api/main/auth"
	notificationPb "github.com/noo8xl/anvil-api/main/notifications"
	"github.com/noo8xl/anvil-gateway/middlewares"
)

// @description -> Get notifications list by customer id
//
// @route -> /api/v1/notifications/get-notifications-list/{skip}/
//
// @method -> GET
//
// @body -> an empty one
//
// @response 200:
//
//		{
//		    List: {
//		    Id: uint64
//		    CustomerId: uint64
//		    Area: string
//	    	Title: string
//	    	Body: string
//	    	CreatedAt: string
//	    }
//	   }
//
// @response 400 {object} ErrorResponse
//
// @response 401 {object} ErrorResponse
//
// @response 403 {object} ErrorResponse
//
// @response 500 {object} ErrorResponse
func (h *Handler) GetNotificationsListHandler(w http.ResponseWriter, r *http.Request) {

	customerId := r.Context().Value(middlewares.CustomerKey).(*authPb.CustomerDto).CustomerId
	skip, err := strconv.ParseUint(r.PathValue("skip"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	response, err := h.notificationsClient.GetNotificationsList(context.Background(), &notificationPb.GetNotificationsListRequest{
		CustomerId: customerId,
		Skip:       uint32(skip),
	})
	if err != nil {
		if strings.Contains(err.Error(), "customer not found") {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// @description -> Delete a notification by notificationId
//
// @route -> /api/v1/notifications/delete-notification/{notificationId}/
//
// @method -> DELETE
//
// @body -> an empty one
//
// @response 204
//
// @response 400 {object} ErrorResponse
//
// @response 401 {object} ErrorResponse
//
// @response 403 {object} ErrorResponse
//
// @response 500 {object} ErrorResponse
func (h *Handler) DeleteNotificationHandler(w http.ResponseWriter, r *http.Request) {

	notificationId, err := strconv.ParseUint(r.PathValue("notificationId"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	_, err = h.notificationsClient.DeleteNotification(context.Background(), &notificationPb.DeleteNotificationRequest{
		NotificationId: notificationId,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @description -> Clear all notifications
//
// @route -> /api/v1/notifications/clear-notifications/
//
// @method -> DELETE
//
// @body -> an empty one
//
// Response list:
//
//   - @response 200
//
// @response 400 {object} ErrorResponse
//
// @response 401 {object} ErrorResponse
//
// @response 403 {object} ErrorResponse
//
// @response 500 {object} ErrorResponse
func (h *Handler) ClearNotificationsHandler(w http.ResponseWriter, r *http.Request) {

	customerId := r.Context().Value(middlewares.CustomerKey).(*authPb.CustomerDto).CustomerId

	_, err := h.notificationsClient.ClearNotifications(context.Background(), &notificationPb.ClearNotificationsRequest{
		CustomerId: customerId,
	})
	if err != nil {
		if strings.Contains(err.Error(), "customer not found") {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
}
