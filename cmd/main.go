package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	adminRoutes "github.com/noo8xl/anvil-gateway/routes/admin"
	userRoutes "github.com/noo8xl/anvil-gateway/routes/user"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	authPb "github.com/noo8xl/anvil-api/main/auth"
	blogPb "github.com/noo8xl/anvil-api/main/blog"
	notificationsPb "github.com/noo8xl/anvil-api/main/notifications"
	offersPb "github.com/noo8xl/anvil-api/main/offers"
	ordersPb "github.com/noo8xl/anvil-api/main/orders"
	profilePb "github.com/noo8xl/anvil-api/main/profile"
	reviewsPb "github.com/noo8xl/anvil-api/main/reviews"

	"github.com/noo8xl/anvil-common/exceptions"
	"github.com/noo8xl/anvil-gateway/config"
	"github.com/noo8xl/anvil-gateway/middlewares"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	httpServerAddress := config.GetServerAddressByServiceName("gateway")

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	cfg := config.LoadConfig()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	clients := initClients(cfg, logger)
	defer func() {
		for _, client := range clients {
			if conn, ok := client.(*grpc.ClientConn); ok {
				if err := conn.Close(); err != nil {
					logger.Error("failed to close gRPC connection", zap.Error(err))
				}
			}
		}
	}()

	mux := http.NewServeMux()

	mux.Handle("/health/", healthCheckHandler(clients))
	mux.Handle("/metrics/", promhttp.Handler())

	userHandler := userRoutes.InitHandler(
		clients["auth"].(authPb.AuthServiceClient),
		clients["profile"].(profilePb.ProfileServiceClient),
		clients["orders"].(ordersPb.OrdersServiceClient),
		clients["reviews"].(reviewsPb.ReviewsServiceClient),
		clients["blog"].(blogPb.BlogServiceClient),
		nil,
		// clients["payments"].(paymentsPb.PaymentsServiceClient),
		clients["offers"].(offersPb.OffersServiceClient),
		clients["notifications"].(notificationsPb.NotificationsServiceClient),
	)

	adminHandler := adminRoutes.InitAdminHandler(
		clients["auth"].(authPb.AuthServiceClient),
		clients["profile"].(profilePb.ProfileServiceClient),
		clients["orders"].(ordersPb.OrdersServiceClient),
		clients["reviews"].(reviewsPb.ReviewsServiceClient),
		clients["blog"].(blogPb.BlogServiceClient),
		nil,
		// clients["payments"].(paymentsPb.PaymentsServiceClient),
		clients["offers"].(offersPb.OffersServiceClient),
		clients["notifications"].(notificationsPb.NotificationsServiceClient),
	)

	if err := registerRoutes(mux, userHandler, adminHandler); err != nil {
		logger.Fatal("failed to register routes", zap.Error(err))
	}

	// Apply middlewares in order
	serverHandler := middlewares.Logger(
		middlewares.MetricsMiddleware(
			middlewares.AuthMiddleware(mux, clients["auth"].(authPb.AuthServiceClient)),
		),
	)
	server := config.GetServerConfig(httpServerAddress, serverHandler)

	go func() {
		log.Printf("gateway server is running successfully on %s", httpServerAddress)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			msg := fmt.Errorf("failed to start http server: %v", err)
			exceptions.HandleAnException(msg)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("shutting down gracefully...")
	if err := server.Shutdown(shutdownCtx); err != nil {
		msg := fmt.Errorf("server forced to shutdown: %v", err)
		exceptions.HandleAnException(msg)
	}

	log.Println("server exited properly")
}

func initClients(cfg *config.Config, logger *zap.Logger) map[string]any {
	clients := make(map[string]any)

	env := os.Getenv("GO_ENV")
	retryOpts := []grpc_retry.CallOption{
		grpc_retry.WithMax(uint(cfg.RetryConfig.MaxAttempts)),
		grpc_retry.WithBackoff(grpc_retry.BackoffLinear(cfg.RetryConfig.Backoff)),
	}

	for service, address := range cfg.Services {
		var conn *grpc.ClientConn
		var err error

		if env == "production" {
			conn, err = initClientWithRetry(address, retryOpts, cfg.TLS)
		} else {
			conn, err = initClient(address, retryOpts)
		}

		if err != nil {
			logger.Fatal("failed to initialize client",
				zap.String("service", service),
				zap.Error(err))
		}

		client := initServiceClient(service, conn)
		clients[service] = client

	}

	return clients
}

