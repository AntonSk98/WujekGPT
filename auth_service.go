package main

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

// Define possible authentication errors
var (
	ErrNotAuthenticated = errors.New("user not authenticated")
	ErrInvalidToken     = errors.New("invalid authentication token")
)

// Value to store in cache for authenticated users
const authenticatedValue = "authenticated"

// Token holds the components of the authentication token
type Token struct {
	prefix    string
	separator string
	username  string
	password  string
}

// AuthService handles authentication checks using a persistent cache
type AuthService struct {
	persistentCache *PersistentCache
	token           Token
}

// NewAuthService creates a new AuthService with the provided persistent cache
func NewAuthService(cache *PersistentCache, token Token) *AuthService {
	return &AuthService{
		persistentCache: cache,
		token:           token,
	}
}

// CheckAuth verifies if the given ID is authenticated or attempts to authenticate using the provided message
func (f *AuthService) CheckAuth(id, message string) error {
	if f.isAuthenticated(id) {
		return nil
	}

	if !f.isAuthRequest(message) {
		return ErrNotAuthenticated
	}

	if err := f.authenticate(id, message); err != nil {
		return err
	}

	return nil
}

// isAuthenticated checks if the given ID is authenticated
func (f *AuthService) isAuthenticated(id string) bool {
	val, found, err := f.persistentCache.Get(id)
	if err != nil {
		log.Printf("Error checking authentication for id %s: %v", id, err)
		return false
	}

	isAuthenticated := found && string(val) == authenticatedValue
	if isAuthenticated {
		return true
	}

	log.Printf("No authentication data found for id: %s", id)

	return false

}

// authenticate stores the authentication status for the given ID if the token is valid
func (f *AuthService) authenticate(id string, token string) error {
	if id == "" {
		return fmt.Errorf("cannot store authentication for empty ID")
	}

	expectedToken := f.token.prefix + f.token.username + f.token.separator + f.token.password

	if token != expectedToken {
		return ErrInvalidToken
	}

	err := f.persistentCache.Set(id, []byte(authenticatedValue), 30*24*time.Hour)
	if err != nil {
		return fmt.Errorf("failed to store authentication: %w", err)
	}

	return nil
}

// isAuthRequest checks if the provided string is a valid authentication request
func (f *AuthService) isAuthRequest(maybeToken string) bool {
	return strings.HasPrefix(maybeToken, f.token.prefix) &&
		strings.Contains(maybeToken, f.token.separator)
}
