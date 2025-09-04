# Anvil Gateway API Documentation

## Generated on: Fri May 23 03:32:10 PM CEST 2025

package routes // import "github.com/noo8xl/anvil-gateway/routes/user"



func (h *Handler) AddReviewCommentHandler(w http.ResponseWriter, r *http.Request)
    @description -> Add a comment to a review

    @route -> /api/v1/reviews/add-review-comment/

    @method -> POST

    @body -> comment body should follow the following structure:

        {
        	Comment: {
        		Comment: string
        		PostedBy: uint64
        		Title: string
        		Body: string
        		CreatedAt: string
        		UpdatedAt: string
        	}
        }

    @response 201

    @response 400 {object} ErrorResponse {error: err text}

    @response 401 {object} ErrorResponse {error: err text}

    @response 403 {object} ErrorResponse {error: err text}

    @response 500 {object} ErrorResponse {error: err text}

func (h *Handler) ApplyToTheOfferHandler(w http.ResponseWriter, r *http.Request)
    @description -> Apply to an offer

    @route -> /api/v1/offers/apply/

    @method -> POST

    @body -> body should follow the following structure:

        {
        	OfferId: uint64
        	CustomerId: uint64
        	Title: string
        	Body: string
        }

    @response 200

    @response 400 {object} ErrorResponse

    @response 401 {object} ErrorResponse

    @response 403 {object} ErrorResponse

    @response 500 {object} ErrorResponse

func (h *Handler) ApplyToTheOrderHandler(w http.ResponseWriter, r *http.Request)
    @description -> Apply to the order

    @route -> /api/v1/orders/apply/

    @method -> POST

    @body -> body should follow the following structure:

        {
        	orderBasics: {
        		CustomerId: uint64
        		ApplicantId: uint64
        		OfferId: uint64
        		OrderId: uint64
        	}
        }

    @response 200

    @response 400 {error: err text}

    @response 401 {error: err text}

    @response 403 {error: err text}

    @response 500 {error: err text}

func (h *Handler) BuyInternalCoinsHandler(w http.ResponseWriter, r *http.Request)
    description: Buy internal coins to use in the app

    @method -> POST

    @body -> body should follow the following structure:
      - not implemented yet

    @response 200

    @response 400 {error: err text}

    @response 401 {error: err text}

    @response 403 {error: err text}

    @response 500 {error: err text}

func (h *Handler) ChangeCustomerEmailHandler(w http.ResponseWriter, r *http.Request)
    @description -> Change customer email

    @route -> /api/v1/profile/security/update/change-email/

    @method -> PATCH

    @body -> body should follow the following structure:

          {
        		CustomerId: uint64
        		OldEmail: string
        		NewEmail: string
          }

    @response 200

    @response 400 {object} {error: err text}

    @response 401 {object} {error: err text}

    @response 403 {object} {error: err text}

    @response 500 {object} {error: err text}

func (h *Handler) ChangePasswordHandler(w http.ResponseWriter, r *http.Request)
    @description -> Change customer password

    @route -> /api/v1/profile/security/update/change-password/

    @method -> PATCH

    @body -> body should follow the following structure:

          {
        		CustomerId: uint64
        		OldPassword: string
        		NewPassword: string
          }

    @response 200

    @response 400 {object} {error: err text}

    @response 401 {object} {error: err text}

    @response 403 {object} {error: err text}

    @response 500 {object} {error: err text}

func (h *Handler) ChangeTwoStepStatusHandler(w http.ResponseWriter, r *http.Request)
    @description -> Change customer two step status

    @route -> /api/v1/profile/security/update/change-two-step-status/

    @method -> PATCH

    @body -> body should follow the following structure:

          {
        		CustomerId: uint64
        		IsEnabled: bool
        		Email: string
        		Code: string
          }

    @response 200

    @response 400 {object} {error: err text}

    @response 401 {object} {error: err text}

    @response 403 {object} {error: err text}

    @response 500 {object} {error: err text}

func (h *Handler) ClearNotificationsHandler(w http.ResponseWriter, r *http.Request)
    @description -> Clear all notifications

    @route -> /api/v1/notifications/clear-notifications/

    @method -> DELETE

    @body -> an empty one

    Response list:

      - @response 200

    @response 400 {object} ErrorResponse

    @response 401 {object} ErrorResponse

    @response 403 {object} ErrorResponse

    @response 500 {object} ErrorResponse

