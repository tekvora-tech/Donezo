// Package jwt provides JWT token generation, validation, and context helpers
// for the Donezo application.
package jwt

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// ============================================================================
// Constants
// ============================================================================

const (
	// TokenTypeAccess identifies an access token in the "type" claim.
	TokenTypeAccess = "access"
	// TokenTypeRefresh identifies a refresh token in the "type" claim.
	TokenTypeRefresh = "refresh"

	// ContextKeyUserID is the key used to store user_id in context.
	ContextKeyUserID = "user_id"
	// ContextKeyEmail is the key used to store email in context.
	ContextKeyEmail = "email"
)

// ============================================================================
// Errors
// ============================================================================

var (
	ErrInvalidToken     = errors.New("token tidak valid")
	ErrExpiredToken     = errors.New("token sudah expired")
	ErrInvalidSignature = errors.New("signature token tidak valid")
	ErrInvalidType      = errors.New("tipe token tidak valid")
	ErrMissingToken     = errors.New("token tidak ditemukan")
	ErrMalformedToken   = errors.New("format token tidak valid")
)

// ============================================================================
// Config
// ============================================================================

// Config holds JWT configuration.
type Config struct {
	// Secret is the HMAC secret key used for signing tokens.
	// Must be at least 32 characters for HS256.
	Secret string

	// AccessTokenTTL is the lifetime of an access token.
	// Default: 15 minutes for production.
	AccessTokenTTL time.Duration

	// RefreshTokenTTL is the lifetime of a refresh token.
	// Default: 7 days.
	RefreshTokenTTL time.Duration

	// Issuer identifies the token issuer (your app name).
	// Default: "donezo".
	Issuer string
}

// Validate checks if the configuration is valid.
func (c *Config) Validate() error {
	if c.Secret == "" {
		return errors.New("JWT secret tidak boleh kosong")
	}
	if len(c.Secret) < 32 {
		return errors.New("JWT secret minimal 32 karakter")
	}
	if c.AccessTokenTTL <= 0 {
		c.AccessTokenTTL = 15 * time.Minute
	}
	if c.RefreshTokenTTL <= 0 {
		c.RefreshTokenTTL = 7 * 24 * time.Hour
	}
	if c.Issuer == "" {
		c.Issuer = "donezo"
	}
	return nil
}

// ============================================================================
// TokenPair
// ============================================================================

// TokenPair holds both access and refresh tokens returned after login/register.
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // seconds until access token expires
}

// ============================================================================
// Claims
// ============================================================================

// CustomClaims extends jwt.RegisteredClaims with application-specific claims.
type CustomClaims struct {
	jwt.RegisteredClaims

	// Email is the user's email address.
	Email string `json:"email"`

	// Type indicates whether this is an access or refresh token.
	Type string `json:"type"`
}

// ============================================================================
// Service
// ============================================================================

// Service defines the contract for JWT operations.
type Service interface {
	// GenerateTokenPair creates a new access token and refresh token for a user.
	GenerateTokenPair(userID uuid.UUID, email string) (*TokenPair, error)

	// GenerateAccessToken creates only an access token (used after refresh).
	GenerateAccessToken(userID uuid.UUID, email string) (string, error)

	// ValidateAccessToken parses and validates an access token.
	// Returns the claims if valid.
	ValidateAccessToken(tokenString string) (*CustomClaims, error)

	// ValidateRefreshToken parses and validates a refresh token.
	// Returns the claims if valid.
	ValidateRefreshToken(tokenString string) (*CustomClaims, error)

	// ExtractBearerToken extracts the token string from "Bearer <token>" header.
	ExtractBearerToken(authHeader string) (string, error)

	// GetTokenRemainingTime returns the remaining time until token expiration.
	GetTokenRemainingTime(tokenString string) (time.Duration, error)
}

// ============================================================================
// Service Implementation
// ============================================================================

type service struct {
	config Config
}

