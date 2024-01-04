package api

import (
	db "FitnessClubManagerAPI/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"net/http"
)

type createGalleryItemRequest struct {
	Image string `json:"image" binding:"required"`
}

// TODO: Handle Auth by role
func (s *Server) createGalleryItem(ctx *gin.Context) {
	var req createGalleryItemRequest
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

	galleryItem, err := s.store.CreateGalleryItem(ctx, image)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, galleryItem)
}

func (s *Server) getGalleryItem(ctx *gin.Context) {
	var req idRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	galleryItem, err := s.store.GetGalleryItem(ctx, req.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, galleryItem)
}

func (s *Server) listAllGalleryItems(ctx *gin.Context) {
	allGalleryItems, err := s.store.ListAllGalleryItems(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, allGalleryItems)
}

func (s *Server) listGalleryItemsByPages(ctx *gin.Context) {
	var req listResourceByPagesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.ListGalleryItemsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	galleryItems, err := s.store.ListGalleryItems(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, galleryItems)
}

// TODO: Handle Auth by role
func (s *Server) deleteGalleryItem(ctx *gin.Context) {
	var req idRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	galleryItem, err := s.store.DeleteGalleryItem(ctx, req.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, galleryItem)
}
