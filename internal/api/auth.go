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

// contextKey is used for context values
type contextKey string

// Context keys
const (
	userContextKey contextKey = "user"
)

// LoginRequest represents a login request
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// User represents a user in the system
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

// Claims represents the JWT claims
type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// authMiddleware is middleware for JWT authentication
func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			s.sendError(w, http.StatusUnauthorized, "Authorization header is required")
			return
		}

		// Check that it's a Bearer token
		if !strings.HasPrefix(authHeader, "Bearer ") {
			s.sendError(w, http.StatusUnauthorized, "Invalid authorization format")
			return
		}

		// Extract token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// Return secret key
			secret := s.config.TokenAuthSecret
			if secret == "" {
				secret = "sentinel-default-secret-key-change-in-production"
			}
			return []byte(secret), nil
		})

		// Handle parsing errors
		if err != nil {
			s.log.Printf("Token validation error: %v", err)
			s.sendError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		// Check if token is valid
		if !token.Valid {
			s.sendError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		// Store claims in context for use by handlers
		ctx := r.Context()
		ctx = context.WithValue(ctx, "user", claims)

		// Call next handler with updated context
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

// @Summary User login
// @Description Authenticate a user and get a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param user body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// For demo purposes, accept any non-empty username/password
	if req.Username == "" || req.Password == "" {
		s.sendError(w, http.StatusBadRequest, "Username and password are required")
		return
	}

	// In a real implementation, validate credentials against a database
	// For now, mock a successful login for any non-empty credentials

	// Create user data
	user := User{
		ID:       "user-1",
		Username: req.Username,
		Email:    fmt.Sprintf("%s@example.com", req.Username),
		Role:     "admin", // For demo purposes, make everyone an admin
	}

	// Create JWT token
	token, err := s.createJWTToken(user)
	if err != nil {
		s.log.Printf("Error creating JWT token: %v", err)
		s.sendError(w, http.StatusInternalServerError, "Error creating authentication token")
		return
	}

	// Return token and user info
	resp := LoginResponse{
		Token: token,
		User:  user,
	}

	s.sendJSON(w, http.StatusOK, resp)
}

// createJWTToken creates a new JWT token for a user
func (s *Server) createJWTToken(user User) (string, error) {
	// Use config secret or a default one
	secret := s.config.TokenAuthSecret
	if secret == "" {
		secret = "sentinel-default-secret-key-change-in-production"
	}

	// Create claims
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "sentinel-api",
			Subject:   user.ID,
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	return token.SignedString([]byte(secret))
}
