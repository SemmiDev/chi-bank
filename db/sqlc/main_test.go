package db

import (
	"database/sql"
	"github.com/SemmiDev/chi-bank/common"
	"log"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := common.LoadConfig("../..")
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
