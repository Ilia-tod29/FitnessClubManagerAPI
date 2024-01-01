package api

import (
	db "FitnessClubManagerAPI/db/sqlc"
	"github.com/gin-gonic/gin"
	"net/http"
)

type createUserRequest struct {
	Email          string `json:"email" binding:"required"`
	HashedPassword string `json:"hashed_password" binding:"required"`
	Role           string `json:"role" binding:"required"`
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
		Role:           req.Role,
	}

	user, err := s.store.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}
