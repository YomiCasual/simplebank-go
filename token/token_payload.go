package token

import (
	"errors"
	lib "simplebank/libs"
	"time"

	"github.com/google/uuid"
)


var  (
ErrExpiredToken = errors.New("token has expired")
ErrInvalidToken = errors.New("token is invalid")
)

type Payload struct {
	UUID uuid.UUID `json:"uid"`
	UserId int32 `json:"user_id"`
	Username string `json:"username"`
	IssuedAt time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`

	
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}

	return nil
}

func NewPayload(username string, userId int32, duration time.Duration) (*Payload, error) {
		tokenID, err := uuid.NewRandom()

		if lib.HasError(err) {
			return nil, err
		}


		payload := &Payload{
			UUID: tokenID,
			UserId: userId,
			Username: username,
			IssuedAt: time.Now(),
			ExpiredAt: time.Now().Add(duration),
		}

		return payload, nil
}