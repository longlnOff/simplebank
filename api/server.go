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
	server := &Server{store: store}
	router := gin.Default()

	// register routes

	// we need bind to Server struct because we need to interact with DB to create account and of course, routing
	router.POST("/account", server.createAccount)	// create account

	server.router = router
	return server
}