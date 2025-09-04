package routes

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	authPb "github.com/noo8xl/anvil-api/main/auth"
	notificationPb "github.com/noo8xl/anvil-api/main/notifications"
	offersPb "github.com/noo8xl/anvil-api/main/offers"
	profilePb "github.com/noo8xl/anvil-api/main/profile"
	"github.com/noo8xl/anvil-gateway/middlewares"
)

// @description -> Create a new offer
//
// @route -> /api/v1/offers/create/
//
// @method -> POST
//
// @body -> body should follow the following structure:
//
//	{
//		OfferId: uint64 -> can be omitted here
//		Title: string
//		Description: string
//		Location: string -> use an empty string if not available
//		Price: uint64
//		Status: string
//		ExpireAt: string
//		PostedBy: uint64
//		TagsList: []string
//		CategoriesList: []string
//	}
//
// @response 201
//
// @response 400 {error: err text}
//
// @response 401 {error: err text}
//
// @response 403 {error: err text}
//
// @response 500 {error: err text}
func (h *Handler) CreateOfferHandler(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	defer r.Body.Close()

	var dto *offersPb.OfferRequest
	if err := json.Unmarshal(body, &dto); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	customerId := r.Context().Value(middlewares.CustomerKey).(*authPb.CustomerDto).CustomerId
	if err := validateCustomer(customerId, dto.PostedBy); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	_, err = h.offersClient.CreateOffer(context.Background(), dto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// @description -> Update an offer
//
// @route -> /api/v1/offers/update/
//
// @method -> POST
//
// @body -> body should follow the following structure:
//
//	{
//		OfferId: uint64
//		Title: string
//		Description: string
//		Location: string
//		Price: uint64
//		Status: string
//		ExpireAt: string
//		PostedBy: uint64
//		TagsList: []string
//		CategoriesList: []string
//	}
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
func (h *Handler) UpdateOfferHandler(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	defer r.Body.Close()

	var dto *offersPb.OfferRequest
	if err := json.Unmarshal(body, &dto); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	customerId := r.Context().Value(middlewares.CustomerKey).(*authPb.CustomerDto).CustomerId
	if err := validateCustomer(customerId, dto.PostedBy); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	offer, err := h.cacheService.GetOfferDetails(dto.OfferId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if offer != nil {
		h.cacheService.ClearOfferDetails(dto.OfferId)
	}

	_, err = h.offersClient.UpdateOffer(context.Background(), dto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @description -> Get offer details by offer id
//
// @route -> /api/v1/offers/get-offer-details/{offerId}/
//
// @method -> GET
//
// @body -> an empty one
//
// @response 200 {object} GetOfferDetailsResponse
//
//	{
//	Offer: {
//	OfferId: uint64
//	PostedBy: uint64
//	AppliedBy: []uint64
//	ApprovedFor: uint64
//	ExpireAt: string
//	Details: {
//	Title: string
//	Description: string
//	Status: string
//	Price: uint64
//	Location: string
//	CreatedAt: string
//	UpdatedAt: string
//	ExpireAt: string
//	TagsList: []string
//	CategoriesList: []string
//	}
//	}
//
//	}
//
// @response 400 {object} ErrorResponse {error: err text}
//
// @response 401 {object} ErrorResponse {error: err text}
//
// @response 403 {object} ErrorResponse {error: err text}
//
// @response 500 {object} ErrorResponse {error: err text}
func (h *Handler) GetOfferDetailsHandler(w http.ResponseWriter, r *http.Request) {

	offerId, err := strconv.ParseUint(r.PathValue("offerId"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	offer, err := h.cacheService.GetOfferDetails(offerId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if offer != nil {
		w.WriteHeader(http.StatusNotModified)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(offer)
		return
	}

	payload := &offersPb.GetOfferDetailsRequest{
		OfferId: offerId,
	}

	offer, err = h.offersClient.GetOfferDetails(context.Background(), payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	h.cacheService.SetOfferDetails(offerId, offer)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(offer)
}

// @description -> Get a list of offers
//
// @route -> /api/v1/offers/get-offers-list/
//
// @method -> POST
//
// @body -> body should follow the following structure:
//
//	{
//		Skip: uint64
//		CustomerId: uint64
//		MinPrice: uint64
//		MaxPrice: uint64
//		DateFrom: string
//		DateTo: string
//		Location: string
//		TagsList: []string
//		CategoriesList: []string
//	}
//
// @response 200 {object} GetOffersListResponse
//
//	{
//		CardList: [20]offersPb.OfferShortCard{
//			OfferId: uint64
//			PostedBy: uint64
//			Title: string
//			Price: uint64
//			CreatedAt: string
//			ExpireAt: string
//			TagsList: []string
//			CategoriesList: []string
//		}
//	}
//
// @response 400 {object} {error: err text}
//
// @response 401 {object} {error: err text}
//
// @response 403 {object} {error: err text}
//
// @response 500 {object} {error: err text}
func (h *Handler) GetOffersListHandler(w http.ResponseWriter, r *http.Request) {

	customerId := r.Context().Value(middlewares.CustomerKey).(*authPb.CustomerDto).CustomerId

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	defer r.Body.Close()

	var filter *offersPb.OffersFilter
	err = json.Unmarshal(body, &filter)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	filter.CustomerId = customerId
	offers, err := h.offersClient.GetOffersList(context.Background(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(offers)
}

// @description -> handles requests to get a list of offers posted by the authenticated user.
//
// @route -> /api/v1/offers/get-my-offers/
//
// @method -> POST
//
// @body -> empty
//
// @response 200 {object}
//
//	{
//	  "CardList": [
//	    {
//	      "OfferId": "uint64",
//	      "PostedBy": "uint64",
//	      "Title": "string",
//	      "Price": "float64",
//	      "CreatedAt": "string",
//	      "ExpireAt": "string",
//	      "TagsList": ["string"],
//	      "CategoriesList": ["string"]
//	    }
//	  ]
//	}
//
// @response 400 {object} ErrorResponse
//
// @response 401 {object} ErrorResponse
//
// @response 403 {object} ErrorResponse
//
// @response 500 {object} ErrorResponse
func (h *Handler) GetMyOffersHandler(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	defer r.Body.Close()

	var filter *offersPb.GetMyOffersRequest
	err = json.Unmarshal(body, &filter)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	customerId := r.Context().Value(middlewares.CustomerKey).(*authPb.CustomerDto).CustomerId
	if err := validateCustomer(customerId, filter.CustomerId); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	offers, err := h.offersClient.GetMyOffers(context.Background(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(offers)
}

// @description -> Delete an offer by offerId
//
// @route -> /api/v1/offers/delete/{offerId}/
//
// @method -> DELETE
//
// @body -> an empty one
//
// @response 204
//
// @response 400 {object} ErrorResponse
//
// @response 401 {object} ErrorResponse
//
// @response 403 {object} ErrorResponse
//
// @response 500 {object} ErrorResponse
func (h *Handler) DeleteOfferHandler(w http.ResponseWriter, r *http.Request) {

	offerId, err := strconv.ParseUint(r.PathValue("offerId"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	payload := &offersPb.DeleteOfferRequest{
		OfferId: offerId,
	}

	_, err = h.offersClient.DeleteOffer(context.Background(), payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	h.cacheService.ClearOfferDetails(offerId)
	h.cacheService.ClearApplicantsList(offerId)

	w.WriteHeader(http.StatusNoContent)
}

// ####################### applications area #######################

// @description -> Apply to an offer
//
// @route -> /api/v1/offers/apply/
//
// @method -> POST
//
// @body -> body should follow the following structure:
//
//	{
//		OfferId: uint64
//		CustomerId: uint64
//		Title: string
//		Body: string
//	}
//
// @response 200
//
// @response 400 {object} ErrorResponse
//
// @response 401 {object} ErrorResponse
//
// @response 403 {object} ErrorResponse
//
// @response 500 {object} ErrorResponse
func (h *Handler) ApplyToTheOfferHandler(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	defer r.Body.Close()

	var dto *offersPb.ApplyToTheOfferRequest
	err = json.Unmarshal(body, &dto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	_, err = h.offersClient.ApplyToTheOffer(context.Background(), dto)
	if err != nil {
		if strings.Contains(err.Error(), "offer not found") {
			json.NewEncoder(w).Encode(map[string]string{"error": "Offer not found"})
			w.WriteHeader(http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	payload := &offersPb.GetOfferDetailsRequest{
		OfferId: dto.OfferId,
	}

	offer, err := h.offersClient.GetOfferDetails(context.Background(), payload)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	notificationPayload := &notificationPb.CreateNotificationRequest{
		CustomerId: offer.PostedBy,
		Title:      "New Application",
		Body:       "You have a new application to the offer. To view it, please go to the offers section.",
		Area:       "offers",
	}

	_, err = h.notificationsClient.CreateNotification(context.Background(), notificationPayload)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @description -> Get a list of applicants for an offer by offer id
//
// @route -> /api/v1/offers/get-applicants-list/{offerId}/{skip}/
//
// @method -> GET
//
// @body -> an empty one
//
// @response 200:
//
//	{
//	    Applicants: [20]offersPb.ApplicantDto{
//	    OfferId: uint64
//	    CustomerId: uint64
//	    AppliedAt: string
//	    Title: string
//	    Body: string
//	    Customer: profilePb.PublicProfileDto {
//	    Id: uint64
//	    Email: string
//	    Name: string
//	    Avatar: string
//	    Title: string
//	    Description: string
//	    Tags: []string
//	    }
//	    }
//	    Total: uint64
//	    }
//
// @response 400 {object} ErrorResponse
//
// @response 401 {object} ErrorResponse
//
// @response 403 {object} ErrorResponse
//
// @response 500 {object} ErrorResponse
func (h *Handler) GetApplicantsListHandler(w http.ResponseWriter, r *http.Request) {

	offerId, err := strconv.ParseUint(r.PathValue("offerId"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	skip, err := strconv.ParseUint(r.PathValue("skip"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	applicantsList, err := h.cacheService.GetApplicantsList(offerId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if applicantsList != nil {
		w.WriteHeader(http.StatusNotModified)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(applicantsList)
		return
	}

	payload := &offersPb.GetApplicantsListRequest{
		OfferId: offerId,
		Skip:    uint32(skip),
	}

	applicantsList, err = h.offersClient.GetApplicantsList(context.Background(), payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// get customer profile data to applicants applicantsList.Applicant.Customer
	for i, applicant := range applicantsList.Applicant {
		profilePayload := &profilePb.GetPublicProfileRequest{
			CustomerId: applicant.CustomerId,
		}

		applicantCard, err := h.profileClient.GetPublicProfile(context.Background(), profilePayload)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		applicantsList.Applicant[i].Customer = applicantCard
	}

	h.cacheService.SetApplicantsList(offerId, applicantsList)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(applicantsList)

}
