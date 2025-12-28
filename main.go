package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/asaskevich/EventBus"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	// Load .env file if it exists
	_ = godotenv.Load()

	eventBus := EventBus.New()

	// Read required configuration from environment variables
	cachePath := requireEnvVar("CACHE_PATH")
	dbPath := requireEnvVar("DB_PATH")
	authTokenPrefix := requireEnvVar("AUTH_TOKEN_PREFIX")
	authTokenSeparator := requireEnvVar("AUTH_TOKEN_SEPARATOR")
	authUsername := requireEnvVar("AUTH_USERNAME")
	authPassword := requireEnvVar("AUTH_PASSWORD")

	persistentCache := NewCache(cachePath)

	authFilter := NewAuthService(persistentCache, Token{
		prefix:    authTokenPrefix,
		separator: authTokenSeparator,
		username:  authUsername,
		password:  authPassword,
	})

	NewWhatsAppGateway(eventBus)

	app := NewWhatsmeowApp(dbPath, authFilter, eventBus)

	fmt.Println("App running. Press Ctrl+C to exit.")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	<-sigChan

	fmt.Println("\nShutting down...")

	persistentCache.Close()
	app.Close()
}

// requireEnvVar retrieves an environment variable or exits with a fatal error if not found
func requireEnvVar(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("environment variable %s is required but not set", key)
	}
	return value
}
