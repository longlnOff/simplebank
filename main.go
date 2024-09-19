package main

import (
	"database/sql"
	"fmt"
	"log"
	"github.com/longln/simplebank/api"
	_ "github.com/lib/pq"
	db "github.com/longln/simplebank/db/sqlc"
)

const (
    host     = "localhost"
    port     = 5432
    user     = "root"
    password = "secret"
    dbname   = "simple_bank"
	address	 = "0.0.0.0:8080"
)

func main() {
    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)

    conn, err := sql.Open("postgres", psqlInfo)
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
	err = server.StartServer(address)
	if err != nil {
		log.Fatal(err)
	}
}