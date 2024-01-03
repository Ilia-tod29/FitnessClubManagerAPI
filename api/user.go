package api

import (
	db "FitnessClubManagerAPI/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"net/http"
)

type createUserRequest struct {
	Email          string `json:"email" binding:"required"`
	HashedPassword string `json:"hashed_password" binding:"required"`
	Role           string `json:"role" binding:"required,oneof=admin user"`
}

type listUserRequest struct {
	PageID   int64 `form:"page_id" binding:"required,min=1"`
	PageSize int64 `form:"page_size" binding:"required,min=5,max=20"`
}

type getUserRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type updateUserRequest struct {
	Suspended string `json:"suspended" binding:"required,oneof=true false"`
}

func (s *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Email:          req.Email,
		HashedPassword: req.HashedPassword,
		// We don't suspend a user on creation
		Suspended: false,
		Role:      req.Role,
	}

	user, err := s.store.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (s *Server) getUser(ctx *gin.Context) {
	var req getUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := s.store.GetUser(ctx, req.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (s *Server) listAllUsers(ctx *gin.Context) {
	allUsers, err := s.store.ListAllUsers(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, allUsers)
}

func (s *Server) listUsers(ctx *gin.Context) {
	var req listUserRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.ListUsersParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	users, err := s.store.ListUsers(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (s *Server) updateUser(ctx *gin.Context) {
	var req getUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var upd updateUserRequest
	if err := ctx.ShouldBindJSON(&upd); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	isSuspended := false
	if upd.Suspended == "true" {
		isSuspended = true
	}
	args := db.UpdateUserParams{
		ID:        req.ID,
		Suspended: isSuspended,
	}

	user, err := s.store.UpdateUser(ctx, args)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (s *Server) deleteUser(ctx *gin.Context) {
	var req getUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := s.store.DeleteUser(ctx, req.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}
