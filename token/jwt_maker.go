package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)


const minSecretKeySize = 32

// JWTMaker is a JSON Web Token Maker
type JWTMaker struct {
	secretKey string
}

// CreateToken implements Maker.
func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", payload, err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(maker.secretKey))
	return token, payload, err
}

// VertifyToken implements Maker.
func (maker *JWTMaker) VertifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token)  (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}

		return []byte(maker.secretKey), nil
	}
	
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}


func NewJWTMaker(secretKey string) (*JWTMaker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}

	return &JWTMaker{
		secretKey: secretKey,
	}, nil
}

var _ Maker = (*JWTMaker)(nil)
