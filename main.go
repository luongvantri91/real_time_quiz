package main

import (
	"log"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/you/rt-quiz/cmd"
)

func main() {
	// Initialize environment (config, Redis, PostgreSQL)
	env := cmd.InitializeServerEnv()
	defer env.Close()

	// Initialize all dependencies (services, handlers)
	deps := cmd.InitializeDependencies(env)

	// Create Echo server
	e := echo.New()

	// Setup all routes
	cmd.SetupRouter(e, deps)

	// Start server
	log.Printf("Server listening on :%s", env.Config.ServerPort)
	if err := e.Start(":" + env.Config.ServerPort); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
