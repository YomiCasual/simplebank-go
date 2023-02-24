package lib

import (
	"log"
	"math/rand"
	"strings"
	"time"
)


const alphabet = "abcdefghijklmnopqrstuvwxyz"

func HandleError(err error) {
	if (err != nil) {
		log.Fatal("Cannot connect to database",err)
	}
}


func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max - min + 1)
}

func RandomString(n int) string {
	var sb strings.Builder

	k := len(alphabet)

	for i:=0; i < n; i++ {
		c:= alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return RandomInt(0, 5000)
}

func RandomCurrency() string {
	currencies := []string{ "EUR", "USD", "CAD"}

	k := len(currencies);

	return currencies[rand.Intn(k)]

}