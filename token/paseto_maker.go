package token

import (
	"fmt"
	lib "simplebank/libs"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)


type PasetoMaker struct {
	paseto *paseto.V2
	symmetricKey []byte
}


func NewPasetoMaker(symmetricKey string) (Maker[int32], error) {
	if (len(symmetricKey) < chacha20poly1305.KeySize) {
		return nil, fmt.Errorf("invalid secret key size less than 32")
	}

	maker := &PasetoMaker{
		symmetricKey: []byte(symmetricKey),
		paseto: paseto.NewV2(),
	}

	return maker, nil
}

func (maker *PasetoMaker) 	CreateToken(username string, userId int32, duration time.Duration)(string, error)  {
	payload, err := NewPayload(username, userId, duration )

	if lib.HasError(err) {
		return "", err
	}

	return maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
		

}


func (maker *PasetoMaker)  VerifyToken(token string) (*Payload, error) {

	payload := &Payload{}

	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)

	if lib.HasError(err) {
		return nil, ErrInvalidToken
	}

	err = payload.Valid() 
	if lib.HasError(err) {
		return nil, err
	}

	return payload, nil
}