package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/longln/simplebank/api"
	db "github.com/longln/simplebank/db/sqlc"
	"github.com/longln/simplebank/gapi"
	"github.com/longln/simplebank/pb"
	"github.com/longln/simplebank/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
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

	// Use 2 goroutines to run both server
	// run GRPC server
	go RunGRPCServer(config, store)
	RunGatewayServer(config, store)

}


func RunGatewayServer(config utils.Config, store db.Store) {

	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("cannot register handler server:", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	// start server
	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot create listener:", err)
	}
	log.Printf("start HTTP Gateway server on %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("cannot start HTTTP Gateway server:", err)
	}

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