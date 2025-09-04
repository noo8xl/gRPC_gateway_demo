package routes

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	authPb "github.com/noo8xl/anvil-api/main/auth"
	blogPb "github.com/noo8xl/anvil-api/main/blog"
	"github.com/noo8xl/anvil-gateway/middlewares"
)

// @description -> Create a blog post
//
// @route -> /api/v1/blog/create/
//
// @method -> POST
//
// @body -> body should follow the following structure:
//
//	{
//		Id: uint64 -> can be omitted here
//		CustomerId: uint64
//		Title: string
//		Description: string
//		TagList: []string
//		ImageList: []string
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
func (h *Handler) CreateBlogHandler(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	defer r.Body.Close()

	var dto *blogPb.PostRequest
	if err = json.Unmarshal(body, &dto); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	customerId := r.Context().Value(middlewares.CustomerKey).(*authPb.CustomerDto).CustomerId
	if err := validateCustomer(customerId, dto.CustomerId); err != nil {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if _, err = h.blogClient.CreatePost(context.Background(), dto); err != nil {
		if strings.Split(err.Error(), "desc = ")[1] == "customer not found" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "customer not found"})
			return
		}
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(201)
}

// @description -> Update a post
//
// @route -> /api/v1/blog/update/
//
// @method -> POST
//
// @body -> body should follow the following structure:
//
//	{
//		Id: uint64
//		CustomerId: uint64
//		Title: string
//		Description: string
//		TagList: []string
//		ImageList: []string
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
func (h *Handler) UpdateBlogHandler(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	defer r.Body.Close()

	var dto *blogPb.PostRequest
	if err = json.Unmarshal(body, &dto); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	customerId := r.Context().Value(middlewares.CustomerKey).(*authPb.CustomerDto).CustomerId
	if err := validateCustomer(customerId, dto.CustomerId); err != nil {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if _, err := h.blogClient.UpdatePost(context.Background(), dto); err != nil {
		if strings.Split(err.Error(), "desc = ")[1] == "customer not found" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "customer not found"})
			return
		}
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(204)

}

// @description -> Get a list of posts in the customer's portfolio
//
// @route -> /api/v1/blog/get/{skip}/
//
// @method -> GET
//
// @body -> an empty one
//
// @response 200
//
//	  {
//	  Blog: [20]blogPb.BlogItem{
//	  {
//	  Id: uint64
//	  Title: string
//	  Description: string
//	  Likes: uint64
//	  Dislikes: uint64
//	  CreatedAt: string
//	  UpdatedAt: string
//	  Tags: []string
//	  Images: []string
//	  },
//	  {
//	  Id: uint64
//	  Title: string
//	  Description: string
//	  Likes: uint64
//	  Dislikes: uint64
//	  CreatedAt: string
//	  UpdatedAt: string
//	  Tags: []string
//	  Images: []string
//	  },
//	  }
//	  }
//
//	@response 400 {object} ErrorResponse {error: err text}
//
//	@response 401 {object} ErrorResponse {error: err text}
//
//	@response 403 {object} ErrorResponse {error: err text}
//
//	@response 500 {object} ErrorResponse {error: err text}
func (h *Handler) GetBlogHandler(w http.ResponseWriter, r *http.Request) {

	customerId := r.Context().Value(middlewares.CustomerKey).(*authPb.CustomerDto).CustomerId
	skip, err := strconv.ParseUint(r.PathValue("skip"), 10, 32)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	blog, err := h.cacheService.GetBlog(customerId)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if blog != nil {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(blog); err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		return
	} else {

		dto := &blogPb.GetBlogRequest{
			CustomerId: customerId,
			Skip:       uint32(skip),
		}

		response, err := h.blogClient.GetBlog(context.Background(), dto)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		h.cacheService.SetBlog(customerId, response)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

	}
}

func (h *Handler) DeleteBlogHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(204)
}

func (h *Handler) SetReactionHandler(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(202)
}