func initClientWithRetry(address string, retryOpts []grpc_retry.CallOption, tlsConfig config.TLSConfig) (*grpc.ClientConn, error) {
	var creds credentials.TransportCredentials

	if os.Getenv("GO_ENV") == "production" {
		if tlsConfig.CertFile == "" || tlsConfig.KeyFile == "" {
			return nil, fmt.Errorf("TLS cert file or key file is not set")
		}

		cert, err := tls.LoadX509KeyPair(tlsConfig.CertFile, tlsConfig.KeyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load TLS cert: %v", err)
		}
		creds = credentials.NewServerTLSFromCert(&cert)
	} else {
		creds = insecure.NewCredentials()
	}

	return grpc.NewClient(
		address,
		grpc.WithTransportCredentials(creds),
		grpc.WithChainUnaryInterceptor(
			grpc_retry.UnaryClientInterceptor(retryOpts...),
		),
	)
}

// func closeClients(clients map[string]interface{}) {
// 	for _, client := range clients {
// 		if conn, ok := client.(*grpc.ClientConn); ok {
// 			conn.Close()
// 		}
// 	}
// }

func registerRoutes(mux *http.ServeMux, userHandler *userRoutes.Handler, adminHandler *adminRoutes.AdminHandler) error {
	// register user routes
	userHandler.RegisterAuthRoutes(mux)
	userHandler.RegisterProfileRoutes(mux)
	userHandler.RegisterBlogRoutes(mux)
	userHandler.RegisterOffersRoutes(mux)
	userHandler.RegisterOrdersRoutes(mux)
	userHandler.RegisterReviewsRoutes(mux)
	userHandler.RegisterPaymentsRoutes(mux)
	userHandler.RegisterNotificationsRoutes(mux)

	// register admin routes
	adminHandler.RegisterAdminRoutes(mux)

	return nil
}

func healthCheckHandler(clients map[string]any) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for service, client := range clients {
			if !isServiceHealthy(client) {
				http.Error(w, fmt.Sprintf("%s service unhealthy", service), http.StatusServiceUnavailable)
				return
			}
		}
		w.WriteHeader(http.StatusOK)
	})
}

func initClient(serverAddress string, retryOpts []grpc_retry.CallOption) (*grpc.ClientConn, error) {

	var creds credentials.TransportCredentials
	var err error
	creds = insecure.NewCredentials()

	conn, err := grpc.NewClient(
		serverAddress,
		grpc.WithTransportCredentials(creds),
		grpc.WithChainUnaryInterceptor(
			grpc_retry.UnaryClientInterceptor(retryOpts...),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create grpc client <serverAddress: %s>: %v", serverAddress, err)
	}

	return conn, nil
}

func initServiceClient(service string, conn *grpc.ClientConn) any {
	switch service {
	case "auth":
		return authPb.NewAuthServiceClient(conn)
	case "profile":
		return profilePb.NewProfileServiceClient(conn)
	case "orders":
		return ordersPb.NewOrdersServiceClient(conn)
	case "reviews":
		return reviewsPb.NewReviewsServiceClient(conn)
	case "blog":
		return blogPb.NewBlogServiceClient(conn)
	case "payments":
		return nil
	// 	return paymentsPb.NewPaymentsServiceClient(conn)
	case "offers":
		return offersPb.NewOffersServiceClient(conn)
	case "notifications":
		return notificationsPb.NewNotificationsServiceClient(conn)
	// case "promotions":
	// 	return promotionsPb.NewPromotionsServiceClient(conn)
	default:
		msg := fmt.Errorf("unknown service: %s", service)
		exceptions.HandleAnException(msg)
		return nil
	}
}

func isServiceHealthy(client interface{}) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch c := client.(type) {
	case authPb.AuthServiceClient:
		_, err := c.HealthCheck(ctx, &authPb.HealthCheckRequest{})
		return err == nil
	case profilePb.ProfileServiceClient:
		_, err := c.HealthCheck(ctx, &profilePb.HealthCheckRequest{})
		return err == nil
	case ordersPb.OrdersServiceClient:
		_, err := c.HealthCheck(ctx, &ordersPb.HealthCheckRequest{})
		return err == nil
	case reviewsPb.ReviewsServiceClient:
		_, err := c.HealthCheck(ctx, &reviewsPb.HealthCheckRequest{})
		return err == nil
	case blogPb.BlogServiceClient:
		_, err := c.HealthCheck(ctx, &blogPb.HealthCheckRequest{})
		return err == nil
	// case paymentsPb.PaymentsServiceClient:
	// 	_, err := c.HealthCheck(ctx, &paymentsPb.HealthCheckRequest{})
	// 	return err == nil
	case offersPb.OffersServiceClient:
		_, err := c.HealthCheck(ctx, &offersPb.HealthCheckRequest{})
		return err == nil
	case notificationsPb.NotificationsServiceClient:
		_, err := c.HealthCheck(ctx, &notificationsPb.HealthCheckRequest{})
		return err == nil
	default:
		return false
	}
}