func (h *Handler) ComplianceApproveHandler(w http.ResponseWriter, r *http.Request)
    @description -> Approve a compliance request by offer owner

    @route -> /api/v1/orders/compliance/approve/

    @method -> POST

    @body -> body should follow the following structure:

        {
        	Id: uint64 -> can be omitted here
        	Round: uint32
        	Status: string
        	Claim: string
        	CreatedAt: string -> can be omitted here
        	OrderBasics: {
        		CustomerId: uint64
        		ApplicantId: uint64спать уже собрался
        	}
        }

    @response 200

    @response 400 {error: err text}

    @response 401 {error: err text}

    @response 403 {error: err text}

    @response 500 {error: err text}

func (h *Handler) CreateBlogHandler(w http.ResponseWriter, r *http.Request)
    @description -> Create a blog post

    @route -> /api/v1/blog/create/

    @method -> POST

    @body -> body should follow the following structure:

        {
        	Id: uint64 -> can be omitted here
        	CustomerId: uint64
        	Title: string
        	Description: string
        	TagList: []string
        	ImageList: []string
        }

    @response 201

    @response 400 {object} ErrorResponse {error: err text}

    @response 401 {object} ErrorResponse {error: err text}

    @response 403 {object} ErrorResponse {error: err text}

    @response 500 {object} ErrorResponse {error: err text}

func (h *Handler) CreateComplianceRequestHandler(w http.ResponseWriter, r *http.Request)
    @description -> Create a new compliance request if customer has done with
    the order / round

    @route -> /api/v1/orders/compliance/create/

    @method -> POST

    @body -> body should follow the following structure:

        {
        	Id: uint64 -> can be omitted here
        	Round: uint32
        	Status: string
        	Claim: string
        	CreatedAt: string -> can be omitted here
        	OrderBasics: {
        		CustomerId: uint64
        		ApplicantId: uint64
        		OrderId: uint64
        		OfferId: uint64
        	}
        }

    @response 200

    @response 400 {error: err text}

    @response 401 {error: err text}

    @response 403 {error: err text}

    @response 500 {error: err text}

func (h *Handler) CreateCustomeKycHandler(w http.ResponseWriter, r *http.Request)
    @description -> Create a customer kyc

    @route -> /api/v1/profile/kyc/create/

    @method -> POST

    @body -> not implemented yet

    @response 201 {object}

    @response 400 {object} ErrorResponse {error: err text}

    @response 401 {object} ErrorResponse {error: err text}

    @response 403 {object} ErrorResponse {error: err text}

    @response 500 {object} ErrorResponse {error: err text}

func (h *Handler) CreateCustomer(w http.ResponseWriter, r *http.Request)
    @description -> Create a customer

    @route -> this route is available from the client side without the admin
    permission *

    @method -> POST

    @response 201

    @response 400 {object} {error: err text}

    @response 401 {object} {error: err text}

    @response 403 {object} {error: err text}

    @response 500 {object} {error: err text}

func (h *Handler) CreateOfferHandler(w http.ResponseWriter, r *http.Request)
    @description -> Create a new offer

    @route -> /api/v1/offers/create/

    @method -> POST

    @body -> body should follow the following structure:

        {
        	OfferId: uint64 -> can be omitted here
        	Title: string
        	Description: string
        	Location: string -> use an empty string if not available
        	Price: uint64
        	Status: string
        	ExpireAt: string
        	PostedBy: uint64
        	TagsList: []string
        	CategoriesList: []string
        }

    @response 201

    @response 400 {error: err text}

    @response 401 {error: err text}

    @response 403 {error: err text}

    @response 500 {error: err text}

func (h *Handler) CreateOrderHandler(w http.ResponseWriter, r *http.Request)
    @description -> Create a new order

    @route -> /api/v1/orders/create/

    @method -> POST

    @body -> body should follow the following structure:

        {
        	OrderBasics: {
        		CustomerId: uint64
        		ApplicantId: uint64
        		OfferId: uint64
        		OrderId: uint64
        	},
        	OrderDetails: {
        		Title: string
        		Body: string
        		Price: float64
        		Status: string -> can be omitted here
        		CreatedAt: string -> can be omitted here
        		UpdatedAt: string -> can be omitted here
        		ExpireAt: string
        		Status: string -> can be omitted here
        		Price: uint64
        	},
        	PaymentDetails: {
        		TotalPrice: float64
        		RoundAmount: uint32
        		IsPaid: bool -> can be omitted here
        		PaidAt: string -> can be omitted here
        		UpdatedAt: string -> can be omitted here
        		PaymentRounds: []PaymentRound{

        {
        				Id: uint64 -> can be omitted here
        				Round: uint32
        				Amount: float64
        				PaidStatus: string
        				PaidAt: string -> can be omitted here
        				UpdatedAt: string -> can be omitted here
        },

        {
        				Id: uint64 -> can be omitted here
        				Round: uint32
        				Amount: float64
        				PaidStatus: string
        				PaidAt: string -> can be omitted here
        				UpdatedAt: string -> can be omitted here
        },

        		}
        	}
        }

    @response 201

    @response 400 {object} ErrorResponse {error: err text}

    @response 401 {object} ErrorResponse {error: err text}

    @response 403 {object} ErrorResponse {error: err text}

    @response 500 {object} ErrorResponse {error: err text}

