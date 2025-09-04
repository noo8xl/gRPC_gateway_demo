package routes

import "net/http"

func (h *AdminHandler) RegisterAdminRoutes(mux *http.ServeMux) {
	// mux.HandleFunc("GET /admin/orders/get-orders-list/{skip}/", h.GetOrdersListHandler)

	// mux.HandleFunc("POST /admin/notifications/create-notification/", h.CreateNotificationHandler)

	// mux.HandleFunc("POST /profile/create/", h.CreateCustomer) // TODO: admin permission only
}
