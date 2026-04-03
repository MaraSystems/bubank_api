package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/MaraSystems/bubank_api/utils"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var testDB *sql.DB
var testQueries *Queries

func TestMain(m *testing.M) {
	var err error

	config, err := utils.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to the database: ", err)
	}

	testQueries = New(testDB)
	os.Exit(m.Run())
}
