package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/longln/simplebank/api"
	db "github.com/longln/simplebank/db/sqlc"
	"github.com/longln/simplebank/utils"
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
	server := api.NewServer(store)
	err = server.StartServer(config.ServerAddress)
	if err != nil {
		log.Fatal(err)
	}
}