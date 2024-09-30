package main

import (
	"context"
	"database/sql"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"
	_ "github.com/lib/pq"
	"github.com/longln/simplebank/api"
	db "github.com/longln/simplebank/db/sqlc"
	"github.com/longln/simplebank/gapi"
	"github.com/longln/simplebank/pb"
	"github.com/longln/simplebank/utils"
	"github.com/longln/simplebank/worker"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {

    config, err := utils.LoadConfig(".")
	if config.Enviroment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

    if err != nil {
        log.Fatal().Msg("cannot load config")
    }


    conn, err := sql.Open(config.DBDriver, config.DBSource)
    if err != nil {
        log.Fatal().Msg("cannot connect to database")
    }
    defer conn.Close()


    log.Info().Msg("Successfully connected to Database!")

	store := db.NewStore(conn)

	// redis connection
	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)


	// // run HTTP Gin server
	// RunGinServer(config, store)

	// Use 2 goroutines to run both server
	// run GRPC server
	go runTaskProcessor(redisOpt, store)
	go RunGRPCServer(config, store, taskDistributor)
	RunGatewayServer(config, store, taskDistributor)

}


func runTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) {
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store)
	log.Info().Msg("start task processor")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Msg("failed to start task processor")
	}
}

func RunGatewayServer(config utils.Config, store db.Store, taskDistributor worker.TaskDistributor) {

	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Msg("cannot create server:")
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
		log.Fatal().Msg("cannot register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	// start server
	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Msg("cannot create listener:")
	}
	log.Printf("start HTTP Gateway server on %s", listener.Addr().String())
	handler := gapi.HttpLogger(mux)
	err = http.Serve(listener, handler)
	if err != nil {
		log.Fatal().Msg("cannot start HTTTP Gateway server:")
	}

}

func RunGRPCServer(config utils.Config, store db.Store, taskDistributor worker.TaskDistributor) {

	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Msg("cannot create server")
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	// start server
	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Msg("cannot create listener")
	}
	log.Printf("start gRPC server on %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Msg("cannot start gRPC server")
	}

}

func RunGinServer(config utils.Config, store db.Store) {
	server, err := api.NewServer(config, store)
    if err != nil {
		log.Fatal().Msgf("%s", err)
	}
	err = server.StartServer(config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Msgf("%s", err)
	}
}