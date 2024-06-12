package crypto

import (
	"fmt"
	"strconv"
	"sync"

	configs "beli_mang/cfg"

	"golang.org/x/crypto/bcrypt"
)

var (
	bcryptCost     int
	bcryptCostOnce sync.Once
)

// initBcryptCost initializes the bcrypt cost once, making subsequent password hashing faster.
func initBcryptCost() {
	costStr := configs.GetEnvOrDefault("BCRYPT_SALT", "10")
	cost, err := strconv.Atoi(costStr)
	if err != nil {
		// Fallback to default cost of 10 if the environment variable is not a valid integer.
		cost = 10
	}
	bcryptCost = cost
}

// GenerateHashedPassword hashes a password using bcrypt with a cost specified by the configuration.
func GenerateHashedPassword(password string) (string, error) {
	// Ensure bcrypt cost is initialized only once.
	bcryptCostOnce.Do(initBcryptCost)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", fmt.Errorf("failed to generate hashed password: %w", err)
	}

	return string(hashedPassword), nil
}

// VerifyPassword compares a plaintext password with a hashed password and returns an error if they don't match.
func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
