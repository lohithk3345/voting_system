package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password []byte) ([]byte, error) {
	return hashPassword(password)
}

func hashPassword(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}

func VerifyHash(hash []byte, password string) error {
	return verifyHash(hash, password)
}

func verifyHash(hash []byte, password string) error {
	return bcrypt.CompareHashAndPassword(hash, []byte(password))
}
