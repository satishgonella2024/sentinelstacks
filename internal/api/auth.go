package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// UserClaims contains the JWT claims for authentication
type UserClaims struct {
	Username string `json:"username"`
	Admin    bool   `json:"admin"`
	jwt.RegisteredClaims
}

// contextKey is a custom type for context keys
type contextKey string

// Context keys
const (
	userContextKey contextKey = "user"
)

// authMiddleware is middleware for JWT authentication
func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			s.sendError(w, http.StatusUnauthorized, "Authorization header required")
			return
		}

		// Check that the header is in the right format
		if !strings.HasPrefix(authHeader, "Bearer ") {
			s.sendError(w, http.StatusUnauthorized, "Invalid authorization format, expected Bearer token")
			return
		}

		// Extract the token string
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and validate the token
		token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Make sure we're using the expected signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			// Return the secret key used for signing
			return []byte(s.config.TokenAuthSecret), nil
		})

		if err != nil {
			s.sendError(w, http.StatusUnauthorized, fmt.Sprintf("Invalid token: %v", err))
			return
		}

		// Check if the token is valid
		if !token.Valid {
			s.sendError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		// Extract user claims
		claims, ok := token.Claims.(*UserClaims)
		if !ok {
			s.sendError(w, http.StatusUnauthorized, "Invalid token claims")
			return
		}

		// Store user information in the request context
		ctx := context.WithValue(r.Context(), userContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GenerateToken generates a new JWT token for a user
func (s *Server) GenerateToken(username string, admin bool) (string, error) {
	if s.config.TokenAuthSecret == "" {
		return "", fmt.Errorf("token auth secret not configured")
	}

	// Set expiration time
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create claims
	claims := &UserClaims{
		Username: username,
		Admin:    admin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "sentinelstacks-api",
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret
	tokenString, err := token.SignedString([]byte(s.config.TokenAuthSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// getUserFromContext extracts the user claims from the request context
func getUserFromContext(ctx context.Context) (*UserClaims, bool) {
	user, ok := ctx.Value(userContextKey).(*UserClaims)
	return user, ok
}

// requireAdmin is middleware that ensures the user is an admin
func (s *Server) requireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := getUserFromContext(r.Context())
		if !ok || !user.Admin {
			s.sendError(w, http.StatusForbidden, "Admin privileges required")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// handleLogin handles user login and token generation
func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	// This is a simplified login handler
	// In a real implementation, you'd validate user credentials against a database

	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		s.sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// TODO: Implement proper authentication
	// For now, accept any credentials for development

	isAdmin := credentials.Username == "admin"

	// Generate token
	token, err := s.GenerateToken(credentials.Username, isAdmin)
	if err != nil {
		s.sendError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	s.sendJSON(w, http.StatusOK, map[string]string{
		"token":    token,
		"username": credentials.Username,
		"expires":  time.Now().Add(24 * time.Hour).Format(time.RFC3339),
	})
}
