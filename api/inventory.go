package api

import (
	db "FitnessClubManagerAPI/db/sqlc"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"net/http"
)

type createInventoryItemRequest struct {
	Name  string `json:"name" binding:"required"`
	Image string `json:"image"`
}

type updateInventoryItemRequest struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

// TODO: Handle Auth by role
func (s *Server) createInventoryItem(ctx *gin.Context) {
	var req createInventoryItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var image pgtype.Text
	err := image.Scan(req.Image)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.CreateInventoryItemParams{
		Name:  req.Name,
		Image: image,
	}

	inventoryItem, err := s.store.CreateInventoryItem(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, inventoryItem)
}

func (s *Server) getInventoryItem(ctx *gin.Context) {
	var req idRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	inventoryItem, err := s.store.GetInventoryItem(ctx, req.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, inventoryItem)
}

func (s *Server) listAllInventoryItems(ctx *gin.Context) {
	allInventoryItems, err := s.store.ListAllInventoryItems(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, allInventoryItems)
}

func (s *Server) listInventoryItemsByPages(ctx *gin.Context) {
	var req listResourceByPagesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.ListInventoryItemsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	inventoryItems, err := s.store.ListInventoryItems(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, inventoryItems)
}

// TODO: Handle Auth by role
func (s *Server) updateInventoryItem(ctx *gin.Context) {
	var req idRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var upd updateInventoryItemRequest
	if err := ctx.ShouldBindJSON(&upd); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if upd.Image == "" && upd.Name == "" {
		err := fmt.Errorf("at least one of the two parameters 'name' or 'image' should be provided")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var image pgtype.Text

	if upd.Image != "" {
		err := image.Scan(upd.Image)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}

	if upd.Name == "" {
		args := db.UpdateInventoryItemImageParams{
			ID:    req.ID,
			Image: image,
		}

		inventoryItem, err := s.store.UpdateInventoryItemImage(ctx, args)
		if err != nil {
			if err == pgx.ErrNoRows {
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusOK, inventoryItem)
		return
	}

	if upd.Image == "" {
		args := db.UpdateInventoryItemNameParams{
			ID:   req.ID,
			Name: upd.Name,
		}

		inventoryItem, err := s.store.UpdateInventoryItemName(ctx, args)
		if err != nil {
			if err == pgx.ErrNoRows {
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusOK, inventoryItem)
		return
	}

	args := db.UpdateInventoryItemParams{
		ID:    req.ID,
		Name:  upd.Name,
		Image: image,
	}

	inventoryItem, err := s.store.UpdateInventoryItem(ctx, args)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, inventoryItem)
}

// TODO: Handle Auth by role
func (s *Server) deleteInventoryItem(ctx *gin.Context) {
	var req idRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	inventoryItem, err := s.store.DeleteInventoryItem(ctx, req.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, inventoryItem)
}
