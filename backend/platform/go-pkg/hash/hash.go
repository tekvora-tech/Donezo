// Package utils provides password hashing utilities using bcrypt.
package hash

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// ============================================================================
// HashConfig - Configuration for password hashing
// ============================================================================

// HashConfig holds configuration for bcrypt password hashing.
type HashConfig struct {
	// Cost is the bcrypt cost factor (4-31).
	// Higher = more secure but slower.
	// Default: 10 (recommended for production).
	// Development: 4 (faster for testing).
	Cost int
}

// DefaultHashConfig returns the default hash configuration.
func DefaultHashConfig() HashConfig {
	return HashConfig{Cost: 10}
}

// DevHashConfig returns a fast hash configuration for development/testing.
func DevHashConfig() HashConfig {
	return HashConfig{Cost: 4}
}

// ============================================================================
// HashPassword - Hash a plain text password
// ============================================================================

// HashPassword hashes a plain text password using bcrypt.
//
// Logic:
//   1. Generate a random salt internally (bcrypt handles this).
//   2. Hash the password with the configured cost factor.
//   3. Return the bcrypt hash string (includes salt + cost + hash).
//
// Parameters:
//   - password: plain text password to hash
//   - cfg: hash configuration (use DefaultHashConfig() for production)
//
// Returns:
//   - hashed password string (safe to store in database)
//   - error if hashing fails
//
// Usage:
//
//	hash, err := utils.HashPassword("SecurePass123", utils.DefaultHashConfig())
//	if err != nil {
//	    return err
//	}
//	// store 'hash' in users.password_hash
//
// Simpen di mana:
//   - Hashed password → kolom `password_hash` di table `users` (PostgreSQL).
//   - Plain password → TIDAK PERNAH disimpan di mana pun.
//   - Salt → sudah termasuk dalam bcrypt hash string.
func HashPassword(password string, cfg HashConfig) (string, error) {
	if password == "" {
		return "", fmt.Errorf("password cannot be empty")
	}

	// bcrypt.GenerateFromPassword automatically generates salt
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), cfg.Cost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hashedBytes), nil
}

// HashPasswordDefault hashes a password with default configuration (cost=10).
//
// Usage:
//
//	hash, err := utils.HashPasswordDefault("SecurePass123")
func HashPasswordDefault(password string) (string, error) {
	return HashPassword(password, DefaultHashConfig())
}

// ============================================================================
// ComparePassword - Verify a password against a hash
// ============================================================================

// ComparePassword verifies a plain text password against a bcrypt hash.
//
// Logic:
//   1. Extract salt and cost from the stored hash.
//   2. Hash the input password with the same salt and cost.
//   3. Compare the result with the stored hash.
//   4. Return nil if match, error if not match.
//
// Parameters:
//   - password: plain text password from user input
//   - hash: stored bcrypt hash from database
//
// Returns:
//   - nil if password matches
//   - error if password doesn't match or hash is invalid
//
// Usage:
//
//	err := utils.ComparePassword(inputPassword, storedHash)
//	if err != nil {
//	    // password salah
//	    return errors.New("invalid credentials")
//	}
//	// password benar
//
// Simpen di mana:
//   - storedHash → diambil dari `users.password_hash` (PostgreSQL).
//   - inputPassword → dari request body (tidak disimpan).
func ComparePassword(password, hash string) error {
	if password == "" || hash == "" {
		return bcrypt.ErrMismatchedHashAndPassword
	}

	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// ComparePasswordBool returns boolean instead of error.
//
// Usage:
//
//	if !utils.ComparePasswordBool(inputPassword, storedHash) {
//	    return errors.New("invalid credentials")
//	}
func ComparePasswordBool(password, hash string) bool {
	return ComparePassword(password, hash) == nil
}

// ============================================================================
// Hash Validation
// ============================================================================

// IsValidHash checks if a string is a valid bcrypt hash format.
//
// Usage:
//
//	if !utils.IsValidHash(storedHash) {
//	    return errors.New("invalid password hash in database")
//	}
func IsValidHash(hash string) bool {
	// bcrypt hash format: $2a$<cost>$<salt+hash>
	// Minimum length: 59 characters
	if len(hash) < 59 {
		return false
	}
	return hash[:4] == "$2a$" || hash[:4] == "$2b$" || hash[:4] == "$2y$" || hash[:4] == "$2x$"
}

// GetHashCost extracts the cost factor from a bcrypt hash.
//
// Usage:
//
//	cost, err := utils.GetHashCost(storedHash)
//	if err != nil {
//	    return err
//	}
//	if cost < 10 {
//	    // rehash with higher cost
//	}
func GetHashCost(hash string) (int, error) {
	return bcrypt.Cost([]byte(hash))
}

// NeedsRehash checks if a password hash needs to be rehashed with a different cost.
//
// Usage:
//
//	if utils.NeedsRehash(storedHash, utils.DefaultHashConfig()) {
//	    // rehash password with new cost
//	    newHash, _ := utils.HashPassword(password, utils.DefaultHashConfig())
//	    // update database
//	}
func NeedsRehash(hash string, cfg HashConfig) bool {
	currentCost, err := bcrypt.Cost([]byte(hash))
	if err != nil {
		return true
	}
	return currentCost != cfg.Cost
}

// ============================================================================
// Password Validation
// ============================================================================

// ValidatePassword checks if a password meets security requirements.
//
// Rules:
//   - Minimum 8 characters
//   - At least 1 uppercase letter
//   - At least 1 lowercase letter
//   - At least 1 digit
//   - At least 1 special character (!@#$%^&* etc.)
//
// Returns:
//   - nil if password is valid
//   - error with specific validation message
//
// Usage:
//
//	if err := utils.ValidatePassword("SecurePass123!"); err != nil {
//	    response.BadRequest(c, "Validation failed", err.Error())
//	    return
//	}
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasDigit   bool
		hasSpecial bool
	)

	for _, ch := range password {
		switch {
		case ch >= 'A' && ch <= 'Z':
			hasUpper = true
		case ch >= 'a' && ch <= 'z':
			hasLower = true
		case ch >= '0' && ch <= '9':
			hasDigit = true
		case ch >= 33 && ch <= 126:
			hasSpecial = true
		}
	}

	if !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	if !hasDigit {
		return fmt.Errorf("password must contain at least one digit")
	}
	if !hasSpecial {
		return fmt.Errorf("password must contain at least one special character")
	}

	return nil
}

// ValidatePasswordSimple checks only minimum length (for less strict requirements).
func ValidatePasswordSimple(password string, minLength int) error {
	if len(password) < minLength {
		return fmt.Errorf("password must be at least %d characters", minLength)
	}
	return nil
}