// NewService creates a new JWT service with the given configuration.
func NewService(cfg Config) (Service, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid jwt config: %w", err)
	}
	return &service{config: cfg}, nil
}

// GenerateTokenPair creates both access and refresh tokens for a user.
func (s *service) GenerateTokenPair(userID uuid.UUID, email string) (*TokenPair, error) {
	accessToken, err := s.generateToken(userID, email, TokenTypeAccess, s.config.AccessTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("gagal generate access token: %w", err)
	}

	refreshToken, err := s.generateToken(userID, email, TokenTypeRefresh, s.config.RefreshTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("gagal generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.config.AccessTokenTTL.Seconds()),
	}, nil
}

// GenerateAccessToken creates a single access token (used during token refresh).
func (s *service) GenerateAccessToken(userID uuid.UUID, email string) (string, error) {
	return s.generateToken(userID, email, TokenTypeAccess, s.config.AccessTokenTTL)
}

// generateToken is the internal method that creates a JWT with custom claims.
func (s *service) generateToken(userID uuid.UUID, email, tokenType string, ttl time.Duration) (string, error) {
	now := time.Now().UTC()

	claims := CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			Issuer:    s.config.Issuer,
			NotBefore: jwt.NewNumericDate(now),
		},
		Email:  email,
		Type:   tokenType,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.Secret))
	if err != nil {
		return "", fmt.Errorf("gagal sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateAccessToken parses and validates an access token string.
func (s *service) ValidateAccessToken(tokenString string) (*CustomClaims, error) {
	return s.validateToken(tokenString, TokenTypeAccess)
}

// ValidateRefreshToken parses and validates a refresh token string.
func (s *service) ValidateRefreshToken(tokenString string) (*CustomClaims, error) {
	return s.validateToken(tokenString, TokenTypeRefresh)
}

// validateToken is the internal validation logic shared by access and refresh token validation.
func (s *service) validateToken(tokenString, expectedType string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.Secret), nil
	})

	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, ErrExpiredToken
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			return nil, ErrInvalidSignature
		default:
			return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
		}
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	if claims.Type != expectedType {
		return nil, fmt.Errorf("%w: expected %s, got %s", ErrInvalidType, expectedType, claims.Type)
	}

	_, err = uuid.Parse(claims.Subject)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid user_id format", ErrInvalidToken)
	}

	return claims, nil
}

// ExtractBearerToken extracts the token from an Authorization header.
func (s *service) ExtractBearerToken(authHeader string) (string, error) {
	if authHeader == "" {
		return "", ErrMissingToken
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", ErrMalformedToken
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", ErrMissingToken
	}

	return token, nil
}

// GetTokenRemainingTime returns the time remaining until a token expires.
func (s *service) GetTokenRemainingTime(tokenString string) (time.Duration, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &CustomClaims{})
	if err != nil {
		return 0, fmt.Errorf("gagal parse token: %w", err)
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return 0, ErrInvalidToken
	}

	if claims.ExpiresAt == nil {
		return 0, ErrInvalidToken
	}

	return time.Until(claims.ExpiresAt.Time), nil
}

// ============================================================================
// Context Helpers
// ============================================================================

// ContextWithUserID stores the userID in the context.
func ContextWithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, ContextKeyUserID, userID)
}

// UserIDFromContext retrieves the userID from context.
// Returns uuid.Nil if not found.
func UserIDFromContext(ctx context.Context) uuid.UUID {
	userID, ok := ctx.Value(ContextKeyUserID).(uuid.UUID)
	if !ok {
		return uuid.Nil
	}
	return userID
}

// ContextWithEmail stores the email in the context.
func ContextWithEmail(ctx context.Context, email string) context.Context {
	return context.WithValue(ctx, ContextKeyEmail, email)
}

// EmailFromContext retrieves the email from context.
func EmailFromContext(ctx context.Context) string {
	email, ok := ctx.Value(ContextKeyEmail).(string)
	if !ok {
		return ""
	}
	return email
}