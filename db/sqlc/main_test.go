package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
    host     = "localhost"
    port     = 5432
    user     = "root"
    password = "secret"
    dbname   = "simple_bank"
)

var testDB *sql.DB
var testQueries *Queries

func TestMain(m *testing.M) {
    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)

    var err error
    testDB, err = sql.Open("postgres", psqlInfo)
    if err != nil {
        log.Fatal(err)
    }
    defer testDB.Close()

    err = testDB.Ping()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Successfully connected to Database!")

    testQueries = New(testDB)

	os.Exit(m.Run())
}