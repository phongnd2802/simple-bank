package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
)


// Payload contains the payload data of the token
type Payload struct {
	jwt.RegisteredClaims
	Username  string    `json:"username"`
}


// NewPayload creates a new token payload with a specific username and duration
func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenId.String(),
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(duration)},
			IssuedAt:  &jwt.NumericDate{Time: time.Now()},
			Issuer:    "SimpleBank",
		},
		Username: username,
	}

	return payload, nil
}
