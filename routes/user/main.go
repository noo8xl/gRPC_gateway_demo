package routes

import (
	"errors"

	authPb "github.com/noo8xl/anvil-api/main/auth"
	blogPb "github.com/noo8xl/anvil-api/main/blog"
	notificationsPb "github.com/noo8xl/anvil-api/main/notifications"
	offersPb "github.com/noo8xl/anvil-api/main/offers"
	ordersPb "github.com/noo8xl/anvil-api/main/orders"
	paymentsPb "github.com/noo8xl/anvil-api/main/payments"
	profilePb "github.com/noo8xl/anvil-api/main/profile"

	// promotionsPb "github.com/noo8xl/anvil-api/main/promotions"
	reviewsPb "github.com/noo8xl/anvil-api/main/reviews"

	"github.com/noo8xl/anvil-gateway/cache"
)

type Handler struct {
	authClient          authPb.AuthServiceClient
	profileClient       profilePb.ProfileServiceClient
	ordersClient        ordersPb.OrdersServiceClient
	reviewsClient       reviewsPb.ReviewsServiceClient
	blogClient          blogPb.BlogServiceClient
	paymentsClient      paymentsPb.PaymentsServiceClient
	offersClient        offersPb.OffersServiceClient
	notificationsClient notificationsPb.NotificationsServiceClient
	// promotionsClient    promotionsPb.PromotionsServiceClient
	cacheService *cache.CacheService
}

func InitHandler(
	authClient authPb.AuthServiceClient,
	profileClient profilePb.ProfileServiceClient,
	ordersClient ordersPb.OrdersServiceClient,
	reviewsClient reviewsPb.ReviewsServiceClient,
	blogClient blogPb.BlogServiceClient,
	paymentsClient paymentsPb.PaymentsServiceClient,
	offersClient offersPb.OffersServiceClient,
	notificationsClient notificationsPb.NotificationsServiceClient,
	// promotionsClient promotionsPb.PromotionsServiceClient,
) *Handler {
	cs := cache.InitCacheService()
	return &Handler{
		authClient:          authClient,
		profileClient:       profileClient,
		ordersClient:        ordersClient,
		reviewsClient:       reviewsClient,
		blogClient:          blogClient,
		paymentsClient:      paymentsClient,
		offersClient:        offersClient,
		notificationsClient: notificationsClient,
		// promotionsClient:    promotionsClient,
		cacheService: cs,
	}
}

// ############################################################
// ############################################################
// ############################################################

// validateCustomer -> validate customer id from request header and payload
func validateCustomer(headerId uint64, payloadId uint64) error {
	if headerId != payloadId {
		return errors.New("forbidden: not allowed to access this resource")
	}
	return nil
}
