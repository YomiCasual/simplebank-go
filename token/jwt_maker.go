package token

import (
	"errors"
	"fmt"
	lib "simplebank/libs"
	"time"

	"github.com/golang-jwt/jwt"
)

var minSecretKeyLength = 32

type JWTMaker struct {
	secretKey string
}


func NewJWTMaker(secretKey string) (Maker[int32], error) {

	if (len(secretKey) < minSecretKeyLength) {
		return nil, fmt.Errorf("invalid secret key size less than 32")
	}

	return &JWTMaker{secretKey: secretKey}, nil
}

func (maker *JWTMaker) 	CreateToken(username string, userId int32, duration time.Duration)(string, error)  {
	payload, err := NewPayload(username, userId, duration )

	if lib.HasError(err) {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	 return jwtToken.SignedString([]byte(maker.secretKey))

}


func (maker *JWTMaker)  VerifyToken(token string) (*Payload, error) {

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if (!ok) {
			return nil, ErrInvalidToken
		}

		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{},keyFunc )

	if lib.HasError(err) {
		verr, ok := err.(*jwt.ValidationError)

		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		if ok && errors.Is(verr.Inner, ErrInvalidToken) {
			return nil, ErrInvalidToken
		}

		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload) 

	if (!ok) {
		return nil, ErrInvalidToken
	}

	return payload, nil
}