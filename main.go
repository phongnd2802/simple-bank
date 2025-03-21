package main

import (
	"database/sql"
	"log"

	"github.com/phongnd2802/simple-bank/api"
	db "github.com/phongnd2802/simple-bank/db/sqlc"
	"github.com/phongnd2802/simple-bank/util"

	_ "github.com/lib/pq"
)



func main() {
	config, err := util.LoadConfig("./config")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}
	conn, err := sql.Open(config.DB.Driver, config.DB.Source)
	if err != nil {
		log.Fatalf("cannot connect db: %v", err)
	}	
	
	store := db.NewStore(conn)
	server := api.NewServer(store)


	err = server.Start(config.Server.Addr())
	if err != nil {
		log.Fatalf("cannot start server: %v", err)
	}
}