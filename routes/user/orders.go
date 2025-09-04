package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	authPb "github.com/noo8xl/anvil-api/main/auth"
	notificationsPb "github.com/noo8xl/anvil-api/main/notifications"
	offersPb "github.com/noo8xl/anvil-api/main/offers"
	ordersPb "github.com/noo8xl/anvil-api/main/orders"
	reviewsPb "github.com/noo8xl/anvil-api/main/reviews"
	"github.com/noo8xl/anvil-gateway/middlewares"
)

// @description -> Create a new order
//
// @route -> /api/v1/orders/create/
//
// @method -> POST
//
// @body -> body should follow the following structure:
//
//	{
//		OrderBasics: {
//			CustomerId: uint64
//			ApplicantId: uint64
//			OfferId: uint64
//			OrderId: uint64
//		},
//		OrderDetails: {
//			Title: string
//			Body: string
//			Price: float64
//			Status: string -> can be omitted here
//			CreatedAt: string -> can be omitted here
//			UpdatedAt: string -> can be omitted here
//			ExpireAt: string
//			Status: string -> can be omitted here
//			Price: uint64
//		},
//		PaymentDetails: {
//			TotalPrice: float64
//			RoundAmount: uint32
//			IsPaid: bool -> can be omitted here
//			PaidAt: string -> can be omitted here
//			UpdatedAt: string -> can be omitted here
//			PaymentRounds: []PaymentRound{
//
//	{
//					Id: uint64 -> can be omitted here
//					Round: uint32
//					Amount: float64
//					PaidStatus: string
//					PaidAt: string -> can be omitted here
//					UpdatedAt: string -> can be omitted here
//	},
//
//	{
//					Id: uint64 -> can be omitted here
//					Round: uint32
//					Amount: float64
//					PaidStatus: string
//					PaidAt: string -> can be omitted here
//					UpdatedAt: string -> can be omitted here
//	},
//
//			}
//		}
//	}
//
// @response 201
//
// @response 400 {object} ErrorResponse {error: err text}
//
// @response 401 {object} ErrorResponse {error: err text}
//
// @response 403 {object} ErrorResponse {error: err text}
//
// @response 500 {object} ErrorResponse {error: err text}
func (h *Handler) CreateOrderHandler(w http.ResponseWriter, r *http.Request) {

	customerEmail := r.Context().Value(middlewares.CustomerKey).(*authPb.CustomerDto).Email

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	var payload *ordersPb.OrderRequest
	err = json.Unmarshal(body, &payload)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	var totalSum float64
	for _, round := range payload.PaymentDetails.PaymentRounds {
		totalSum += round.Amount
	}

	if totalSum != payload.OrderDetails.Price {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "total sum of payment rounds is not equal to the order price",
		})
		return
	}

	// get title and body from offers by offerId
	offer, err := h.offersClient.GetOfferDetails(context.Background(), &offersPb.GetOfferDetailsRequest{OfferId: payload.OrderBasics.OfferId})
	if err != nil {
		if strings.Contains(err.Error(), "offer not found") {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "offer not found"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	payload.OrderDetails.Title = offer.OfferDetails.Title
	payload.OrderDetails.Body = offer.OfferDetails.Description

	notificationBody := fmt.Sprintf("You have a new invitation to order: %s. To see details visit your profile.", payload.OrderDetails.Title)
	type Result struct {
		service string
		err     error
	}
	var wg sync.WaitGroup
	creationResults := make(chan Result, 2)
	notificationResults := make(chan Result, 2)
	wg.Add(2)

	go func() {
		defer wg.Done()
		_, err = h.ordersClient.CreateOrder(context.Background(), payload)
		creationResults <- Result{service: "orders", err: err}
	}()

	go func() {
		defer wg.Done()
		_, err = h.offersClient.ChangeOfferStatus(context.Background(), &offersPb.ChangeOfferStatusRequest{
			OfferId: payload.OrderBasics.OfferId,
			Status:  "PENDING",
		})
		creationResults <- Result{service: "offers", err: err}
	}()

	wg.Wait()
	close(creationResults)

	for result := range creationResults {
		if result.err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": fmt.Sprintf("failed to %s: %s", result.service, result.err.Error()),
			})
			return
		}
	}

	wg.Add(2)
	go func() {
		defer wg.Done()
		_, err = h.notificationsClient.CreateNotification(context.Background(), &notificationsPb.CreateNotificationRequest{
			CustomerId: payload.OrderBasics.ApplicantId,
			Title:      "New Order Invitation",
			Body:       notificationBody,
			Area:       "orders",
			CreatedAt:  time.Now().Format(time.RFC3339),
		})
		notificationResults <- Result{service: "notifications", err: err}
	}()

	go func() {
		defer wg.Done()
		_, err = h.notificationsClient.SendEmail(context.Background(), &notificationsPb.SendEmailRequest{
			Email:   customerEmail,
			Subject: "New Order Invitation",
			Body:    notificationBody,
		})
		notificationResults <- Result{service: "notifications", err: err}
	}()

	wg.Wait()
	close(notificationResults)

	for result := range notificationResults {
		if result.err != nil {
			json.NewEncoder(w).Encode(map[string]string{
				"error": fmt.Sprintf("failed to %s: %s", result.service, result.err.Error()),
			})
		}
	}

	w.WriteHeader(http.StatusCreated)
}

