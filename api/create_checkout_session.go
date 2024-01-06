package api

import (
	db "FitnessClubManagerAPI/db/sqlc"
	"FitnessClubManagerAPI/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
	"net/http"
	"strconv"
)

const (
	homePageNavigator          = "/home.html"
	subscriptionsPageNavigator = "/subscriptions.html"
)

type createCheckoutSessionRequest struct {
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
}

func (s *Server) createCheckoutSession(ctx *gin.Context) {
	var currentUser db.User
	err := s.getCurrentUser(ctx, &currentUser)
	if err != nil {
		return
	}
	if currentUser.Role == util.AdminRole {
		err := fmt.Errorf("user don't have permissions to access this resource")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	var req createCheckoutSessionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	params := stripe.CheckoutSessionParams{
		PaymentMethodTypes: getPaymentMethods(),
		Mode:               stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				// Provide the exact Price ID (for example, pr_1234) of the product you want to sell
				Price:    stripe.String("price_1OVMJPETgyr4pC7GeEwfnQHE"),
				Quantity: stripe.Int64(1),
			},
		},
		Metadata: map[string]string{
			"start_date": req.StartDate,
			"end_date":   req.EndDate,
			"user_id":    strconv.FormatInt(currentUser.ID, 10),
		},
		SuccessURL: stripe.String(s.config.ClientBaseUrl + homePageNavigator),
		CancelURL:  stripe.String(s.config.ClientBaseUrl + subscriptionsPageNavigator),
	}

	currentSession, err := session.New(&params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, currentSession.URL)
}

func getPaymentMethods() []*string {
	card := stripe.String("card")
	return []*string{card}
}
