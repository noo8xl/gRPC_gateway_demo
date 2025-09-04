package routes

import (
	authPb "github.com/noo8xl/anvil-api/main/auth"
	blogPb "github.com/noo8xl/anvil-api/main/blog"
	notificationsPb "github.com/noo8xl/anvil-api/main/notifications"
	offersPb "github.com/noo8xl/anvil-api/main/offers"
	ordersPb "github.com/noo8xl/anvil-api/main/orders"
	paymentsPb "github.com/noo8xl/anvil-api/main/payments"
	profilePb "github.com/noo8xl/anvil-api/main/profile"

	// promotionsPb "github.com/noo8xl/anvil-api/main/promotions"
	reviewsPb "github.com/noo8xl/anvil-api/main/reviews"
)

type AdminHandler struct {
	authClient          authPb.AuthServiceClient
	profileClient       profilePb.ProfileServiceClient
	ordersClient        ordersPb.OrdersServiceClient
	reviewsClient       reviewsPb.ReviewsServiceClient
	blogClient          blogPb.BlogServiceClient
	paymentsClient      paymentsPb.PaymentsServiceClient
	offersClient        offersPb.OffersServiceClient
	notificationsClient notificationsPb.NotificationsServiceClient
	// promotionsClient    promotionsPb.PromotionsServiceClient
}

func InitAdminHandler(
	authClient authPb.AuthServiceClient,
	profileClient profilePb.ProfileServiceClient,
	ordersClient ordersPb.OrdersServiceClient,
	reviewsClient reviewsPb.ReviewsServiceClient,
	blogClient blogPb.BlogServiceClient,
	paymentsClient paymentsPb.PaymentsServiceClient,
	offersClient offersPb.OffersServiceClient,
	notificationsClient notificationsPb.NotificationsServiceClient,
	// promotionsClient promotionsPb.PromotionsServiceClient,
) *AdminHandler {
	return &AdminHandler{
		authClient:          authClient,
		profileClient:       profileClient,
		ordersClient:        ordersClient,
		reviewsClient:       reviewsClient,
		blogClient:          blogClient,
		paymentsClient:      paymentsClient,
		offersClient:        offersClient,
		notificationsClient: notificationsClient,
		// promotionsClient:    promotionsClient,
	}
}
