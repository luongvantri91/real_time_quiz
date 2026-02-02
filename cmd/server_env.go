package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/you/rt-quiz/client/redis"
	"github.com/you/rt-quiz/infrastructure/postgres"
)

// Config holds all application configuration
type Config struct {
	// Server
	ServerPort string

	// Redis
	RedisURL string

	// PostgreSQL
	PostgresURL string

	// App
	Environment string
}

// ServerEnv holds all initialized server dependencies
type ServerEnv struct {
	Config       *Config
	RedisClient  redis.Client
	PostgresRepo *postgres.PostgresRepository
	Database     *sql.DB
	RedisURL     string // Store RedisURL for initialization
}

// getEnv retrieves environment variable with default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// InitializeConfig loads and returns configuration
func InitializeConfig() *Config {
	cfg := &Config{
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379/0"),
		PostgresURL: getEnv("POSTGRES_URL", "postgres://tri:123456@localhost:5432/rt_quiz?sslmode=disable"),
		Environment: getEnv("ENV", "development"),
	}
	log.Printf("Starting RT-Quiz server on port %s (env: %s)", cfg.ServerPort, cfg.Environment)
	return cfg
}

// String returns string representation of config
func (c *Config) String() string {
	return fmt.Sprintf("ServerPort=%s, Environment=%s", c.ServerPort, c.Environment)
}

// InitializeServerEnv initializes all environment dependencies
func InitializeServerEnv() *ServerEnv {
	// Load configuration
	cfg := InitializeConfig()

	// Initialize Redis (CRITICAL: Single source of truth for real-time)
	redisClient, err := redis.NewRedisClient(cfg.RedisURL)
	if err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}
	log.Println("✓ Redis initialized")

	// Initialize PostgreSQL (optional: only for final results persistence)
	var pgRepo *postgres.PostgresRepository
	var db *sql.DB
	if cfg.PostgresURL != "" {
		var err error
		db, err = sql.Open("postgres", cfg.PostgresURL)
		if err != nil {
			log.Fatalf("Failed to connect to PostgreSQL: %v", err)
		}
		if err := db.Ping(); err != nil {
			log.Fatalf("PostgreSQL connection failed: %v", err)
		}
		pgRepo = postgres.NewPostgresRepository(db)
		log.Println("✓ PostgreSQL initialized")
	} else {
		log.Println("⚠ PostgreSQL disabled (not storing results)")
	}

	return &ServerEnv{
		Config:       cfg,
		RedisClient:  redisClient,
		PostgresRepo: pgRepo,
		Database:     db,
		RedisURL:     cfg.RedisURL,
	}
}

// Close closes all resources
func (env *ServerEnv) Close() error {
	if env.RedisClient != nil {
		env.RedisClient.Close()
	}
	if env.Database != nil {
		env.Database.Close()
	}
	return nil
}
