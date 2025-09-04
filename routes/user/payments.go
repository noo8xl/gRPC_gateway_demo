package routes

import (
	"net/http"
)

// description: Buy internal coins to use in the app
//
// @method -> POST
//
// @body -> body should follow the following structure:
//   - not implemented yet
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
func (h *Handler) BuyInternalCoinsHandler(w http.ResponseWriter, r *http.Request) {

	// customerId := r.Context().Value(middlewares.CustomerKey).(*authPb.CustomerDto).CustomerId

	w.WriteHeader(http.StatusOK)
}
