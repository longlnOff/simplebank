package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/longln/simplebank/db/sqlc"
)


// type server for handling database and routing
type Server struct {
	store *db.Store
	router *gin.Engine
}

// create new server for handling
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}		// To interact with DB
	router := gin.Default()				// To route requests

	// register routes

	// we need bind to Server struct because we need to interact with DB to create account and of course, routing
	router.POST("/accounts", server.createAccount)	// create account
	router.GET("/accounts/:id", server.getAccount)	// get account
	router.GET("/accounts", server.listAccount)	// get account

	server.router = router
	return server
}


func (server *Server) StartServer(address string) error {
	return server.router.Run(address)

}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}