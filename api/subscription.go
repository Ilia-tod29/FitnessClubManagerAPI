package api

import (
	db "FitnessClubManagerAPI/db/sqlc"
	"FitnessClubManagerAPI/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"net/http"
	"time"
)

type createSubscriptionRequest struct {
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
}

func (s *Server) createSubscription(ctx *gin.Context) {
	var req createSubscriptionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	startDate, err := time.Parse("02.01.2006", req.StartDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	endDate, err := time.Parse("02.01.2006", req.EndDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if startDate.After(endDate) {
		err := fmt.Errorf("the start date cannot be after the end date")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var currentUser db.User
	err = s.getCurrentUser(ctx, &currentUser)
	if err != nil {
		return
	}

	allSubscriptions, err := s.store.ListAllSubscriptionsForAGivenUser(ctx, currentUser.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	for _, subscription := range allSubscriptions {
		if isDateWithinRange(startDate, subscription.StartDate.Time, subscription.EndDate.Time) {
			err := fmt.Errorf("the start date cannot be within the validity period of another subscription")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}

	var pgDateStartDate pgtype.Date
	var pgDateEndDate pgtype.Date
	err = pgDateStartDate.Scan(startDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	err = pgDateEndDate.Scan(endDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateSubscriptionParams{
		UserID:    currentUser.ID,
		StartDate: pgDateStartDate,
		EndDate:   pgDateEndDate,
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

func (s *Server) getSubscription(ctx *gin.Context) {
	var req idRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	subscription, err := s.store.GetSubscription(ctx, req.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var currentUser db.User
	err = s.getCurrentUser(ctx, &currentUser)
	if err != nil {
		return
	}

	if currentUser.Role != util.AdminRole && currentUser.ID != subscription.UserID {
		err := fmt.Errorf("subscription doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, subscription)
}

func (s *Server) getAllSubscriptionsForAGivenUser(ctx *gin.Context) {
	var req idRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	subscriptions, err := s.store.ListAllSubscriptionsForAGivenUser(ctx, req.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var currentUser db.User
	err = s.getCurrentUser(ctx, &currentUser)
	if err != nil {
		return
	}

	if currentUser.Role != util.AdminRole && currentUser.ID != req.ID {
		err := fmt.Errorf("subscriptions doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, subscriptions)
}

func (s *Server) listAllSubscriptions(ctx *gin.Context) {
	err := s.validateAdminPermissions(ctx)
	if err != nil {
		return
	}

	allSubscription, err := s.store.ListAllSubscriptions(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, allSubscription)
}

func (s *Server) listSubscriptionsByPages(ctx *gin.Context) {
	var req listResourceByPagesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := s.validateAdminPermissions(ctx)
	if err != nil {
		return
	}

	args := db.ListSubscriptionsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	users, err := s.store.ListSubscriptions(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (s *Server) deleteSubscription(ctx *gin.Context) {
	var req idRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := s.validateAdminPermissions(ctx)
	if err != nil {
		return
	}

	subscription, err := s.store.DeleteSubscription(ctx, req.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, subscription)
}

func (s Server) deleteOutdatedSubscriptions(ctx *gin.Context, userID int64, allSubscriptions []db.Subscription) ([]db.Subscription, error) {
	var deletedSubscriptions []db.Subscription
	for _, subscription := range allSubscriptions {
		if subscription.EndDate.Time.Before(time.Now()) {
			subscription, err := s.store.DeleteSubscription(ctx, userID)
			if err != nil {
				if err == pgx.ErrNoRows {
					ctx.JSON(http.StatusNotFound, errorResponse(err))
					return nil, err
				}
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return nil, err
			}

			deletedSubscriptions = append(deletedSubscriptions, subscription)
		}
	}
	return deletedSubscriptions, nil
}

func isDateWithinRange(targetDate, startDate, endDate time.Time) bool {
	return !targetDate.Before(startDate) && !targetDate.After(endDate)
}
