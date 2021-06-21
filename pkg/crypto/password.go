package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Hash password
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("Failed to hash password %s", password)
	}

	return string(hashedPassword), nil
}

// Check password
func CheckPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
