package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password %s", err)
	}

	return string(hashedPassword), nil
}

func CheckHashPassword(password string, hasedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hasedPassword), []byte(password))
}