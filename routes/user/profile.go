package routes

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"

	authPb "github.com/noo8xl/anvil-api/main/auth"
	ordersPb "github.com/noo8xl/anvil-api/main/orders"
	profilePb "github.com/noo8xl/anvil-api/main/profile"
	reviewsPb "github.com/noo8xl/anvil-api/main/reviews"
	"github.com/noo8xl/anvil-gateway/middlewares"
)

// profile area ##############################################################

// @description -> Create a customer
//
// @route -> this route is available from the client side without the admin permission *
//
// @method -> POST
//
// @response 201
//
// @response 400 {object} {error: err text}
//
// @response 401 {object} {error: err text}
//
// @response 403 {object} {error: err text}
//
// @response 500 {object} {error: err text}
func (h *Handler) CreateCustomer(w http.ResponseWriter, r *http.Request) {

	dto := &profilePb.CreateCustomerRequest{
		Email:    r.PathValue("email"),
		Name:     r.PathValue("name"),
		Password: r.PathValue("password"),
	}
	_, err := h.profileClient.CreateCustomer(context.Background(), dto)
	if err != nil {
		if err.Error() == "customer already exists" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(map[string]string{"error": "customer already exists"})
		} else {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		}
	}

	w.WriteHeader(201)
}

// @description -> Fill a customer profile (bio, which use in public profile)
//
// @route -> /api/v1/profile/fill/
//
// @method -> POST
//
// @body -> body should follow the following structure:
//
//	{
//
//			Base: {
//				CustomerId: uint64
//				Email: string -> use customer current email by default
//	     	Name: string
//				Avatar: string
//			}
//			Details: {
//				Title: string
//				Description: string
//				Tags: []string
//				Categories: []string
//			}
//
//	}
//
// @response 201
//
// @response 400 {object} {error: err text}
//
// @response 401 {object} {error: err text}
//
// @response 403 {object} {error: err text}
//
// @response 500 {object} {error: err text}
func (h *Handler) FillProfileHandler(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	defer r.Body.Close()

	var dto *profilePb.ProfileBioRequest
	err = json.Unmarshal(body, &dto)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	customerId := r.Context().Value(middlewares.CustomerKey).(*authPb.CustomerDto).CustomerId
	if err := validateCustomer(customerId, dto.Base.CustomerId); err != nil {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	_, err = h.profileClient.FillProfile(context.Background(), dto)
	if err != nil {
		if strings.Contains(err.Error(), "customer not found") {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(map[string]string{"error": "customer not found"})
			return
		}
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(201)
}

// @description -> Update a customer bio only
//
// @route -> /api/v1/profile/update/
//
// @method -> POST
//
// @body -> body should follow the following structure:
//
//	{
//
//			Base: {
//				CustomerId: uint64
//				Email: string -> can be omitted here
//	     	Name: string
//				Avatar: string
//			}
//			Details: {
//				Title: string
//				Description: string
//				Tags: []string
//				Categories: []string
//			}
//
//	}
//
// @response 200
//
// @response 400 {object} {error: err text}
//
// @response 401 {object} {error: err text}
//
// @response 403 {object} {error: err text}
//
// @response 500 {object} {error: err text}
func (h *Handler) UpdateCustomerProfileHandler(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	defer r.Body.Close()

	var dto *profilePb.ProfileBioRequest
	err = json.Unmarshal(body, &dto)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	customerId := r.Context().Value(middlewares.CustomerKey).(*authPb.CustomerDto).CustomerId
	if err := validateCustomer(customerId, dto.Base.CustomerId); err != nil {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	customer, err := h.profileClient.UpdateCustomerProfile(context.Background(), dto)
	if err != nil {
		if strings.Contains(err.Error(), "customer not found") {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(map[string]string{"error": "customer not found"})
			return
		}
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	h.cacheService.ClearPublicProfile(dto.Base.CustomerId)
	h.cacheService.ClearCustomerProfile(dto.Base.CustomerId)
	h.cacheService.SetCustomerProfile(dto.Base.CustomerId, customer)

	w.WriteHeader(200)
}

// @description -> Get a customer profile by id for the customer, not a public one.
// If you don't got something - it cause an empty or false fields was omitted.
//
// @route -> /api/v1/profile/get/
//
// @method -> GET
//
// @body -> an empty one
//
// @response 200 {object}
//
//		{
//			Base:
//			{
//			CustomerId: uint64
//	    Email: string
//	    Name: string
//	    Avatar: string
//	    }
//	    Details:
//	    {
//	    Role: string
//	    CreatedAt: string
//	    TwoStepType: string
//	    }
//	    Params:
//	    {
//	    IsBanned: bool
//	    IsVerified: bool
//	    IsPremium: bool
//	    IsTwoFa: bool
//	    IsKyc: bool
//	    }
//	    Bio:
//	    {
//	    Title: string
//	    Description: string
//	    Tags: []string
//	    Categories: []string
//	    }
//	    Kyc: {}
//	    Orders:
//	    {
//	    CompletedProjects: uint32
//	    }
//	    Reviews:
//	    {
//	    TotalProjects: uint32
//	    Rank: float64
//	    }
//	    }
//
// @response 400 {object} ErrorResponse {error: err text}
//
// @response 401 {object} ErrorResponse {error: err text}
//
// @response 403 {object} ErrorResponse {error: err text}
//
// @response 500 {object} ErrorResponse {error: err text}
func (h *Handler) GetCustomerProfileHandler(w http.ResponseWriter, r *http.Request) {

	customerDto := r.Context().Value(middlewares.CustomerKey).(*authPb.CustomerDto)
	if customerDto == nil {
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized. customer not found in context"})
		return
	}

	customerId := customerDto.CustomerId

	type Result struct {
		customer     *profilePb.CustomerResponse
		orderStats   *ordersPb.CustomerStats
		reviewsStats *reviewsPb.CustomerStats
		customerKyc  *profilePb.CustomerKyc
		err          error
	}

	var response *profilePb.CustomerResponse
	var customer *profilePb.CustomerResponse
	var orderStats *ordersPb.CustomerStats
	var reviewsStats *reviewsPb.CustomerStats
	var customerKyc *profilePb.CustomerKyc
	results := make(chan Result, 3)
	var wg sync.WaitGroup
	wg.Add(3)

	response, err := h.cacheService.GetCustomerProfile(customerId)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	if response != nil {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(200)
		return
	}

	go func() {
		defer wg.Done()
		orderStats, err := h.ordersClient.GetCustomerStats(context.Background(), &ordersPb.GetCustomerStatsRequest{
			CustomerId: customerId,
		})
		results <- Result{orderStats: orderStats, err: err}
	}()

	go func() {
		defer wg.Done()
		reviewsStats, err := h.reviewsClient.GetCustomerStats(context.Background(), &reviewsPb.GetCustomerStatsRequest{
			CustomerId: customerId,
		})
		results <- Result{reviewsStats: reviewsStats, err: err}
	}()

	go func() {
		defer wg.Done()
		customer, err := h.profileClient.GetCustomerProfile(context.Background(), &profilePb.GetCustomerProfileRequest{
			CustomerId: customerId,
		})
		results <- Result{customer: customer, err: err}
	}()

	// go func() {
	// 	defer wg.Done()
	// 	kyc, err := h.profileClient.GetCustomerKycProfile(context.Background(), &profilePb.GetCustomerKycRequest{
	// 		CustomerId: customerId,
	// 	})
	// 	results <- Result{customerKyc: kyc, err: err}
	// }()

	wg.Wait()
	close(results)

	for result := range results {
		if result.err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(map[string]string{"error": result.err.Error()})
			return
		}
		if result.customer != nil {
			customer = result.customer
		}
		if result.orderStats != nil {
			orderStats = result.orderStats
		}
		if result.reviewsStats != nil {
			reviewsStats = result.reviewsStats
		}
		if result.customerKyc != nil {
			customerKyc = result.customerKyc
		}
	}

	response = &profilePb.CustomerResponse{
		Base:    customer.Base,
		Details: customer.Details,
		Params:  customer.Params,
		Bio:     customer.Bio,
		Kyc:     customerKyc,
		Orders:  orderStats,
		Reviews: reviewsStats,
	}

	if err := h.cacheService.SetCustomerProfile(customerId, response); err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// @description -> Get a public profile by id
//
// @route -> /api/v1/profile/public/get/
//
// @method -> GET
//
// @body -> an empty one
//
// @response 200 {object}
//
//		{
//			Id: uint64
//			Name: string
//			Email: string
//			Avatar: string
//	    Title: string
//	    Description: string
//	    Tags: []string
//	    Categories: []string
//	    IsVerified: bool
//	    IsPremium: bool
//	  }
//
// @response 400 {object} ErrorResponse {error: err text}
//
// @response 401 {object} ErrorResponse {error: err text}
//
// @response 403 {object} ErrorResponse {error: err text}
//
// @response 500 {object} ErrorResponse {error: err text}
func (h *Handler) GetPublicProfileHandler(w http.ResponseWriter, r *http.Request) {

	customerDto := r.Context().Value(middlewares.CustomerKey).(*authPb.CustomerDto)
	if customerDto == nil {
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized. customer not found in context"})
		return
	}
	customerId := customerDto.CustomerId

	type Result struct {
		customer *profilePb.PublicCustomer
		stats    *ordersPb.CustomerStats
		reviews  *reviewsPb.CustomerStats
		err      error
	}

	var response *profilePb.PublicProfileResponse
	var customer *profilePb.PublicCustomer
	var stats *ordersPb.CustomerStats
	var reviews *reviewsPb.CustomerStats

	results := make(chan Result, 3)
	var wg sync.WaitGroup
	wg.Add(3)

	response, err := h.cacheService.GetPublicProfile(customerId)
	if err != nil {
		if err.Error() != "redis: nil" {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
	}
	if response != nil {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(200)
		return
	}

	go func() {
		defer wg.Done()
		customer, err := h.profileClient.GetPublicProfile(context.Background(), &profilePb.GetPublicProfileRequest{
			CustomerId: customerId,
		})
		results <- Result{customer: customer, err: err}
	}()

	go func() {
		defer wg.Done()
		stats, err := h.ordersClient.GetCustomerStats(context.Background(), &ordersPb.GetCustomerStatsRequest{
			CustomerId: customerId,
		})
		results <- Result{stats: stats, err: err}
	}()

	go func() {
		defer wg.Done()
		reviews, err := h.reviewsClient.GetCustomerStats(context.Background(), &reviewsPb.GetCustomerStatsRequest{
			CustomerId: customerId,
		})
		results <- Result{reviews: reviews, err: err}
	}()

	wg.Wait()
	close(results)

	for result := range results {
		if result.err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(map[string]string{"error": result.err.Error()})
			return
		}
		if result.customer != nil {
			customer = result.customer
		}
		if result.stats != nil {
			stats = result.stats
		}
		if result.reviews != nil {
			reviews = result.reviews
		}
	}

	response = &profilePb.PublicProfileResponse{
		Customer: customer,
		Orders:   stats,
		Reviews:  reviews,
	}

	h.cacheService.SetPublicProfile(customerId, response)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
}

// kyc area ##############################################################

// @description -> Create a customer kyc
//
// @route -> /api/v1/profile/kyc/create/
//
// @method -> POST
//
// @body -> not implemented yet
//
// @response 201 {object}
//
// @response 400 {object} ErrorResponse {error: err text}
//
// @response 401 {object} ErrorResponse {error: err text}
//
// @response 403 {object} ErrorResponse {error: err text}
//
// @response 500 {object} ErrorResponse {error: err text}
func (h *Handler) CreateCustomeKycHandler(w http.ResponseWriter, r *http.Request) {

	// customerId := r.Context().Value(middlewares.CustomerKey).(*authPb.CustomerDto).CustomerId
	// if err := validateCustomer(customerId, dto.Base.CustomerId); err != nil {
	// 	w.WriteHeader(http.StatusForbidden)
	// 	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
	// 	return
	// }

}

// @description -> Update a customer kyc data
//
// @route -> /api/v1/profile/kyc/update/
//
// @method -> PUT
//
// @body -> not implemented yet
//
// @response 200 {object}
//
// @response 400 {object} ErrorResponse {error: err text}
//
// @response 401 {object} ErrorResponse {error: err text}
//
// @response 403 {object} ErrorResponse {error: err text}
//
// @response 500 {object} ErrorResponse {error: err text}
func (h *Handler) UpdateCustomeKycHandler(w http.ResponseWriter, r *http.Request) {

	// customerId := r.Context().Value(middlewares.CustomerKey).(*authPb.CustomerDto).CustomerId
	// if err := validateCustomer(customerId, dto.Base.CustomerId); err != nil {
	// 	w.WriteHeader(http.StatusForbidden)
	// 	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
	// 	return
	// }

}

// @description -> Get a customer kyc by customerId
//
// @route -> /api/v1/profile/kyc/get/
//
// @method -> GET
//
// @body -> not implemented yet
//
// @response 200 {object} not implemented yet
//
// @response 400 {object} ErrorResponse {error: err text}
//
// @response 401 {object} ErrorResponse {error: err text}
//
// @response 403 {object} ErrorResponse {error: err text}
//
// @response 500 {object} ErrorResponse {error: err text}
func (h *Handler) GetCustomerKycHandler(w http.ResponseWriter, r *http.Request) {

	// customerId := r.Context().Value(middlewares.CustomerKey).(*authPb.CustomerDto).CustomerId

}

// @description -> Report a customer
//
// @route -> /api/v1/profile/report/
//
// @method -> POST
//
// @body -> body should follow the following structure:
//
//	{
//		CustomerId: uint64
//		ReporterId: uint64
//		Reason: string
//		Description: string
//	}
//
// Response list:
//
// @response 200 {object}
//
// @response 400 {object} ErrorResponse {error: err text}
//
// @response 401 {object} ErrorResponse {error: err text}
//
// @response 403 {object} ErrorResponse {error: err text}
//
// @response 500 {object} ErrorResponse {error: err text}
func (h *Handler) ReportCustomerHandler(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	defer r.Body.Close()

	var dto *profilePb.ReportCustomerRequest
	err = json.Unmarshal(body, &dto)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	customerId := r.Context().Value(middlewares.CustomerKey).(*authPb.CustomerDto).CustomerId
	if err := validateCustomer(customerId, dto.ReporterId); err != nil {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	_, err = h.profileClient.ReportCustomer(context.Background(), dto)
	if err != nil {
		if strings.Split(err.Error(), "desc = ")[1] == "customer or reporter not found" {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(map[string]string{"error": "customer or reporter not found"})
			return
		}
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(200)
}
