package gapi

import (
	db "github.com/phongnd2802/simple-bank/db/sqlc"
	"github.com/phongnd2802/simple-bank/pb"
	"github.com/phongnd2802/simple-bank/token"
	"github.com/phongnd2802/simple-bank/util"
)

// Server serves gRPC requests for our banking service.
type Server struct {
	pb.UnimplementedSimpleBankServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

// NewServer creates a new gRPC server.
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.Token.Secret)
	if err != nil {
		return nil, err
	}

	server := &Server{
		config: config,
		store: store,
		tokenMaker: tokenMaker,
	}

	return server, nil
}
