package routes

import (
	"net/http"
)

func (h *Handler) RegisterAuthRoutes(mux *http.ServeMux) {

	mux.HandleFunc("POST /api/v1/auth/sign-up/", h.HandleAuthSignUp)
	mux.HandleFunc("POST /api/v1/auth/sign-in/", h.HandleAuthSignIn)
	mux.HandleFunc("PATCH /api/v1/auth/forgot-password/{email}/", h.HandleAuthForgotPwd)

}

func (h *Handler) RegisterBlogRoutes(mux *http.ServeMux) {

	mux.HandleFunc("POST /api/v1/blog/create/", h.CreateBlogHandler)
	mux.HandleFunc("POST /api/v1/blog/update/", h.UpdateBlogHandler)
	mux.HandleFunc("GET /api/v1/blog/get/{skip}/", h.GetBlogHandler)
	mux.HandleFunc("DELETE /api/v1/blog/delete/{postId}/", h.DeleteBlogHandler)
	mux.HandleFunc("PATCH /api/v1/blog/set-reaction/", h.SetReactionHandler)

}

func (h *Handler) RegisterProfileRoutes(mux *http.ServeMux) {

	// profile -> base interactions
	mux.HandleFunc("GET /api/v1/profile/get/", h.GetCustomerProfileHandler)
	mux.HandleFunc("POST /api/v1/profile/update/", h.UpdateCustomerProfileHandler)
	mux.HandleFunc("POST /api/v1/profile/fill/", h.FillProfileHandler)
	mux.HandleFunc("GET /api/v1/profile/get-public-profile/", h.GetPublicProfileHandler)
	mux.HandleFunc("POST /api/v1/profile/report/", h.ReportCustomerHandler)

	// profile -> kyc area  TODO: not available in the MVP version
	mux.HandleFunc("POST /api/v1/profile/kyc/create/", h.CreateCustomeKycHandler)
	mux.HandleFunc("POST /api/v1/profile/kyc/update/", h.UpdateCustomeKycHandler)
	mux.HandleFunc("GET /api/v1/profile/kyc/get/", h.GetCustomerKycHandler)

	// profile -> security area
	mux.HandleFunc("PATCH /api/v1/profile/security/change-two-step-status/", h.ChangeTwoStepStatusHandler)
	mux.HandleFunc("PATCH /api/v1/profile/security/update/change-password/", h.ChangePasswordHandler)
	mux.HandleFunc("PATCH /api/v1/profile/security/update/change-email/", h.ChangeCustomerEmailHandler)
}

func (h *Handler) RegisterOffersRoutes(mux *http.ServeMux) {

	// offers -> base interactions
	mux.HandleFunc("POST /api/v1/offers/create/", h.CreateOfferHandler)
	mux.HandleFunc("POST /api/v1/offers/update/", h.UpdateOfferHandler)
	mux.HandleFunc("POST /api/v1/offers/get-offers-list/", h.GetOffersListHandler)
	mux.HandleFunc("POST /api/v1/offers/get-my-offers/", h.GetMyOffersHandler)
	mux.HandleFunc("DELETE /api/v1/offers/delete/{offerId}/", h.DeleteOfferHandler)
	mux.HandleFunc("GET /api/v1/offers/get-offer-details/{offerId}/", h.GetOfferDetailsHandler)

	// offers -> applicants interactions
	mux.HandleFunc("POST /api/v1/offers/apply/", h.ApplyToTheOfferHandler)
	mux.HandleFunc("GET /api/v1/offers/get-applicants-list/{offerId}/{skip}/", h.GetApplicantsListHandler)
}

