package routes

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	authPb "github.com/noo8xl/anvil-api/main/auth"
	notificationPb "github.com/noo8xl/anvil-api/main/notifications"
	profilePb "github.com/noo8xl/anvil-api/main/profile"
	"github.com/noo8xl/anvil-common/exceptions"
	helpers "github.com/noo8xl/anvil-common/helpers"
	"github.com/noo8xl/anvil-gateway/middlewares"
)

// @description -> Change customer email
//
// @route -> /api/v1/profile/security/update/change-email/
//
// @method -> PATCH
//
// @body -> body should follow the following structure:
//
//	  {
//			CustomerId: uint64
//			OldEmail: string
//			NewEmail: string
//	  }
//
// @response 200
//
// @response 400 {object} {error: err text}
//
// @response 401 {object} {error: err text}
//
// @response 403 {object} {error: err text}
//
// @response 500 {object} {error: err text}
func (h *Handler) ChangeCustomerEmailHandler(w http.ResponseWriter, r *http.Request) {

	email := r.Context().Value(middlewares.CustomerKey).(*authPb.CustomerDto).Email
	if email == "" {
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized: customer not found in context"})
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	defer r.Body.Close()

	var dto *profilePb.ChangeCustomerEmailRequest
	err = json.Unmarshal(body, &dto)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	customerId := r.Context().Value(middlewares.CustomerKey).(*authPb.CustomerDto).CustomerId
	if err := validateCustomer(customerId, dto.CustomerId); err != nil {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if dto.OldEmail != email {
		w.WriteHeader(401)
		err := exceptions.HandleAnException(errors.New("** attack detected: unknown customer try to change other customer email **"))
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	_, err = h.profileClient.ChangeCustomerEmail(context.Background(), dto)
	if err != nil {
		if strings.Contains(err.Error(), "error: email is not available to be in use") {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(map[string]string{"error": "error: email is not available to be in use"})
			return
		}

		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(200)
}

// @description -> Change customer password
//
// @route -> /api/v1/profile/security/update/change-password/
//
// @method -> PATCH
//
// @body -> body should follow the following structure:
//
//	  {
//			CustomerId: uint64
//			OldPassword: string
//			NewPassword: string
//	  }
//
// @response 200
//
// @response 400 {object} {error: err text}
//
// @response 401 {object} {error: err text}
//
// @response 403 {object} {error: err text}
//
// @response 500 {object} {error: err text}
func (h *Handler) ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	defer r.Body.Close()

	var dto *profilePb.ChangePasswordRequest
	err = json.Unmarshal(body, &dto)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	customerId := r.Context().Value(middlewares.CustomerKey).(*authPb.CustomerDto).CustomerId
	if err := validateCustomer(customerId, dto.CustomerId); err != nil {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	_, err = h.profileClient.ChangePassword(context.Background(), dto)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(200)
}

// @description -> Change customer two step status
//
// @route -> /api/v1/profile/security/update/change-two-step-status/
//
// @method -> PATCH
//
// @body -> body should follow the following structure:
//
//	  {
//			CustomerId: uint64
//			IsEnabled: bool
//			Email: string
//			Code: string
//	  }
//
// @response 200
//
// @response 400 {object} {error: err text}
//
// @response 401 {object} {error: err text}
//
// @response 403 {object} {error: err text}
//
// @response 500 {object} {error: err text}
func (h *Handler) ChangeTwoStepStatusHandler(w http.ResponseWriter, r *http.Request) {

	customerId := r.Context().Value(middlewares.CustomerKey).(*authPb.CustomerDto).CustomerId
	email := r.Context().Value(middlewares.CustomerKey).(*authPb.CustomerDto).Email

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	defer r.Body.Close()

	var dto *profilePb.ChangeTwoStepStatusRequest
	err = json.Unmarshal(body, &dto)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if err := validateCustomer(customerId, dto.CustomerId); err != nil {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if customerId != dto.CustomerId && email != dto.Email {
		w.WriteHeader(401)
		err := exceptions.HandleAnException(errors.New("** attack detected: unknown customer try to change other customer two step status **"))
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if !dto.IsEnabled {
		err = h.disableTwoStepHandler(dto)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
	} else {
		err = h.enableTwoStepHandler(dto)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
	}

	w.WriteHeader(200)
}

func (h *Handler) enableTwoStepHandler(dto *profilePb.ChangeTwoStepStatusRequest) error {

	if _, err := h.profileClient.ChangeTwoStepStatus(context.Background(), dto); err != nil {
		if strings.Contains(err.Error(), "customer not found") {
			return errors.New("customer not found")
		}
		if strings.Contains(err.Error(), "customer already has two step enabled") {
			return errors.New("customer already has two step enabled")
		}
		return err
	}
	return nil
}

func (h *Handler) disableTwoStepHandler(dto *profilePb.ChangeTwoStepStatusRequest) error {

	if dto.Code != "" {
		c, err := h.cacheService.Get2FACode(dto.Email)
		if err != nil {
			if err.Error() == "code not found" {
				return errors.New("code not found")
			}
			return err
		}

		if c != dto.Code {
			return errors.New("invalid code")
		}

		if _, err := h.profileClient.ChangeTwoStepStatus(context.Background(), dto); err != nil {
			if strings.Contains(err.Error(), "customer not found") {
				return errors.New("customer not found")
			}
			if strings.Contains(err.Error(), "customer already has two step enabled") {
				return errors.New("customer already has two step enabled")
			}
			return err
		}
		return nil
	} else {

		code := helpers.GenerateRandomPassword(12)
		h.cacheService.Set2FACode(dto.Email, code)

		_, err := h.notificationsClient.SendEmail(context.Background(), &notificationPb.SendEmailRequest{
			Email:   dto.Email,
			Subject: "Two-Step Verification",
			Body:    code,
		})
		if err != nil {
			return err
		}
	}

	h.cacheService.ClearPublicProfile(dto.CustomerId)
	h.cacheService.ClearCustomerProfile(dto.CustomerId)

	return nil
}
