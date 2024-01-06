package api

import (
	db "FitnessClubManagerAPI/db/sqlc"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentintent"
	"github.com/stripe/stripe-go/v76/webhook"
	"io"
	"net/http"
	"strconv"
)

func (s Server) handleWebhook(ctx *gin.Context) {
	const MaxBodyBytes = int64(65536)
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, MaxBodyBytes)
	payload, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		errToReturn := fmt.Errorf("error reading request body: %v", err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(errToReturn))
		return
	}

	event := stripe.Event{}

	if err := json.Unmarshal(payload, &event); err != nil {
		errToReturn := fmt.Errorf("webhook error while parsing basic request: %v", err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(errToReturn))
		return
	}

	// Replace this endpoint secret with your endpoint's unique secret
	// If you are testing with the CLI, find the secret by running 'stripe listen'
	// If you are using an endpoint defined with the API or dashboard, look in your webhook settings
	// at https://dashboard.stripe.com/webhooks
	endpointSecret := "whsec_709f12babc10f61ace3918a07fc65fce1490bcd2d1beed43efc0809e845924bb"
	signatureHeader := ctx.Request.Header.Get("Stripe-Signature")
	event, err = webhook.ConstructEvent(payload, signatureHeader, endpointSecret)
	if err != nil {
		errToReturn := fmt.Errorf("webhook signature verification failed: %v", err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(errToReturn))
		return
	}

	// Handle the event based on its type
	switch event.Type {
	case "checkout.session.completed":
		var checkoutSessionCompleted stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &checkoutSessionCompleted)
		if err != nil {
			return
		}
		params := &stripe.PaymentIntentParams{}
		result, err := paymentintent.Get(checkoutSessionCompleted.PaymentIntent.ID, params)
		if result.Status == stripe.PaymentIntentStatusSucceeded {
			userId, err := strconv.ParseInt(checkoutSessionCompleted.Metadata["user_id"], 10, 64)
			if err != nil {
				errToReturn := fmt.Errorf("error converting string to int64: %v", err.Error())
				ctx.JSON(http.StatusBadRequest, errorResponse(errToReturn))
				return
			}
			s.registerSubscription(ctx, checkoutSessionCompleted.Metadata["start_date"], checkoutSessionCompleted.Metadata["end_date"], userId)
		}
	default:
		fmt.Printf("Unhandled event type: %s\n", event.Type)
	}
}

func (s Server) registerSubscription(ctx *gin.Context, startDate, endDate string, userId int64) {
	parsedStartDate, parsedEndDate, err := parseDates(ctx, startDate, endDate)
	if err != nil {
		return
	}

	pgStartDate, pgEndDate, err := s.validatePeriod(ctx, parsedStartDate, parsedEndDate, userId)
	if err != nil {
		return
	}

	arg := db.CreateSubscriptionParams{
		UserID:    userId,
		StartDate: pgStartDate,
		EndDate:   pgEndDate,
	}
	subscription, err := s.store.CreateSubscription(ctx, arg)
	if err != nil {
		errCode := db.ErrorCode(err)
		if errCode == db.ForeignKeyViolation {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, subscription)
}
