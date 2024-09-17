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

var TestDB *sql.DB
var Query *Queries

func TestMain(m *testing.M) {
    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)

    var err error
    TestDB, err = sql.Open("postgres", psqlInfo)
    if err != nil {
        log.Fatal(err)
    }
    defer TestDB.Close()

    err = TestDB.Ping()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Successfully connected to Database!")

    Query = New(TestDB)

	os.Exit(m.Run())
}