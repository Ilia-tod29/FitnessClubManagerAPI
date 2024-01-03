package api

import (
	db "FitnessClubManagerAPI/db/sqlc"
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
	store  *db.SQLStore
	router *gin.Engine
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(store *db.SQLStore) (*Server, error) {
	//tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	//if err != nil {
	//	return nil, fmt.Errorf("cannot create token maker: %w", err)
	//}

	server := &Server{
		//config:     config,
		store: store,
		//tokenMaker: tokenMaker,
	}

	router := gin.Default()

	// Users
	// TODO: additional functionality - change password
	router.POST("/users", server.createUser)
	router.PUT("/users/:id", server.updateUser)
	router.GET("/users/:id", server.getUser)
	router.GET("/allusers", server.listAllUsers)
	router.GET("/users", server.listUsersByPages)
	router.DELETE("/users/:id", server.deleteUser)

	// Subscriptions
	router.POST("/subscriptions", server.createSubscription)
	router.GET("/subscriptions/:id", server.getSubscription)
	router.GET("/subscriptions/user/:id", server.getAllSubscriptionsForAGivenUser)
	router.GET("/allsubscriptions", server.listAllSubscriptions)
	router.GET("/subscriptions", server.listSubscriptionsByPages)
	router.DELETE("/subscriptions/:id", server.deleteSubscription)

	// Inventory Items
	router.POST("/inventoryitems", server.createInventoryItem)
	router.PUT("/inventoryitems/:id", server.updateInventoryItem)
	router.GET("/inventoryitems/:id", server.getInventoryItem)
	router.GET("/allinventoryitems", server.listAllInventoryItems)
	router.GET("/inventoryitems", server.listInventoryItemsByPages)
	router.DELETE("/inventoryitems/:id", server.deleteInventoryItem)

	// Gallery Items
	router.POST("/gallery", server.createGalleryItem)
	router.GET("/gallery/:id", server.getGalleryItem)
	router.GET("/allgallery", server.listAllGalleryItems)
	router.GET("/gallery", server.listGalleryItemsByPages)
	router.DELETE("/gallery/:id", server.deleteGalleryItem)

	server.router = router

	//if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
	//	v.RegisterValidation("currency", validCurrency)
	//}
	//
	//server.setupRouter()
	return server, nil
}

// Start runs the HTTP server on a specific address.
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
