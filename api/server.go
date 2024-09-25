package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/longln/simplebank/db/sqlc"
	"github.com/longln/simplebank/token"
	"github.com/longln/simplebank/utils"
)

// type server for handling database and routing
type Server struct {
	config utils.Config
	store db.Store
	tokenMaker token.Maker
	router *gin.Engine
}

// create new server for handling
func NewServer(config utils.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
        return nil, fmt.Errorf("cannot create token maker: %w", err)
    }
	server := &Server{
		tokenMaker: tokenMaker, 
		store: store,
		config: config,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {

	router := gin.Default()				// To route requests
	// register routes
	// we need bind to Server struct because we need to interact with DB to create account and of course, routing
	router.POST("/users", server.createUser)	// create user
	router.POST("/users/login", server.loginUser) // login
	
	// protected by middleware
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRoutes.POST("/accounts", server.createAccount)	// create account
	authRoutes.GET("/accounts/:id", server.getAccount)	// get account
	authRoutes.GET("/accounts", server.listAccount)	// list accounts
	authRoutes.POST("/transfers", server.createTransfer)	// create transfer

	
	server.router = router
}


func (server *Server) StartServer(address string) error {
	return server.router.Run(address)

}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}