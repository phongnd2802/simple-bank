package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/phongnd2802/simple-bank/util"
)

var (
	testQueries *Queries
	testDB      *sql.DB
)

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("local", "../../config")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}
	testDB, err = sql.Open(config.DB.Driver, config.DB.Source)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}

	testQueries = New(testDB)
	os.Exit(m.Run())
}
