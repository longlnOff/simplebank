package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	_ "github.com/lib/pq"
	"github.com/longln/simplebank/utils"
)



var testDB *sql.DB
var testQueries *Queries

func TestMain(m *testing.M) {
    config, err := utils.LoadConfig("../..")
    if err != nil {
        log.Fatal("cannot read config:", err)
    }
    testDB, err = sql.Open(config.DBDriver, config.DBSource)
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