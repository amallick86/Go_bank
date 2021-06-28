package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/dgrijalva/jwt-go"
)

const minSecretKeySize = 32

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}

	return &JWTMaker{secretKey},nil
}


//CreateToken creates a new token for a specific username and duration
func (make *JWTMaker) CreateToken(username string, duration time.Duration) (string , error) {
	payload,err := NewPayoad(username, duration)
	if err != nil {
		return "",err
	}
	jwtToken:= jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	return jwtToken.SignedString([]byte(make.secretKey))
}

//verifyToken checks if the token is valid or not
func (make *JWTMaker) verifyToken(token string)(*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(make.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token , &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil,ErrExpiredToken
		}
	}

		payload, ok := jwtToken.Claims.(*Payload)
		if !ok {
			return nil, ErrInvalidToken
		}
		return payload, nil
	
}