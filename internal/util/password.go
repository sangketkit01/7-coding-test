package util

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	
	return string(hashedPassword), err
}

func CheckPassword(hashedPassword, password string) error{
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}