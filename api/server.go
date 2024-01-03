package api

import (
	db "FitnessClubManagerAPI/db/sqlc"
	"github.com/gin-gonic/gin"
)

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

	// TODO: additional functionality - change password
	router.POST("/users", server.createUser)
	router.PUT("/users/:id", server.updateUser)
	router.GET("/users/:id", server.getUser)
	router.GET("/allusers", server.listAllUsers)
	router.GET("/users", server.listUsers)
	router.DELETE("/users/:id", server.deleteUser)

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
