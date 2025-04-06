package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/phongnd2802/simple-bank/util"
)

var testStore Store

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("local", "../../config")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}
	connPool, err := pgxpool.New(context.Background(), config.DB.Source)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}

	testStore = NewStore(connPool)

	os.Exit(m.Run())
}
