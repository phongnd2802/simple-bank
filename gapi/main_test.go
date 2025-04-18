package gapi

import (
	"context"
	"fmt"
	"testing"
	"time"

	db "github.com/phongnd2802/simple-bank/db/sqlc"
	"github.com/phongnd2802/simple-bank/token"
	"github.com/phongnd2802/simple-bank/util"
	"github.com/phongnd2802/simple-bank/worker"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func newTestServer(t *testing.T, store db.Store, taskDistributor worker.TaskDistributor) *Server {
	config := util.Config{
		Token: util.Token{
			Secret:              util.RandomString(32),
			AccessTokenDuration: time.Minute,
		},
	}

	server, err := NewServer(config, store, taskDistributor)
	require.NoError(t, err)

	return server
}

func newContextWithBearerToken(t *testing.T, tokenMaker token.Maker, username string, duration time.Duration) context.Context {
	accessToken, _, err := tokenMaker.CreateToken(username, duration)
	require.NoError(t, err)

	bearerToken := fmt.Sprintf("%s %s", authorizationBearer, accessToken)
	md := metadata.MD{
		authorizationHeader: []string{
			bearerToken,
		},
	}

	return metadata.NewIncomingContext(context.Background(), md)
}
