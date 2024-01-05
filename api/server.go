package api

import (
	db "FitnessClubManagerAPI/db/sqlc"
	"FitnessClubManagerAPI/token"
	"FitnessClubManagerAPI/util"
	"fmt"
	"github.com/gin-gonic/gin"
)

// idRequest represents the data needed from the user when only ID is needed
type idRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// listResourceByPagesRequest represents the data needed from the user when listing resources by pages
type listResourceByPagesRequest struct {
	PageID   int64 `form:"page_id" binding:"required,min=1"`
	PageSize int64 `form:"page_size" binding:"required,min=5,max=20"`
}

// Server serves HTTP requests for our banking service.
type Server struct {
	config     util.Config
	store      *db.SQLStore
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(config util.Config, store *db.SQLStore) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	//if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
	//	v.RegisterValidation("currency", validCurrency)
	//}

	server.setupRouter()
	return server, nil
}

func (s *Server) setupRouter() {
	router := gin.Default()

	authRoutes := router.Group("/").Use(authMiddleware(s.tokenMaker))

	// Users
	// TODO: additional functionality - change password
	// Unprotected
	router.POST("/users", s.createUser)
	router.POST("/users/login", s.loginUser)
	router.POST("/tokens/renew_access", s.renewAccessToken)
	// Protected
	authRoutes.PUT("/users/:id", s.updateUser)
	authRoutes.GET("/users/:id", s.getUser)
	authRoutes.GET("/allusers", s.listAllUsers)
	authRoutes.GET("/users", s.listUsersByPages)
	authRoutes.DELETE("/users/:id", s.deleteUser)

	// Subscriptions
	// Protected
	authRoutes.POST("/subscriptions", s.createSubscription)
	authRoutes.GET("/subscriptions/:id", s.getSubscription)
	authRoutes.GET("/subscriptions/user/:id", s.getAllSubscriptionsForAGivenUser)
	authRoutes.GET("/allsubscriptions", s.listAllSubscriptions)
	authRoutes.GET("/subscriptions", s.listSubscriptionsByPages)
	authRoutes.DELETE("/subscriptions/:id", s.deleteSubscription)

	// Inventory Items
	// Unprotected
	router.GET("/inventoryitems/:id", s.getInventoryItem)
	router.GET("/allinventoryitems", s.listAllInventoryItems)
	router.GET("/inventoryitems", s.listInventoryItemsByPages)
	// Protected
	authRoutes.POST("/inventoryitems", s.createInventoryItem)
	authRoutes.PUT("/inventoryitems/:id", s.updateInventoryItem)
	authRoutes.DELETE("/inventoryitems/:id", s.deleteInventoryItem)

	// Gallery Items
	// Unprotected
	router.GET("/gallery/:id", s.getGalleryItem)
	router.GET("/allgallery", s.listAllGalleryItems)
	router.GET("/gallery", s.listGalleryItemsByPages)
	// Protected
	authRoutes.POST("/gallery", s.createGalleryItem)
	authRoutes.DELETE("/gallery/:id", s.deleteGalleryItem)

	s.router = router
}

// Start runs the HTTP server on a specific address.
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
