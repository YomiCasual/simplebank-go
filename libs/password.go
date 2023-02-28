package lib

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)



type PasswordManager struct {}


func PasswordCrypt() *PasswordManager {
	return &PasswordManager{}
}




func (passwordManager *PasswordManager) HashPassword(password string)( string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	
	return string(hash), err

}

func (passwordManager *PasswordManager) CheckPassword(hashedPassword string, password string) (isPasswordValid bool) {
	if err :=  bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return false
	}

	return true
}

