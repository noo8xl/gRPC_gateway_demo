package middlewares

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	authPb "github.com/noo8xl/anvil-api/main/auth"
)

type contextKey string

const CustomerKey contextKey = "customerDto"

func AuthMiddleware(next http.Handler, authClient authPb.AuthServiceClient) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// authless routes
		if strings.Contains(r.URL.Path, "/auth/") || strings.Contains(r.URL.Path, "/health") {
			next.ServeHTTP(w, r)
			return
		}

		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		t := strings.TrimSuffix(strings.Split(token, "Bearer ")[1], `"`)

		customer, err := authClient.ValidateToken(r.Context(), &authPb.ValidateTokenRequest{
			Token: t,
		})

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
			return
		}

		// admin routes guard
		if strings.Contains(r.URL.Path, "/admin/") {
			if customer.Role != "ADMIN" && customer.Role != "SUPERVISOR" {
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Forbidden",
				})
				return
			}
			next.ServeHTTP(w, r)
			return
		}

		// Use custom key type instead of string
		ctx := context.WithValue(r.Context(), CustomerKey, customer)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
