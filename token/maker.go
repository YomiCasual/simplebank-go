package token

import (
	"time"

	"golang.org/x/exp/constraints"
)




type  Maker[T constraints.Ordered] interface {
	CreateToken(username string, userId T, duration time.Duration)(string, error) 
	VerifyToken(token string) (*Payload, error)
}