// @description -> Update an order
//
// @route -> /api/v1/orders/update/
//
// @method -> PUT
//
// @body -> body should follow the following structure:
//
//	{
//		OrderBasics: {
//			CustomerId: uint64
//			ApplicantId: uint64
//			OfferId: uint64
//			OrderId: uint64
//		},
//		OrderDetails: {
//			Title: string
//			Body: string
//			Price: float64
//			Status: string
//			CreatedAt: string -> can be omitted here
//			UpdatedAt: string -> can be omitted here
//			ExpireAt: string
//			Status: string -> can be omitted here
//			Price: uint64
//		},
//		PaymentDetails: {
//			TotalPrice: float64
//			RoundAmount: uint32
//			IsPaid: bool -> can be omitted here
//			PaidAt: string -> can be omitted here
//			UpdatedAt: string -> can be omitted here
//			PaymentRounds: []PaymentRound{
//
//	{
//					Id: uint64 -> can be omitted here
//					Round: uint32
//					Amount: float64
//					PaidStatus: string
//					PaidAt: string -> can be omitted here
//					UpdatedAt: string -> can be omitted here
//	},
//
//	{
//					Id: uint64 -> can be omitted here
//					Round: uint32
//					Amount: float64
//					PaidStatus: string
//					PaidAt: string -> can be omitted here
//					UpdatedAt: string -> can be omitted here
//	},
//
//			}
//		}
//	}
//
// @response 204
//
// @response 400 {object} ErrorResponse {error: err text}
//
// @response 401 {object} ErrorResponse {error: err text}
//
// @response 403 {object} ErrorResponse {error: err text}
//
// @response 500 {object} ErrorResponse {error: err text}
func (h *Handler) UpdateOrderHandler(w http.ResponseWriter, r *http.Request) {

	role := r.Context().Value(middlewares.CustomerKey).(*authPb.CustomerDto).Role

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	var dto *ordersPb.OrderRequest
	err = json.Unmarshal(body, &dto)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	var totalSum float64
	for _, round := range dto.PaymentDetails.PaymentRounds {
		totalSum += round.Amount
	}

	if totalSum != dto.OrderDetails.Price {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "total sum of payment rounds is not equal to the order price",
		})
		return
	}

	s, err := h.ordersClient.GetOrderStatus(context.Background(), &ordersPb.GetOrderStatusRequest{OrderId: dto.OrderBasics.OrderId})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	if s.Status == "APPLIED" {

		if role != "ADMIN" && role != "SUPERVISOR" {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "forbidden. You are not allowed to update this order",
			})
			return
		} else {
			_, err = h.ordersClient.UpdateOrder(context.Background(), dto)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{
					"error": err.Error(),
				})
				return
			}

			err = h.cacheService.ClearOrderDetails(dto.OrderBasics.OrderId)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{
					"error": err.Error(),
				})
				return
			}
		}

		w.WriteHeader(http.StatusNoContent)
		return
	}

	type Result struct {
		service string
		err     error
	}

	notificationBody := fmt.Sprintf("Order %d has been updated. To see details visit your profile.", dto.OrderBasics.OrderId)
	var wg sync.WaitGroup
	results := make(chan Result, 3)
	wg.Add(3)

	go func() {
		defer wg.Done()
		_, err = h.ordersClient.UpdateOrder(context.Background(), dto)
		results <- Result{service: "orders", err: err}
	}()

	go func() {
		defer wg.Done()
		_, err = h.notificationsClient.CreateNotification(context.Background(), &notificationsPb.CreateNotificationRequest{
			CustomerId: dto.OrderBasics.ApplicantId,
			Title:      "Order Updated!",
			Body:       notificationBody,
			Area:       "orders",
			CreatedAt:  time.Now().Format(time.RFC3339),
		})
		results <- Result{service: "notifications", err: err}
	}()

	go func() {
		defer wg.Done()
		err = h.cacheService.ClearOrderDetails(dto.OrderBasics.OrderId)
		results <- Result{service: "cache", err: err}
	}()

	wg.Wait()

	for result := range results {
		if result.err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": fmt.Sprintf("failed to %s: %s", result.service, result.err.Error()),
			})
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

