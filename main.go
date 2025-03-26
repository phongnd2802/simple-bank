package main

import (
	"context"
	"database/sql"
	"embed"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/phongnd2802/simple-bank/api"
	db "github.com/phongnd2802/simple-bank/db/sqlc"
	"github.com/phongnd2802/simple-bank/gapi"
	"github.com/phongnd2802/simple-bank/pb"
	"github.com/phongnd2802/simple-bank/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

	_ "github.com/lib/pq"
)

//go:embed docs/swagger/*
var fileSwagger embed.FS

func main() {
	mode := os.Getenv("MODE")
	if mode == "" {
		mode = "local"
	}

	config, err := util.LoadConfig(mode, "./config")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}
	conn, err := sql.Open(config.DB.Driver, config.DB.Source)
	if err != nil {
		log.Fatalf("cannot connect db: %v", err)
	}

	store := db.NewStore(conn)

	go runGatewayServer(config, store)
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

func runGatewayServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatalf("cannot create new server: %v", err)
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("cannot register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	swaggerFS, _ := fs.Sub(fileSwagger, "docs/swagger")
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", http.FileServer(http.FS(swaggerFS))))

	listener, err := net.Listen("tcp", config.Server.Http.Addr())
	if err != nil {
		log.Fatal("cannot create listener")
	}

	log.Printf("start HTTP gateway server at %s", listener.Addr().String())

	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("cannot create HTTP gateway server")
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