func (h *Handler) CreateReviewHandler(w http.ResponseWriter, r *http.Request)
    @des cription -> Create a review

    @route -> /api/v1/reviews/create/

    @method -> POST

    @body -> body should follow the following structure:

        {
        	ReviewerId: uint64 -> can be omitted here
        	CustomerId: uint64
        	ReviewerId: uint64
        	OrderId: uint64
        	Rank: float32
        	Type: bool
        	Title: string
        	Body: string
        }

    @response 201

    @response 400 {object} ErrorResponse {error: err text}

    @response 401 {object} ErrorResponse {error: err text}

    @response 403 {object} ErrorResponse {error: err text}

    @response 500 {object} ErrorResponse {error: err text}

func (h *Handler) DeleteBlogHandler(w http.ResponseWriter, r *http.Request)
    @description -> Delete a post by postId

    @route -> /api/v1/blog/delete/{postId}/

    @method -> DELETE

    @body -> an empty one

    @response 204

    @response 400 {object} ErrorResponse {error: err text}

    @response 401 {object} ErrorResponse {error: err text}

    @response 403 {object} ErrorResponse {error: err text}

    @response 500 {object} ErrorResponse {error: err text}

func (h *Handler) DeleteNotificationHandler(w http.ResponseWriter, r *http.Request)
    @description -> Delete a notification by notificationId

    @route -> /api/v1/notifications/delete-notification/{notificationId}/

    @method -> DELETE

    @body -> an empty one

    @response 204

    @response 400 {object} ErrorResponse

    @response 401 {object} ErrorResponse

    @response 403 {object} ErrorResponse

    @response 500 {object} ErrorResponse

func (h *Handler) DeleteOfferHandler(w http.ResponseWriter, r *http.Request)
    @description -> Delete an offer by offerId

    @route -> /api/v1/offers/delete/{offerId}/

    @method -> DELETE

    @body -> an empty one

    @response 204

    @response 400 {object} ErrorResponse

    @response 401 {object} ErrorResponse

    @response 403 {object} ErrorResponse

    @response 500 {object} ErrorResponse

func (h *Handler) DeleteOrderHandler(w http.ResponseWriter, r *http.Request)
    @description -> Delete an order by orderId

    @route -> /api/v1/orders/delete/{orderId}/

    @method -> DELETE

    @body -> an empty one

    @response 204

    @response 400 {error: err text}

    @response 401 {error: err text}

    @response 403 {error: err text}

    @response 500 {error: err text}

func (h *Handler) DeleteReviewHandler(w http.ResponseWriter, r *http.Request)
    @description -> Delete a review by reviewId

    @route -> /api/v1/reviews/delete/{reviewId}/

    @method -> DELETE

    @body -> an empty one

    @response 204

    @response 400 {object} ErrorResponse {error: err text}

    @response 401 {object} ErrorResponse {error: err text}

    @response 403 {object} ErrorResponse {error: err text}

    @response 500 {object} ErrorResponse {error: err text}

func (h *Handler) FillProfileHandler(w http.ResponseWriter, r *http.Request)
    @description -> Fill a customer profile (bio, which use in public profile)

    @route -> /api/v1/profile/fill/

    @method -> POST

    @body -> body should follow the following structure:

        {

        		Base: {
        			CustomerId: uint64
        			Email: string -> use customer current email by default
             	Name: string
        			Avatar: string
        		}
        		Details: {
        			Title: string
        			Description: string
        			Tags: []string
        			Categories: []string
        		}

        }

    @response 201

    @response 400 {object} {error: err text}

    @response 401 {object} {error: err text}

    @response 403 {object} {error: err text}

    @response 500 {object} {error: err text}

