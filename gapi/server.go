package gapi

import (
	"fmt"

	db "github.com/longln/simplebank/db/sqlc"
	"github.com/longln/simplebank/pb"
	"github.com/longln/simplebank/token"
	"github.com/longln/simplebank/utils"
)

// type server for handling database and routing
type Server struct {
	pb.UnimplementedSimpleBankServer
	config utils.Config
	store db.Store
	tokenMaker token.Maker
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

	return server, nil
}