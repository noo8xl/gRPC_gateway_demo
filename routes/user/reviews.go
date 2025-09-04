package routes

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	authPb "github.com/noo8xl/anvil-api/main/auth"
	ordersPb "github.com/noo8xl/anvil-api/main/orders"
	reviewsPb "github.com/noo8xl/anvil-api/main/reviews"

	"github.com/noo8xl/anvil-gateway/middlewares"
)

// @des	cription -> Create a review
//
// @route -> /api/v1/reviews/create/
//
// @method -> POST
//
// @body -> body should follow the following structure:
//
//	{
//		ReviewerId: uint64 -> can be omitted here
//		CustomerId: uint64
//		ReviewerId: uint64
//		OrderId: uint64
//		Rank: float32
//		Type: bool
//		Title: string
//		Body: string
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
func (h *Handler) CreateReviewHandler(w http.ResponseWriter, r *http.Request) {

	customerId := r.Context().Value(middlewares.CustomerKey).(*authPb.CustomerDto).CustomerId

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	var payload *reviewsPb.ReviewRequest
	err = json.Unmarshal(body, &payload)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	if err := validateCustomer(customerId, payload.ReviewerId); err != nil {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	validateRelation, err := h.ordersClient.ValidateCustomersRelation(context.Background(), &ordersPb.ValidateCustomersRelationRequest{
		CustomerId:  customerId,
		ApplicantId: payload.ReviewerId,
	})

	if err != nil {
		if strings.Contains(err.Error(), "error: no relation found") {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{
				"forbidden": "You can only review customers you have worked with",
			})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	if validateRelation.OfferId == 0 {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{
			"forbidden": "You can only review customers you have worked with",
		})
		return
	}

	if _, err = h.reviewsClient.CreateReview(context.Background(), payload); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// @description -> Update a review only if you are the reviewer
//
// @route -> /api/v1/reviews/update/
//
// @method -> PUT
//
// @body -> body should follow the following structure:
//
//	{
//		ReviewerId: uint64
//		CustomerId: uint64
//		ReviewerId: uint64
//		OrderId: uint64
//		Rank: float32
//		Type: bool
//		Title: string
//		Body: string
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
func (h *Handler) UpdateReviewHandler(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	var payload *reviewsPb.ReviewRequest
	err = json.Unmarshal(body, &payload)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	customerId := r.Context().Value(middlewares.CustomerKey).(*authPb.CustomerDto).CustomerId
	if err := validateCustomer(customerId, payload.ReviewerId); err != nil {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if _, err = h.reviewsClient.UpdateReview(context.Background(), payload); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @description -> Get a list of reviews for a customer
//
// @route -> /api/v1/reviews/get-reviews-list/{skip}/
//
// @method -> GET
//
// @response 200 {object} reviewsPb.GetReviewsListResponse
//
//	{
//		Reviews: [20]*reviewsPb.ReviewResponse
//
//		{
//			Review: {
//				ReviewId: uint64
//				CustomerId: uint64
//				ReviewerId: uint64
//			},
//			Details: {
//				Rank: float32
//				Title: string
//				Body: string
//				OrderId: uint64
//				CreatedAt: string
//				UpdatedAt: string
//				Helpful: uint32
//				Useless: uint32
//			}
//		}
//	}

// @response 400 {object} ErrorResponse {error: err text}
//
// @response 401 {object} ErrorResponse {error: err text}
//
// @response 403 {object} ErrorResponse {error: err text}
//
// @response 500 {object} ErrorResponse {error: err text}
func (h *Handler) GetReviewsListHandler(w http.ResponseWriter, r *http.Request) {

	customerId := r.Context().Value(middlewares.CustomerKey).(*authPb.CustomerDto).CustomerId

	skip, err := strconv.ParseUint(r.PathValue("skip"), 10, 64)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	payload := &reviewsPb.GetReviewsListRequest{
		CustomerId: customerId,
		Skip:       uint32(skip),
	}

	reviews, err := h.reviewsClient.GetReviewsList(context.Background(), payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	// save it to cache ??

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reviews)

}

// @description -> Delete a review by reviewId
//
// @route -> /api/v1/reviews/delete/{reviewId}/
//
// @method -> DELETE
//
// @body -> an empty one
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
func (h *Handler) DeleteReviewHandler(w http.ResponseWriter, r *http.Request) {

	customerId := r.Context().Value(middlewares.CustomerKey).(*authPb.CustomerDto).CustomerId
	reviewId, err := strconv.ParseUint(r.PathValue("reviewId"), 10, 64)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	_, err = h.reviewsClient.DeleteReview(context.Background(), &reviewsPb.DeleteReviewRequest{
		ReviewId:   reviewId,
		CustomerId: customerId,
	})
	if err != nil {
		if strings.Contains(err.Error(), "error: review not found") {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Review not found",
			})
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @description -> Add a comment to a review
//
// @route -> /api/v1/reviews/add-review-comment/
//
// @method -> POST
//
// @body -> comment body should follow the following structure:
//
//	{
//		Comment: {
//			Comment: string
//			PostedBy: uint64
//			Title: string
//			Body: string
//			CreatedAt: string
//			UpdatedAt: string
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
func (h *Handler) AddReviewCommentHandler(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	var payload *reviewsPb.AddReviewCommentRequest
	err = json.Unmarshal(body, &payload)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	customerId := r.Context().Value(middlewares.CustomerKey).(*authPb.CustomerDto).CustomerId
	if err := validateCustomer(customerId, payload.Comment.PostedBy); err != nil {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if _, err = h.reviewsClient.AddReviewComment(context.Background(), payload); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// @description -> Get a list of comments for a review
//
// @route -> /api/v1/reviews/get-review-comments-list/{reviewId}/{skip}/
//
// @method -> GET
//
// @response 200 {object} reviewsPb.GetReviewCommentsListResponse
//
//	{
//			Comments: [20]*reviewsPb.ReviewCommentResponse
//
//				{
//					Comment: {
//						ReviewId: uint64
//						PostedBy: uint64
//						Title: string
//						Body: string
//						CreatedAt: string
//						UpdatedAt: string
//					}
//				}
//		}
//
// @response 400 {object} ErrorResponse {error: err text}
//
// @response 401 {object} ErrorResponse {error: err text}
//
// @response 403 {object} ErrorResponse {error: err text}
//
// @response 500 {object} ErrorResponse {error: err text}
func (h *Handler) GetReviewCommentsListHandler(w http.ResponseWriter, r *http.Request) {

	reviewId, err := strconv.ParseUint(r.PathValue("reviewId"), 10, 64)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
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

	payload := &reviewsPb.GetReviewCommentsListRequest{
		ReviewId: reviewId,
		Skip:     uint32(skip),
	}

	comments, err := h.reviewsClient.GetReviewCommentsList(context.Background(), payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	if err := h.cacheService.SetReviewCommentsList(reviewId, comments); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

// @description -> Set a reaction to a review
//
// @route -> /api/v1/reviews/reaction/
//
// @method -> POST
//
// @body -> body should follow the following structure:
//
//	{
//		ReviewId: uint64
//		CustomerId: uint64
//		Helpful: bool
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
func (h *Handler) SetReviewReactionHandler(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	var payload *reviewsPb.SetReviewReactionRequest
	err = json.Unmarshal(body, &payload)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	customerId := r.Context().Value(middlewares.CustomerKey).(*authPb.CustomerDto).CustomerId
	if err := validateCustomer(customerId, payload.CustomerId); err != nil {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	_, err = h.reviewsClient.SetReviewReaction(context.Background(), payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
