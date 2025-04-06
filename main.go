package main

import (
	"context"
	"embed"
	"errors"
	"io/fs"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"
	"github.com/phongnd2802/simple-bank/api"
	db "github.com/phongnd2802/simple-bank/db/sqlc"
	"github.com/phongnd2802/simple-bank/email"
	"github.com/phongnd2802/simple-bank/gapi"
	"github.com/phongnd2802/simple-bank/pb"
	"github.com/phongnd2802/simple-bank/util"
	"github.com/phongnd2802/simple-bank/worker"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

var interuptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

//go:embed docs/swagger/*
var fileSwagger embed.FS

func main() {
	mode := os.Getenv("MODE")
	if mode == "" {
		mode = "local"
	}

	if mode == "local" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	config, err := util.LoadConfig(mode, "./config")
	if err != nil {
		log.Fatal().Msgf("cannot load config: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), interuptSignals...)
	defer stop()

	conn, err := pgxpool.New(ctx, config.DB.Source)
	if err != nil {
		log.Fatal().Msgf("cannot connect db: %v", err)
	}

	store := db.NewStore(conn)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.Redis.Addr(),
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)
	mailer := email.NewGmailSender(config.Email.EmailSenderName, config.Email.EmailSenderAddress, config.Email.EmailSenderPassword)

	waitGroup, ctx := errgroup.WithContext(ctx)

	runTaskProcessor(ctx, waitGroup, redisOpt, store, mailer)
	runGatewayServer(ctx, waitGroup, config, store, taskDistributor)
	runGrpcServer(ctx, waitGroup, config, store, taskDistributor)

	err = waitGroup.Wait()
	if err != nil {
		log.Fatal().Err(err).Msg("error from wait group")
	}
}

func runTaskProcessor(ctx context.Context, waitGroup *errgroup.Group, redisOpt asynq.RedisClientOpt, store db.Store, mailer email.EmailSender) {
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, mailer)
	log.Info().Msg("Start task processor")
	waitGroup.Go(func() error {
		err := taskProcessor.Start()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to start task processor")
			return err
		}

		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("shutting down task processor")
		taskProcessor.Shutdown()
		log.Info().Msg("task processor exited")
		return nil
	})
}

func runGrpcServer(ctx context.Context, waitGroup *errgroup.Group, config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Msgf("cannot create new server: %v", err)
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.Server.Grpc.Addr())
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener")
	}

	waitGroup.Go(func() error {
		log.Info().Msgf("start gRPC server at %s", listener.Addr().String())
		err = grpcServer.Serve(listener)
		if err != nil {
			log.Fatal().Msg("cannot create gRPC server")
			return err
		}

		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("shutting down gRPC server")
		grpcServer.GracefulStop()
		log.Info().Msg("gRPC server exited")
		return nil
	})
}

func runGatewayServer(ctx context.Context, waitGroup *errgroup.Group, config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Msgf("cannot create new server: %v", err)
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
	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Msg("cannot register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	swaggerFS, _ := fs.Sub(fileSwagger, "docs/swagger")
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", http.FileServer(http.FS(swaggerFS))))

	listener, err := net.Listen("tcp", config.Server.Http.Addr())
	if err != nil {
		log.Fatal().Msg("cannot create listener")
	}

	log.Info().Msgf("start HTTP gateway server at %s", listener.Addr().String())
	handler := gapi.HttpLogger(mux)

	waitGroup.Go(func() error {
		err = http.Serve(listener, handler)
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}
			log.Fatal().Msg("cannot create HTTP gateway server")
			return err
		}

		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("shutting down HTTP gateway server")
		err = listener.Close()
		if err != nil {
			log.Fatal().Msg("cannot close HTTP gateway server")
			return err
		}
		log.Info().Msg("HTTP gateway server exited")
		return nil
	})
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Msgf("cannot create new server: %v", err)
	}

	err = server.Start(config.Server.Http.Addr())
	if err != nil {
		log.Fatal().Msgf("cannot start server: %v", err)
	}
}