func (h *Handler) GetApplicantsListHandler(w http.ResponseWriter, r *http.Request)
    @description -> Get a list of applicants for an offer by offer id

    @route -> /api/v1/offers/get-applicants-list/{offerId}/{skip}/

    @method -> GET

    @body -> an empty one

    @response 200:

        {
            Applicants: [20]offersPb.ApplicantDto{
            OfferId: uint64
            CustomerId: uint64
            AppliedAt: string
            Title: string
            Body: string
            Customer: profilePb.PublicProfileDto {
            Id: uint64
            Email: string
            Name: string
            Avatar: string
            Title: string
            Description: string
            Tags: []string
            }
            }
            Total: uint64
            }

    @response 400 {object} ErrorResponse

    @response 401 {object} ErrorResponse

    @response 403 {object} ErrorResponse

    @response 500 {object} ErrorResponse

func (h *Handler) GetBlogHandler(w http.ResponseWriter, r *http.Request)
    @description -> Get a list of posts in the customer's portfolio

    @route -> /api/v1/blog/get/{skip}/

    @method -> GET

    @body -> an empty one

    @response 200

          {
          Blog: [20]blogPb.BlogItem{
          {
          Id: uint64
          Title: string
          Description: string
          Likes: uint64
          Dislikes: uint64
          CreatedAt: string
          UpdatedAt: string
          Tags: []string
          Images: []string
          },
          {
          Id: uint64
          Title: string
          Description: string
          Likes: uint64
          Dislikes: uint64
          CreatedAt: string
          UpdatedAt: string
          Tags: []string
          Images: []string
          },
          }
          }

        @response 400 {object} ErrorResponse {error: err text}

        @response 401 {object} ErrorResponse {error: err text}

        @response 403 {object} ErrorResponse {error: err text}

        @response 500 {object} ErrorResponse {error: err text}

