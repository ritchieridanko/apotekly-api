package utils

import (
	"github.com/ritchieridanko/apotekly-api/auth/config"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (hashedPassword string, err error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), config.GetBCryptCost())
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func ValidatePassword(hashedPassword, password string) (err error) {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
