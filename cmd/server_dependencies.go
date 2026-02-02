package cmd

import (
	"github.com/you/rt-quiz/client/redis"
	"github.com/you/rt-quiz/client/ws"
	"github.com/you/rt-quiz/handlers"
	"github.com/you/rt-quiz/repository"
	"github.com/you/rt-quiz/services"
)

// RepositoryDependencies holds all repository dependencies
type RepositoryDependencies struct {
	ParticipantRepo repository.ParticipantRepository
	AnswerRepo      repository.AnswerRepository
	ResultRepo      repository.ResultRepository
	QuizRepo        repository.QuizRepository
}

// InitializeRepositories sets up all repository interfaces from PostgreSQL
func InitializeRepositories(env *ServerEnv) *RepositoryDependencies {
	if env.PostgresRepo == nil {
		return &RepositoryDependencies{} // Empty if no DB configured
	}

	return &RepositoryDependencies{
		ParticipantRepo: env.PostgresRepo,
		AnswerRepo:      env.PostgresRepo,
		ResultRepo:      env.PostgresRepo,
		QuizRepo:        env.PostgresRepo,
	}
}

// ServerDependencies holds all application service dependencies
type ServerDependencies struct {
	QuizService        *services.QuizService
	ParticipantService *services.ParticipantService
	HTTPHandler        *handlers.HTTPHandler
	WSHandler          *handlers.WSHandler
}

// InitializeDependencies sets up all service layer and handler dependencies
func InitializeDependencies(env *ServerEnv) *ServerDependencies {
	// Initialize repositories
	repos := InitializeRepositories(env)

	// Create client implementations directly from environment
	redisClient, err := redis.NewRedisClient(env.RedisURL)
	if err != nil {
		panic("Failed to initialize Redis client: " + err.Error())
	}
	wsServer := ws.NewWebSocketServer()

	// Initialize WebSocket handler
	wsHandler := handlers.NewWSHandler(wsServer, redisClient)

	// Initialize service layer with dependencies
	quizService := services.NewQuizService(repos.QuizRepo, repos.ResultRepo, redisClient)
	participantService := services.NewParticipantService(repos.ParticipantRepo, repos.AnswerRepo, redisClient)

	// Initialize HTTP handler with services (clean architecture)
	httpHandler := handlers.NewHTTPHandlerWithServices(quizService, participantService, redisClient)

	return &ServerDependencies{
		QuizService:        quizService,
		ParticipantService: participantService,
		HTTPHandler:        httpHandler,
		WSHandler:          wsHandler,
	}
}
