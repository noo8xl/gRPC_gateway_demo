package routes

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	authPb "github.com/noo8xl/anvil-api/main/auth"
	notificationPb "github.com/noo8xl/anvil-api/main/notifications"
	profilePb "github.com/noo8xl/anvil-api/main/profile"
	helpers "github.com/noo8xl/anvil-common/helpers"
)

// @description -> Sign up a new customer
//
// @route -> /api/v1/auth/sign-up
//
// @method -> POST
//
// @response 200 {object} profilePb.CreateCustomerResponse
//
//	{
//		Email: string
//		Name: string
//		Password: string
//	}
//
// @response 400 {object} ErrorResponse {error: err text}
//
// @response 401 {object} ErrorResponse {error: err text}
//
// @response 403 {object} ErrorResponse {error: err text}
//
// @response 500 {object} ErrorResponse {error: err text}
func (h *Handler) HandleAuthSignUp(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	defer r.Body.Close()

	var dto *profilePb.CreateCustomerRequest
	err = json.Unmarshal(body, &dto)
	if err != nil {
		if err.Error() == "rpc error: code = Unknown desc = customer already exists" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": strings.Split(err.Error(), "desc = ")[1]})
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	_, err = h.profileClient.CreateCustomer(context.Background(), dto)
	if err != nil {
		if err.Error() == "rpc error: code = Unknown desc = customer already exists" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "customer already exists"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// @description -> Sign in a customer
//
// @route -> /api/v1/auth/sign-in
//
// @method -> POST
//
// @response 200 {object} authPb.SignInResponse
//
//	{
//		Email: string
//		Password: string
//		TwoStepCode: string
//	}
//
// @response 400 {object} ErrorResponse {error: err text}
//
// @response 401 {object} ErrorResponse {error: err text}
//
// @response 403 {object} ErrorResponse {error: err text}
//
// @response 500 {object} ErrorResponse {error: err text}
func (h *Handler) HandleAuthSignIn(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	defer r.Body.Close()

	var dto *authPb.SignInRequest
	err = json.Unmarshal(body, &dto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	customer, err := h.authClient.GetCustomer(context.Background(), &authPb.GetCustomerByEmailRequest{Email: dto.Email})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if customer.IsTwoFa.Value && dto.TwoStepCode != "" {
		c, err := h.cacheService.Get2FACode(dto.Email)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		if c != dto.TwoStepCode {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid code"})
			return
		}
	}

	if customer.IsTwoFa.Value && dto.TwoStepCode == "" {
		code := helpers.GenerateRandomPassword(6)
		h.cacheService.Set2FACode(dto.Email, code)

		notificationDto := &notificationPb.SendEmailRequest{
			Email:   dto.Email,
			Subject: "Two-Step Verification",
			Body:    code,
		}

		_, err = h.notificationsClient.SendEmail(context.Background(), notificationDto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "two-step auth is enabled. code was sent"})
		return
	}

	response, err := h.authClient.SignIn(context.Background(), dto)
	if err != nil {
		if strings.Split(err.Error(), "desc = ")[1] == "invalid code" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid code"})
			return
		}

		if strings.Split(err.Error(), "desc = ")[1] == "two-step auth is enabled" {

			code := helpers.GenerateRandomPassword(6)
			h.cacheService.Set2FACode(dto.Email, code)

			notificationDto := &notificationPb.SendEmailRequest{
				Email:   dto.Email,
				Subject: "Two-Step Verification",
				Body:    code,
			}

			_, err = h.notificationsClient.SendEmail(context.Background(), notificationDto)
			if err != nil {
				log.Println("auth sign in err -> ", err)
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "two-step auth is enabled. code was sent"})
			return
		}

		if strings.Split(err.Error(), "desc = ")[1] == "invalid password" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid password"})
			return
		}

		if strings.Split(err.Error(), "desc = ")[1] == "customer not found" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "customer not found"})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

// @description -> Forgot password
//
// @route -> /api/v1/auth/forgot-password/{email}
//
// @method -> POST
//
// @response 202
//
// @response 400 {object} ErrorResponse {error: err text}
//
// @response 401 {object} ErrorResponse {error: err text}
//
// @response 403 {object} ErrorResponse {error: err text}
//
// @response 500 {object} ErrorResponse {error: err text}
func (h *Handler) HandleAuthForgotPwd(w http.ResponseWriter, r *http.Request) {

	if r.PathValue("email") == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "email is required"})
		return
	}

	response, err := h.authClient.GetCustomerPassword(context.Background(), &authPb.GetCustomerByEmailRequest{Email: r.PathValue("email")})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if response.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "- customer not found"})
		return
	}

	pwd := helpers.GenerateRandomPassword(12)
	hash, err := helpers.EncryptKey(pwd)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// log pwd and hash while debugging
	// log.Println("generated pwd -> ", pwd)
	// log.Println("encrypted pwd -> ", hash)

	changePasswordDto := &profilePb.ChangePasswordRequest{
		CustomerId:  response.CustomerId,
		OldPassword: response.Password,
		NewPassword: hash,
	}

	notificationDto := &notificationPb.SendEmailRequest{
		Email:   response.Email,
		Subject: "Reset Password",
		Body:    pwd,
	}

	if _, err = h.profileClient.ChangePassword(context.Background(), changePasswordDto); err != nil {
		if err.Error() == "got a wrong old password" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "got a wrong old password"})
			return
		}

		if strings.Contains(err.Error(), "sql: no rows in result set") {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "-- customer not found"})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	_, err = h.notificationsClient.SendEmail(context.Background(), notificationDto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
}