// @description -> Get order details by order id
//
// @route -> /api/v1/orders/get-order-details/{orderId}/
//
// @method -> GET
//
// @body -> an empty one
//
// @response 200
//
//	{
//
//	OrderBasics:
//	{
//	CustomerId: uint64
//	ApplicantId: uint64
//	OfferId: uint64
//	OrderId: uint64
//	}
//	OrderDetails:
//	{
//	Id: uint64
//	Rank: float64
//	CompletedProjects: uint32
//	},
//	Applicant:
//	{
//	Id: uint64
//	Rank: float64
//	CompletedProjects: uint32
//	},
//	PaymentDetails:
//	{
//	TotalPrice: float64
//	RoundAmount: uint32
//	IsPaid: bool
//	PaidAt: string
//	PaymentRounds:
//	{
//	Id: uint64
//	Round: uint32
//	Amount: float64
//	PaidStatus: string
//	PaidAt: string
//	UpdatedAt: string
//	},
//	{
//	Id: uint64
//	Round: uint32
//	Amount: float64
//	PaidStatus: string
//	PaidAt: string
//	UpdatedAt: string
//	},
//	},
//
//	}
//
// @response 400 {object} ErrorResponse {error: err text}
//
// @response 401 {object} ErrorResponse {error: err text}
//
// @response 403 {object} ErrorResponse {error: err text}
//
// @response 500 {object} ErrorResponse {error: err text}
func (h *Handler) GetOrderDetailsHandler(w http.ResponseWriter, r *http.Request) {

	orderId, err := strconv.ParseUint(r.PathValue("orderId"), 10, 64)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Try to get from cache first
	order, err := h.cacheService.GetOrderDetails(orderId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	if order != nil {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(order)
		return
	}

	// Get order details if not in cache
	order, err = h.ordersClient.GetOrderDetails(context.Background(), &ordersPb.GetOrderDetailsRequest{
		OrderId: orderId,
	})
	if err != nil {
		if strings.Contains(err.Error(), "order not found") {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "order not found"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	type Result struct {
		service string
		err     error
	}

	var customerReviewsStats *reviewsPb.CustomerStats
	var applicantReviewsStats *reviewsPb.CustomerStats

	var wg sync.WaitGroup
	results := make(chan Result, 2)
	wg.Add(2)

	go func() {
		defer wg.Done()
		customerReviewsStats, err = h.reviewsClient.GetCustomerStats(context.Background(),
			&reviewsPb.GetCustomerStatsRequest{CustomerId: order.OrderBasics.CustomerId})
		results <- Result{service: "reviews", err: err}
	}()

	go func() {
		defer wg.Done()
		applicantReviewsStats, err = h.reviewsClient.GetCustomerStats(context.Background(),
			&reviewsPb.GetCustomerStatsRequest{CustomerId: order.OrderBasics.ApplicantId})
		results <- Result{service: "reviews", err: err}
	}()

	wg.Wait()
	close(results)

	for result := range results {
		if result.err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": fmt.Sprintf("failed to %s: %s", result.service, result.err.Error()),
			})
			return
		}

		order.Customer.Rank = customerReviewsStats.Rank
		order.Applicant.Rank = applicantReviewsStats.Rank
	}

	if err := h.cacheService.SetOrderDetails(orderId, order); err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
	w.WriteHeader(http.StatusOK)
}