func (h *Handler) GetComplianceRequestsListHandler(w http.ResponseWriter, r *http.Request)
    @description -> Get a list of compliance requests by customerId

    @route -> /api/v1/orders/compliance/get-list/{skip}/

    @method -> GET

    @body -> an empty one

    @response 200

        		{

        		ComplianceRequests: {
        		{
        		Id: uint64
        		OrderBasics:
            {
            CustomerId: uint64
            ApplicantId: uint64
            OrderId: uint64
            OfferId: uint64
            }
            Round: uint32
            Status: string
            Claim: string
            CreatedAt: string
            },
            {
            Id: uint64
            OrderBasics:
            {
            CustomerId: uint64
            ApplicantId: uint64
            OrderId: uint64
            OfferId: uint64
            }
            Round: uint32
            Status: string
            Claim: string
            CreatedAt: string
            },
            }

    @response 400 {object} ErrorResponse {error: err text}

    @response 401 {object} ErrorResponse {error: err text}

    @response 403 {object} ErrorResponse {error: err text}

    @response 500 {object} ErrorResponse {error: err text}

func (h *Handler) GetCustomerKycHandler(w http.ResponseWriter, r *http.Request)
    @description -> Get a customer kyc by customerId

    @route -> /api/v1/profile/kyc/get/

    @method -> GET

    @body -> not implemented yet

    @response 200 {object} not implemented yet

    @response 400 {object} ErrorResponse {error: err text}

    @response 401 {object} ErrorResponse {error: err text}

    @response 403 {object} ErrorResponse {error: err text}

    @response 500 {object} ErrorResponse {error: err text}

func (h *Handler) GetCustomerProfileHandler(w http.ResponseWriter, r *http.Request)
    @description -> Get a customer profile by id for the customer, not a public
    one. If you don't got something - it cause an empty or false fields was
    omitted.

    @route -> /api/v1/profile/get/

    @method -> GET

    @body -> an empty one

    @response 200 {object}

        	{
        		Base:
        		{
        		CustomerId: uint64
            Email: string
            Name: string
            Avatar: string
            }
            Details:
            {
            Role: string
            CreatedAt: string
            TwoStepType: string
            }
            Params:
            {
            IsBanned: bool
            IsVerified: bool
            IsPremium: bool
            IsTwoFa: bool
            IsKyc: bool
            }
            Bio:
            {
            Title: string
            Description: string
            Tags: []string
            Categories: []string
            }
            Kyc: {}
            Orders:
            {
            CompletedProjects: uint32
            }
            Reviews:
            {
            TotalProjects: uint32
            Rank: float64
            }
            }

    @response 400 {object} ErrorResponse {error: err text}

    @response 401 {object} ErrorResponse {error: err text}

    @response 403 {object} ErrorResponse {error: err text}

    @response 500 {object} ErrorResponse {error: err text}

func (h *Handler) GetMyOffersHandler(w http.ResponseWriter, r *http.Request)
    @description -> handles requests to get a list of offers posted by the
    authenticated user.

    @route -> /api/v1/offers/get-my-offers/

    @method -> POST

    @body -> empty

    @response 200 {object}

        {
          "CardList": [
            {
              "OfferId": "uint64",
              "PostedBy": "uint64",
              "Title": "string",
              "Price": "float64",
              "CreatedAt": "string",
              "ExpireAt": "string",
              "TagsList": ["string"],
              "CategoriesList": ["string"]
            }
          ]
        }

    @response 400 {object} ErrorResponse

    @response 401 {object} ErrorResponse

    @response 403 {object} ErrorResponse

    @response 500 {object} ErrorResponse

func (h *Handler) GetNotificationsListHandler(w http.ResponseWriter, r *http.Request)
    @description -> Get notifications list by customer id

    @route -> /api/v1/notifications/get-notifications-list/{skip}/

    @method -> GET

    @body -> an empty one

    @response 200:

        	{
        	    List: {
        	    Id: uint64
        	    CustomerId: uint64
        	    Area: string
            	Title: string
            	Body: string
            	CreatedAt: string
            }
           }

    @response 400 {object} ErrorResponse

    @response 401 {object} ErrorResponse

    @response 403 {object} ErrorResponse

    @response 500 {object} ErrorResponse

func (h *Handler) GetOfferDetailsHandler(w http.ResponseWriter, r *http.Request)
    @description -> Get offer details by offer id

    @route -> /api/v1/offers/get-offer-details/{offerId}/

    @method -> GET

    @body -> an empty one

    @response 200 {object} GetOfferDetailsResponse

        {
        Offer: {
        OfferId: uint64
        PostedBy: uint64
        AppliedBy: []uint64
        ApprovedFor: uint64
        ExpireAt: string
        Details: {
        Title: string
        Description: string
        Status: string
        Price: uint64
        Location: string
        CreatedAt: string
        UpdatedAt: string
        ExpireAt: string
        TagsList: []string
        CategoriesList: []string
        }
        }

        }

    @response 400 {object} ErrorResponse {error: err text}

    @response 401 {object} ErrorResponse {error: err text}

    @response 403 {object} ErrorResponse {error: err text}

    @response 500 {object} ErrorResponse {error: err text}

func (h *Handler) GetOffersListHandler(w http.ResponseWriter, r *http.Request)
    @description -> Get a list of offers

    @route -> /api/v1/offers/get-offers-list/

    @method -> POST

    @body -> body should follow the following structure:

        {
        	Skip: uint64
        	CustomerId: uint64
        	MinPrice: uint64
        	MaxPrice: uint64
        	DateFrom: string
        	DateTo: string
        	Location: string
        	TagsList: []string
        	CategoriesList: []string
        }

    @response 200 {object} GetOffersListResponse

        {
        	CardList: [20]offersPb.OfferShortCard{
        		OfferId: uint64
        		PostedBy: uint64
        		Title: string
        		Price: uint64
        		CreatedAt: string
        		ExpireAt: string
        		TagsList: []string
        		CategoriesList: []string
        	}
        }

    @response 400 {object} {error: err text}

    @response 401 {object} {error: err text}

    @response 403 {object} {error: err text}

    @response 500 {object} {error: err text}

func (h *Handler) GetOrderDetailsHandler(w http.ResponseWriter, r *http.Request)
    @description -> Get order details by order id

    @route -> /api/v1/orders/get-order-details/{orderId}/

    @method -> GET

    @body -> an empty one

    @response 200

        {

        OrderBasics:
        {
        CustomerId: uint64
        ApplicantId: uint64
        OfferId: uint64
        OrderId: uint64
        }
        OrderDetails:
        {
        Id: uint64
        Rank: float64
        CompletedProjects: uint32
        },
        Applicant:
        {
        Id: uint64
        Rank: float64
        CompletedProjects: uint32
        },
        PaymentDetails:
        {
        TotalPrice: float64
        RoundAmount: uint32
        IsPaid: bool
        PaidAt: string
        PaymentRounds:
        {
        Id: uint64
        Round: uint32
        Amount: float64
        PaidStatus: string
        PaidAt: string
        UpdatedAt: string
        },
        {
        Id: uint64
        Round: uint32
        Amount: float64
        PaidStatus: string
        PaidAt: string
        UpdatedAt: string
        },
        },

        }

    @response 400 {object} ErrorResponse {error: err text}

    @response 401 {object} ErrorResponse {error: err text}

    @response 403 {object} ErrorResponse {error: err text}

    @response 500 {object} ErrorResponse {error: err text}

func (h *Handler) GetOrdersListByFilterHandler(w http.ResponseWriter, r *http.Request)
    @description -> Get a list of orders by filter

    @route -> /api/v1/orders/get-orders-list-by-filter/

    @method -> POST

    @body -> body should follow the following structure:

        {
        	CustomerId: uint64
        	ApplicantId: uint64
        	Skip: uint32
        	Status: string
        }

    @response 200 {object} GetOrdersListByFilterResponse

        {
        	OrderList: [20]ordersPb.OrderShortCard{
        		OrderId: uint64
        		Title: string
        		Body: string
        		Status: string
        		Price: uint64
        		RoundAmount: uint32
        		CreatedAt: string
        		Skip: uint32
        	}
        }

    @response 400 {object} ErrorResponse {error: err text}

    @response 401 {object} ErrorResponse {error: err text}

    @response 403 {object} ErrorResponse {error: err text}

    @response 500 {object} ErrorResponse {error: err text}

func (h *Handler) GetOrdersRequestsListByApplicantIdHandler(w http.ResponseWriter, r *http.Request)
    @description -> Get a list of orders requests by applicantId

    @route -> /api/v1/orders/get-orders-requests-list/{skip}/

    @method -> GET

    @body -> an empty one

    @response 200 {object} GetOrdersRequestsListByApplicantIdResponse

        {
        	OrderList: [20]ordersPb.OrderShortCard{
        		OrderId: uint64
        		Title: string
        		Body: string
        		Status: string
        		Price: uint64
        		RoundAmount: uint32
        		CreatedAt: string
        		Skip: uint32
        	}
        }

    @response 400 {object} ErrorResponse {error: err text}

    @response 401 {object} ErrorResponse {error: err text}

    @response 403 {object} ErrorResponse {error: err text}

    @response 500 {object} ErrorResponse {error: err text}

func (h *Handler) GetPublicProfileHandler(w http.ResponseWriter, r *http.Request)
    @description -> Get a public profile by id

    @route -> /api/v1/profile/public/get/

    @method -> GET

    @body -> an empty one

    @response 200 {object}

        	{
        		Id: uint64
        		Name: string
        		Email: string
        		Avatar: string
            Title: string
            Description: string
            Tags: []string
            Categories: []string
            IsVerified: bool
            IsPremium: bool
          }

    @response 400 {object} ErrorResponse {error: err text}

    @response 401 {object} ErrorResponse {error: err text}

    @response 403 {object} ErrorResponse {error: err text}

    @response 500 {object} ErrorResponse {error: err text}

func (h *Handler) GetReviewCommentsListHandler(w http.ResponseWriter, r *http.Request)
    @description -> Get a list of comments for a review

    @route -> /api/v1/reviews/get-review-comments-list/{reviewId}/{skip}/

    @method -> GET

    @response 200 {object} reviewsPb.GetReviewCommentsListResponse

        {
        		Comments: [20]*reviewsPb.ReviewCommentResponse

        			{
        				Comment: {
        					ReviewId: uint64
        					PostedBy: uint64
        					Title: string
        					Body: string
        					CreatedAt: string
        					UpdatedAt: string
        				}
        			}
        	}

    @response 400 {object} ErrorResponse {error: err text}

    @response 401 {object} ErrorResponse {error: err text}

    @response 403 {object} ErrorResponse {error: err text}

    @response 500 {object} ErrorResponse {error: err text}

func (h *Handler) GetReviewsListHandler(w http.ResponseWriter, r *http.Request)
    @response 400 {object} ErrorResponse {error: err text}

    @response 401 {object} ErrorResponse {error: err text}

    @response 403 {object} ErrorResponse {error: err text}

    @response 500 {object} ErrorResponse {error: err text}

func (h *Handler) HandleAuthForgotPwd(w http.ResponseWriter, r *http.Request)
    @description -> Forgot password

    @route -> /api/v1/auth/forgot-password/{email}

    @method -> POST

    @response 202

    @response 400 {object} ErrorResponse {error: err text}

    @response 401 {object} ErrorResponse {error: err text}

    @response 403 {object} ErrorResponse {error: err text}

    @response 500 {object} ErrorResponse {error: err text}

func (h *Handler) HandleAuthSignIn(w http.ResponseWriter, r *http.Request)
    @description -> Sign in a customer

    @route -> /api/v1/auth/sign-in

    @method -> POST

    @response 200 {object} authPb.SignInResponse

        {
        	Email: string
        	Password: string
        	TwoStepCode: string
        }

    @response 400 {object} ErrorResponse {error: err text}

    @response 401 {object} ErrorResponse {error: err text}

    @response 403 {object} ErrorResponse {error: err text}

    @response 500 {object} ErrorResponse {error: err text}

func (h *Handler) HandleAuthSignUp(w http.ResponseWriter, r *http.Request)
    @description -> Sign up a new customer

    @route -> /api/v1/auth/sign-up

    @method -> POST

    @response 200 {object} profilePb.CreateCustomerResponse

        {
        	Email: string
        	Name: string
        	Password: string
        }

    @response 400 {object} ErrorResponse {error: err text}

    @response 401 {object} ErrorResponse {error: err text}

    @response 403 {object} ErrorResponse {error: err text}

    @response 500 {object} ErrorResponse {error: err text}

func (h *Handler) RegisterAuthRoutes(mux *http.ServeMux)

func (h *Handler) RegisterBlogRoutes(mux *http.ServeMux)

func (h *Handler) RegisterChatRoutes(mux *http.ServeMux)

func (h *Handler) RegisterNotificationsRoutes(mux *http.ServeMux)

func (h *Handler) RegisterOffersRoutes(mux *http.ServeMux)

func (h *Handler) RegisterOrdersRoutes(mux *http.ServeMux)

func (h *Handler) RegisterPaymentsRoutes(mux *http.ServeMux)

func (h *Handler) RegisterProfileRoutes(mux *http.ServeMux)

func (h *Handler) RegisterPromotionsRoutes(mux *http.ServeMux)

func (h *Handler) RegisterReviewsRoutes(mux *http.ServeMux)

func (h *Handler) RejectAnOrderHandler(w http.ResponseWriter, r *http.Request)
    @description -> Reject an order

    @route -> /api/v1/orders/reject/

    @method -> POST

    @body -> body should follow the following structure:

        {
        	OrderBasics: {
        		CustomerId: uint64
        		ApplicantId: uint64
        		OrderId: uint64
        		OfferId: uint64
        	},
        	Cause: string
        }

    @response 200

    @response 400 {error: err text}

    @response 401 {error: err text}

    @response 403 {error: err text}

    @response 500 {error: err text}

func (h *Handler) RejectComplianceHandler(w http.ResponseWriter, r *http.Request)
    @description -> Reject a compliance request by offer owner

    @route -> /api/v1/orders/compliance/reject/

    @method -> POST

    @body -> body should follow the following structure:

        {
        	Id: uint64
        	Round: uint32
        	Status: string
        	Claim: string
        	CreatedAt: string
        	OrderBasics: {
        		CustomerId: uint64
        		ApplicantId: uint64
        		OrderId: uint64
        		OfferId: uint64
        	}
        }

    @response 200

    @response 400 {error: err text}

    @response 401 {error: err text}

    @response 403 {error: err text}

    @response 500 {error: err text}

func (h *Handler) ReportCustomerHandler(w http.ResponseWriter, r *http.Request)
    @description -> Report a customer

    @route -> /api/v1/profile/report/

    @method -> POST

    @body -> body should follow the following structure:

        {
        	CustomerId: uint64
        	ReporterId: uint64
        	Reason: string
        	Description: string
        }

    Response list:

    @response 200 {object}

    @response 400 {object} ErrorResponse {error: err text}

    @response 401 {object} ErrorResponse {error: err text}

    @response 403 {object} ErrorResponse {error: err text}

    @response 500 {object} ErrorResponse {error: err text}

func (h *Handler) SetReactionHandler(w http.ResponseWriter, r *http.Request)
    @description -> Set a reaction to a post

    @route -> /api/v1/blog/reaction/

    @method -> POST

    @body -> body should follow the following structure:

        {
        	PostId: uint64
        	Like: bool
        	Dislike: bool
        }

    @response 202

    @response 400 {object} ErrorResponse {error: err text}

    @response 401 {object} ErrorResponse {error: err text}

    @response 403 {object} ErrorResponse {error: err text}

    @response 500 {object} ErrorResponse {error: err text}

func (h *Handler) SetReviewReactionHandler(w http.ResponseWriter, r *http.Request)
    @description -> Set a reaction to a review

    @route -> /api/v1/reviews/reaction/

    @method -> POST

    @body -> body should follow the following structure:

        {
        	ReviewId: uint64
        	CustomerId: uint64
        	Helpful: bool
        }

    @response 204

    @response 400 {object} ErrorResponse {error: err text}

    @response 401 {object} ErrorResponse {error: err text}

    @response 403 {object} ErrorResponse {error: err text}

    @response 500 {object} ErrorResponse {error: err text}

func (h *Handler) UpdateBlogHandler(w http.ResponseWriter, r *http.Request)
    @description -> Update a post

    @route -> /api/v1/blog/update/

    @method -> POST

    @body -> body should follow the following structure:

        {
        	Id: uint64
        	CustomerId: uint64
        	Title: string
        	Description: string
        	TagList: []string
        	ImageList: []string
        }

    @response 204

    @response 400 {object} ErrorResponse {error: err text}

    @response 401 {object} ErrorResponse {error: err text}

    @response 403 {object} ErrorResponse {error: err text}

    @response 500 {object} ErrorResponse {error: err text}

func (h *Handler) UpdateCustomeKycHandler(w http.ResponseWriter, r *http.Request)
    @description -> Update a customer kyc data

    @route -> /api/v1/profile/kyc/update/

    @method -> PUT

    @body -> not implemented yet

    @response 200 {object}

    @response 400 {object} ErrorResponse {error: err text}

    @response 401 {object} ErrorResponse {error: err text}

    @response 403 {object} ErrorResponse {error: err text}

    @response 500 {object} ErrorResponse {error: err text}

func (h *Handler) UpdateCustomerProfileHandler(w http.ResponseWriter, r *http.Request)
    @description -> Update a customer bio only

    @route -> /api/v1/profile/update/

    @method -> POST

    @body -> body should follow the following structure:

        {

        		Base: {
        			CustomerId: uint64
        			Email: string -> can be omitted here
             	Name: string
        			Avatar: string
        		}
        		Details: {
        			Title: string
        			Description: string
        			Tags: []string
        			Categories: []string
        		}

        }

    @response 200

    @response 400 {object} {error: err text}

    @response 401 {object} {error: err text}

    @response 403 {object} {error: err text}

    @response 500 {object} {error: err text}

func (h *Handler) UpdateOfferHandler(w http.ResponseWriter, r *http.Request)
    @description -> Update an offer

    @route -> /api/v1/offers/update/

    @method -> POST

    @body -> body should follow the following structure:

        {
        	OfferId: uint64
        	Title: string
        	Description: string
        	Location: string
        	Price: uint64
        	Status: string
        	ExpireAt: string
        	PostedBy: uint64
        	TagsList: []string
        	CategoriesList: []string
        }

    @response 200

    @response 400 {error: err text}

    @response 401 {error: err text}

    @response 403 {error: err text}

    @response 500 {error: err text}

func (h *Handler) UpdateOrderHandler(w http.ResponseWriter, r *http.Request)
    @description -> Update an order

    @route -> /api/v1/orders/update/

    @method -> PUT

    @body -> body should follow the following structure:

        {
        	OrderBasics: {
        		CustomerId: uint64
        		ApplicantId: uint64
        		OfferId: uint64
        		OrderId: uint64
        	},
        	OrderDetails: {
        		Title: string
        		Body: string
        		Price: float64
        		Status: string
        		CreatedAt: string -> can be omitted here
        		UpdatedAt: string -> can be omitted here
        		ExpireAt: string
        		Status: string -> can be omitted here
        		Price: uint64
        	},
        	PaymentDetails: {
        		TotalPrice: float64
        		RoundAmount: uint32
        		IsPaid: bool -> can be omitted here
        		PaidAt: string -> can be omitted here
        		UpdatedAt: string -> can be omitted here
        		PaymentRounds: []PaymentRound{

        {
        				Id: uint64 -> can be omitted here
        				Round: uint32
        				Amount: float64
        				PaidStatus: string
        				PaidAt: string -> can be omitted here
        				UpdatedAt: string -> can be omitted here
        },

        {
        				Id: uint64 -> can be omitted here
        				Round: uint32
        				Amount: float64
        				PaidStatus: string
        				PaidAt: string -> can be omitted here
        				UpdatedAt: string -> can be omitted here
        },

        		}
        	}
        }

    @response 204

    @response 400 {object} ErrorResponse {error: err text}

    @response 401 {object} ErrorResponse {error: err text}

    @response 403 {object} ErrorResponse {error: err text}

    @response 500 {object} ErrorResponse {error: err text}

func (h *Handler) UpdateReviewHandler(w http.ResponseWriter, r *http.Request)
    @description -> Update a review only if you are the reviewer

    @route -> /api/v1/reviews/update/

    @method -> PUT

    @body -> body should follow the following structure:

        {
        	ReviewerId: uint64
        	CustomerId: uint64
        	ReviewerId: uint64
        	OrderId: uint64
        	Rank: float32
        	Type: bool
        	Title: string
        	Body: string
        }

    @response 204

    @response 400 {object} ErrorResponse {error: err text}

    @response 401 {object} ErrorResponse {error: err text}

    @response 403 {object} ErrorResponse {error: err text}

    @response 500 {object} ErrorResponse {error: err text}

