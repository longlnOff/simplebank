package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	_ "github.com/lib/pq"
	"github.com/longln/simplebank/api"
	db "github.com/longln/simplebank/db/sqlc"
	"github.com/longln/simplebank/gapi"
	"github.com/longln/simplebank/pb"
	"github.com/longln/simplebank/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)


func main() {
    config, err := utils.LoadConfig(".")
    if err != nil {
        log.Fatal("cannot load config:", err)
    }


    conn, err := sql.Open(config.DBDriver, config.DBSource)
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    err = conn.Ping()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Successfully connected to Database!")

	store := db.NewStore(conn)


	// // run HTTP Gin server
	// RunGinServer(config, store)

	// run GRPC server
	RunGRPCServer(config, store)

}

func RunGRPCServer(config utils.Config, store db.Store) {

	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	// start server
	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot create listener:", err)
	}
	log.Printf("start gRPC server on %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start gRPC server:", err)
	}

}

func RunGinServer(config utils.Config, store db.Store) {
	server, err := api.NewServer(config, store)
    if err != nil {
		log.Fatal(err)
	}
	err = server.StartServer(config.HTTPServerAddress)
	if err != nil {
		log.Fatal(err)
	}
}