// @description -> Get a list of orders requests by applicantId
//
// @route -> /api/v1/orders/get-orders-requests-list/{skip}/
//
// @method -> GET
//
// @body -> an empty one
//
// @response 200 {object} GetOrdersRequestsListByApplicantIdResponse
//
//	{
//		OrderList: [20]ordersPb.OrderShortCard{
//			OrderId: uint64
//			Title: string
//			Body: string
//			Status: string
//			Price: uint64
//			RoundAmount: uint32
//			CreatedAt: string
//			Skip: uint32
//		}
//	}
//
// @response 400 {object} ErrorResponse {error: err text}
//
// @response 401 {object} ErrorResponse {error: err text}
//
// @response 403 {object} ErrorResponse {error: err text}
//
// @response 500 {object} ErrorResponse {error: err text}
func (h *Handler) GetOrdersRequestsListByApplicantIdHandler(w http.ResponseWriter, r *http.Request) {

	customerId := r.Context().Value(middlewares.CustomerKey).(*authPb.CustomerDto).CustomerId

	skip, err := strconv.ParseUint(r.PathValue("skip"), 10, 64)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	payload := &ordersPb.GetOrderRequestsListByApplicantIdRequest{
		ApplicantId: customerId,
		Skip:        uint32(skip),
	}

	orderList, err := h.ordersClient.GetOrdersRequestsListByApplicantId(context.Background(), payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orderList)
}

// @description -> Get a list of orders by filter
//
// @route -> /api/v1/orders/get-orders-list-by-filter/
//
// @method -> POST
//
// @body -> body should follow the following structure:
//
//	{
//		CustomerId: uint64
//		ApplicantId: uint64
//		Skip: uint32
//		Status: string
//	}
//
// @response 200 {object} GetOrdersListByFilterResponse
//
//	{
//		OrderList: [20]ordersPb.OrderShortCard{
//			OrderId: uint64
//			Title: string
//			Body: string
//			Status: string
//			Price: uint64
//			RoundAmount: uint32
//			CreatedAt: string
//			Skip: uint32
//		}
//	}
//
// @response 400 {object} ErrorResponse {error: err text}
//
// @response 401 {object} ErrorResponse {error: err text}
//
// @response 403 {object} ErrorResponse {error: err text}
//
// @response 500 {object} ErrorResponse {error: err text}
func (h *Handler) GetOrdersListByFilterHandler(w http.ResponseWriter, r *http.Request) {

	customerId := r.Context().Value(middlewares.CustomerKey).(*authPb.CustomerDto).CustomerId

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	orderList, err := h.cacheService.GetFilteredOrdersList(customerId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	if orderList != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(orderList)
		w.WriteHeader(http.StatusOK)
		return
	}

	var payload *ordersPb.GetOrdersListByFilterRequest
	err = json.Unmarshal(body, &payload)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	orderList, err = h.ordersClient.GetOrdersListByFilter(context.Background(), payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	err = h.cacheService.SetFilteredOrdersList(customerId, orderList)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orderList)
}

