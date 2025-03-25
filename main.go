package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/phongnd2802/simple-bank/api"
	db "github.com/phongnd2802/simple-bank/db/sqlc"
	"github.com/phongnd2802/simple-bank/gapi"
	"github.com/phongnd2802/simple-bank/pb"
	"github.com/phongnd2802/simple-bank/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig("local", "./config")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}
	conn, err := sql.Open(config.DB.Driver, config.DB.Source)
	if err != nil {
		log.Fatalf("cannot connect db: %v", err)
	}

	store := db.NewStore(conn)

	runGrpcServer(config, store)
}

func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatalf("cannot create new server: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.Server.Grpc.Addr())
	if err != nil {
		log.Fatal("cannot create listener")
	}

	log.Printf("start gRPC server at %s", listener.Addr().String())

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot create gRPC server")
	}
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatalf("cannot create new server: %v", err)
	}

	err = server.Start(config.Server.Http.Addr())
	if err != nil {
		log.Fatalf("cannot start server: %v", err)
	}
}
