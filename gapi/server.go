package gapi

import (
	"fmt"

	db "github.com/longln/simplebank/db/sqlc"
	"github.com/longln/simplebank/pb"
	"github.com/longln/simplebank/token"
	"github.com/longln/simplebank/utils"
	"github.com/longln/simplebank/worker"
)

// type server for handling database and routing
type Server struct {
	config utils.Config
	store db.Store
	tokenMaker token.Maker
	taskDistributor worker.TaskDistributor
	pb.UnimplementedSimpleBankServer
}

// create new server for handling
func NewServer(config utils.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
        return nil, fmt.Errorf("cannot create token maker: %w", err)
    }
	server := &Server{
		tokenMaker: tokenMaker, 
		store: store,
		config: config,
		taskDistributor: taskDistributor,
	}

	return server, nil
}