func (h *Handler) RegisterOrdersRoutes(mux *http.ServeMux) {

	// orders -> base interactions
	mux.HandleFunc("POST /api/v1/orders/create/", h.CreateOrderHandler)
	mux.HandleFunc("PUT /api/v1/orders/update/", h.UpdateOrderHandler)

	mux.HandleFunc("POST /api/v1/orders/get-orders-list-by-filter/", h.GetOrdersListByFilterHandler)
	mux.HandleFunc("GET /api/v1/orders/get-order-details/{orderId}/", h.GetOrderDetailsHandler)
	mux.HandleFunc("DELETE /api/v1/orders/delete/{orderId}/", h.DeleteOrderHandler)

	// orders -> requests list (as an applicant)
	mux.HandleFunc("GET /api/v1/orders/get-orders-requests-list/{skip}/", h.GetOrdersRequestsListByApplicantIdHandler)

	// orders -> status interactions
	mux.HandleFunc("POST /api/v1/orders/apply/", h.ApplyToTheOrderHandler)
	mux.HandleFunc("POST /api/v1/orders/reject/", h.RejectAnOrderHandler)

	// orders -> compliance interactions
	mux.HandleFunc("POST /api/v1/orders/compliance/create/", h.CreateComplianceRequestHandler)
	mux.HandleFunc("POST /api/v1/orders/compliance/approve/", h.ComplianceApproveHandler)
	mux.HandleFunc("POST /api/v1/orders/compliance/reject/", h.RejectComplianceHandler)
	mux.HandleFunc("GET /api/v1/orders/compliance/get-list/{skip}/", h.GetComplianceRequestsListHandler)

}

func (h *Handler) RegisterReviewsRoutes(mux *http.ServeMux) {

	// reviews -> base interactions
	mux.HandleFunc("POST /api/v1/reviews/create/", h.CreateReviewHandler)
	mux.HandleFunc("PUT /api/v1/reviews/update/", h.UpdateReviewHandler)
	mux.HandleFunc("GET /api/v1/reviews/get-reviews-list/{skip}/", h.GetReviewsListHandler)
	mux.HandleFunc("DELETE /api/v1/reviews/delete/{reviewId}/", h.DeleteReviewHandler)

	// reviews -> comments interactions
	mux.HandleFunc("POST /api/v1/reviews/add-review-comment/", h.AddReviewCommentHandler)
	mux.HandleFunc("GET /api/v1/reviews/get-review-comments-list/{reviewId}/{skip}/", h.GetReviewCommentsListHandler)
	mux.HandleFunc("POST /api/v1/reviews/set-review-reaction/", h.SetReviewReactionHandler)

}

func (h *Handler) RegisterNotificationsRoutes(mux *http.ServeMux) {

	mux.HandleFunc("GET /api/v1/notifications/get-notifications-list/{skip}/", h.GetNotificationsListHandler)
	mux.HandleFunc("DELETE /api/v1/notifications/delete-notification/{notificationId}/", h.DeleteNotificationHandler)
	mux.HandleFunc("DELETE /api/v1/notifications/clear-notifications/", h.ClearNotificationsHandler)

}

func (h *Handler) RegisterChatRoutes(mux *http.ServeMux) {

	// mux.Handle("POST /api/v1/chat/create-chat/", middlewares.AuthMiddleware(http.HandlerFunc(h.CreateChatHandler)))
	// mux.Handle("GET /api/v1/chat/get-chats-list/{customerId}/{skip}/", middlewares.AuthMiddleware(http.HandlerFunc(h.GetChatsListHandler)))
	// mux.Handle("GET /api/v1/chat/get-chat-messages-list/{chatId}/{skip}/", middlewares.AuthMiddleware(http.HandlerFunc(h.GetChatMessagesListHandler)))

}

func (h *Handler) RegisterPromotionsRoutes(mux *http.ServeMux) {

	// mux.HandleFunc("POST /api/v1/promotions/create/", h.CreatePromotionHandler)
	// mux.HandleFunc("PUT /api/v1/promotions/update/", h.UpdatePromotionHandler)
	// mux.HandleFunc("GET /api/v1/promotions/get-promotions-list/{skip}/", h.GetPromotionsListHandler)

}

func (h *Handler) RegisterPaymentsRoutes(mux *http.ServeMux) {

	// payments -> base interactions
	// NOTE: not available for the MVP version

	// mux.HandleFunc("POST /api/v1/payments/buy-internal-coins/", h.BuyInternalCoinsHandler)
	// mux.HandleFunc("POST /api/v1/payments/create/", h.CreatePaymentHandler)
	// mux.HandleFunc("PUT /api/v1/payments/update/", h.UpdatePaymentHandler)
	// mux.HandleFunc("GET /api/v1/payments/get-payments-list/{customerId}/{skip}/", h.GetPaymentsListHandler)
	// mux.HandleFunc("DELETE /api/v1/payments/delete/{paymentId}/", h.DeletePaymentHandler)

}
