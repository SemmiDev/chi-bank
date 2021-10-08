package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/SemmiDev/chi-bank/util/config"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := config.LoadEnv("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