// @description -> Delete an order by orderId
//
// @route -> /api/v1/orders/delete/{orderId}/
//
// @method -> DELETE
//
// @body -> an empty one
//
// @response 204
//
// @response 400 {error: err text}
//
// @response 401 {error: err text}
//
// @response 403 {error: err text}
//
// @response 500 {error: err text}
func (h *Handler) DeleteOrderHandler(w http.ResponseWriter, r *http.Request) {

	// if order status is applied -> check if the customer has admin or supervisor permission -> delete the order
	// if order status is not applied -> delete the order without any permission check

	orderId, err := strconv.ParseUint(r.PathValue("orderId"), 10, 64)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	customerId := r.Context().Value(middlewares.CustomerKey).(*authPb.CustomerDto).CustomerId
	role := r.Context().Value(middlewares.CustomerKey).(*authPb.CustomerDto).Role

	payload := &ordersPb.DeleteOrderRequest{
		OrderId:    orderId,
		CustomerId: customerId,
	}

	s, err := h.ordersClient.GetOrderStatus(context.Background(), &ordersPb.GetOrderStatusRequest{OrderId: orderId})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	if s.Status == "APPLIED" {

		if role != "ADMIN" && role != "SUPERVISOR" {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "forbidden: you are not allowed to delete this order because it is applied",
			})
			return
		} else {
			_, err = h.ordersClient.DeleteOrder(context.Background(), payload)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{
					"error": err.Error(),
				})
				return
			}
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, err = h.ordersClient.DeleteOrder(context.Background(), payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	err = h.cacheService.ClearOrderDetails(orderId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ############################## orders applications area

// @description -> Apply to the order
//
// @route -> /api/v1/orders/apply/
//
// @method -> POST
//
// @body -> body should follow the following structure:
//
//	{
//		orderBasics: {
//			CustomerId: uint64
//			ApplicantId: uint64
//			OfferId: uint64
//			OrderId: uint64
//		}
//	}
//
// @response 200
//
// @response 400 {error: err text}
//
// @response 401 {error: err text}
//
// @response 403 {error: err text}
//
// @response 500 {error: err text}
func (h *Handler) ApplyToTheOrderHandler(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	var payload *ordersPb.ApplyToTheOrderRequest
	err = json.Unmarshal(body, &payload)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	_, err = h.ordersClient.ApplyToTheOrder(context.Background(), payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	notificationBody := "Congratulations! You have been accepted to the order. To see details visit your profile."
	_, err = h.notificationsClient.CreateNotification(context.Background(), &notificationsPb.CreateNotificationRequest{
		CustomerId: payload.OrderBasics.CustomerId,
		Title:      "New Application",
		Body:       notificationBody,
		Area:       "orders",
		CreatedAt:  time.Now().Format(time.RFC3339),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @description -> Reject an order
//
// @route -> /api/v1/orders/reject/
//
// @method -> POST
//
// @body -> body should follow the following structure:
//
//	{
//		OrderBasics: {
//			CustomerId: uint64
//			ApplicantId: uint64
//			OrderId: uint64
//			OfferId: uint64
//		},
//		Cause: string
//	}
//
// @response 200
//
// @response 400 {error: err text}
//
// @response 401 {error: err text}
//
// @response 403 {error: err text}
//
// @response 500 {error: err text}
func (h *Handler) RejectAnOrderHandler(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	var payload *ordersPb.RejectAnOrderRequest
	err = json.Unmarshal(body, &payload)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	_, err = h.ordersClient.RejectAnOrder(context.Background(), payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	notificationBody := "Your order has been rejected. To see details visit your profile."

	_, err = h.notificationsClient.CreateNotification(context.Background(), &notificationsPb.CreateNotificationRequest{
		CustomerId: payload.OrderBasics.ApplicantId,
		Title:      "Order Rejected",
		Body:       notificationBody,
		Area:       "orders",
		CreatedAt:  time.Now().Format(time.RFC3339),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ############################## orders compliance area

// @description -> Create a new compliance request if customer has done with the order / round
//
// @route -> /api/v1/orders/compliance/create/
//
// @method -> POST
//
// @body -> body should follow the following structure:
//
//	{
//		Id: uint64 -> can be omitted here
//		Round: uint32
//		Status: string
//		Claim: string
//		CreatedAt: string -> can be omitted here
//		OrderBasics: {
//			CustomerId: uint64
//			ApplicantId: uint64
//			OrderId: uint64
//			OfferId: uint64
//		}
//	}
//
// @response 200
//
// @response 400 {error: err text}
//
// @response 401 {error: err text}
//
// @response 403 {error: err text}
//
// @response 500 {error: err text}
func (h *Handler) CreateComplianceRequestHandler(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	var payload *ordersPb.ComplianceRequest
	err = json.Unmarshal(body, &payload)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	_, err = h.ordersClient.CreateComplianceRequest(context.Background(), payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	notificationBody := "You have a new compliance request. To see details visit your profile."

	_, err = h.notificationsClient.CreateNotification(context.Background(), &notificationsPb.CreateNotificationRequest{
		CustomerId: payload.OrderBasics.CustomerId,
		Title:      "New Compliance Request",
		Body:       notificationBody,
		Area:       "orders",
		CreatedAt:  time.Now().Format(time.RFC3339),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @description -> Approve a compliance request by offer owner
//
// @route -> /api/v1/orders/compliance/approve/
//
// @method -> POST
//
// @body -> body should follow the following structure:
//
//	{
//		Id: uint64 -> can be omitted here
//		Round: uint32
//		Status: string
//		Claim: string
//		CreatedAt: string -> can be omitted here
//		OrderBasics: {
//			CustomerId: uint64
//			ApplicantId: uint64спать уже собрался
//		}
//	}
//
// @response 200
//
// @response 400 {error: err text}
//
// @response 401 {error: err text}
//
// @response 403 {error: err text}
//
// @response 500 {error: err text}
func (h *Handler) ComplianceApproveHandler(w http.ResponseWriter, r *http.Request) {

	var order *ordersPb.Order

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	var payload *ordersPb.ComplianceRequest
	err = json.Unmarshal(body, &payload)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	order, err = h.cacheService.GetOrderDetails(payload.OrderBasics.OrderId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	if order == nil {
		order, err = h.ordersClient.GetOrderDetails(context.Background(), &ordersPb.GetOrderDetailsRequest{OrderId: payload.OrderBasics.OrderId})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
			return
		}
	}

	// ------------------------------------
	// TODO:
	// check if the order is paid
	// if paid -> return error
	// else -> continue

	// call the payments service to pay the order
	// update the order status to "APPROVED"
	// update the order payment details
	// update the offer status details
	// send notification to the customer
	// send notification to the applicant
	// clear the order details from the cache
	// return http status 204

	// ------------------------------------

	// ------------------------------------
	// TODO:
	// payment ##
	// round: 1, 2, 3
	// check if the order is paid
	// if paid -> return error
	// else -> continue
	// validate the round amount
	// if round amount more than 1 ->
	// check if previous round is paid
	// if not paid -> return error
	// if paid -> continue
	// if round = 3 -> change order status to "COMPLETED"
	// else -> just complete the round
	// ------------------------------------

	if order.PaymentDetails.IsPaid {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "order is already paid",
		})
		return
	}

	_, err = h.ordersClient.ApproveCompliance(context.Background(), payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	err = h.cacheService.ClearOrderDetails(payload.OrderBasics.OrderId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	notificationBody := "Your compliance request has been approved. To see details visit your profile."
	if payload.Round == order.PaymentDetails.RoundAmount {
		notificationBody = "Your compliance request has been approved and order has been completed! To see details visit your profile."
	}

	_, err = h.notificationsClient.CreateNotification(context.Background(), &notificationsPb.CreateNotificationRequest{
		CustomerId: order.OrderBasics.ApplicantId,
		Title:      "Compliance Request Approved",
		Body:       notificationBody,
		Area:       "orders",
		CreatedAt:  time.Now().Format(time.RFC3339),
	})
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @description -> Reject a compliance request by offer owner
//
// @route -> /api/v1/orders/compliance/reject/
//
// @method -> POST
//
// @body -> body should follow the following structure:
//
//	{
//		Id: uint64
//		Round: uint32
//		Status: string
//		Claim: string
//		CreatedAt: string
//		OrderBasics: {
//			CustomerId: uint64
//			ApplicantId: uint64
//			OrderId: uint64
//			OfferId: uint64
//		}
//	}
//
// @response 200
//
// @response 400 {error: err text}
//
// @response 401 {error: err text}
//
// @response 403 {error: err text}
//
// @response 500 {error: err text}
func (h *Handler) RejectComplianceHandler(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	var payload *ordersPb.ComplianceRequest
	err = json.Unmarshal(body, &payload)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	_, err = h.ordersClient.RejectCompliance(context.Background(), payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	notificationBody := "Your compliance request has been rejected. To see details visit your profile."

	_, err = h.notificationsClient.CreateNotification(context.Background(), &notificationsPb.CreateNotificationRequest{
		CustomerId: payload.OrderBasics.ApplicantId,
		Title:      "Compliance Request Rejected",
		Body:       notificationBody,
		Area:       "orders",
		CreatedAt:  time.Now().Format(time.RFC3339),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	err = h.cacheService.ClearComplianceRequestsList(payload.OrderBasics.CustomerId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @description -> Get a list of compliance requests by customerId
//
// @route -> /api/v1/orders/compliance/get-list/{skip}/
//
// @method -> GET
//
// @body -> an empty one
//
// @response 200
//
//			{
//
//			ComplianceRequests: {
//			{
//			Id: uint64
//			OrderBasics:
//	    {
//	    CustomerId: uint64
//	    ApplicantId: uint64
//	    OrderId: uint64
//	    OfferId: uint64
//	    }
//	    Round: uint32
//	    Status: string
//	    Claim: string
//	    CreatedAt: string
//	    },
//	    {
//	    Id: uint64
//	    OrderBasics:
//	    {
//	    CustomerId: uint64
//	    ApplicantId: uint64
//	    OrderId: uint64
//	    OfferId: uint64
//	    }
//	    Round: uint32
//	    Status: string
//	    Claim: string
//	    CreatedAt: string
//	    },
//	    }
//
// @response 400 {object} ErrorResponse {error: err text}
//
// @response 401 {object} ErrorResponse {error: err text}
//
// @response 403 {object} ErrorResponse {error: err text}
//
// @response 500 {object} ErrorResponse {error: err text}
func (h *Handler) GetComplianceRequestsListHandler(w http.ResponseWriter, r *http.Request) {

	customerId := r.Context().Value(middlewares.CustomerKey).(*authPb.CustomerDto).CustomerId

	list, err := h.cacheService.GetComplianceRequestsList(customerId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	if list != nil {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(list)
		return
	}

	skip, err := strconv.ParseUint(r.PathValue("skip"), 10, 64)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	payload := &ordersPb.GetComplianceRequestsListRequest{
		CustomerId: customerId,
		Skip:       uint32(skip),
	}

	complianceRequestsList, err := h.ordersClient.GetComplianceRequestsList(context.Background(), payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	if err := h.cacheService.SetComplianceRequestsList(customerId, complianceRequestsList); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(complianceRequestsList)
}
