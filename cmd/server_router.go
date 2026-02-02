package cmd

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// SetupRouter configures all API routes
func SetupRouter(e *echo.Echo, deps *ServerDependencies) {
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Content-Type"},
	}))

	// ===== Admin APIs (Quiz Lifecycle Management) =====
	// Admin: Create quiz (pending status)
	e.POST("/admin/quizzes", deps.HTTPHandler.CreateQuiz)

	// Admin: Get specific quiz details
	e.GET("/admin/quizzes/:quizId", deps.HTTPHandler.GetQuiz)

	// Admin: List all quizzes
	e.GET("/admin/quizzes", deps.HTTPHandler.ListQuizzes)

	// Admin: Get quiz status and time remaining
	e.GET("/admin/quizzes/:quizId/status", deps.HTTPHandler.GetQuizStatus)

	// Admin: Initialize quiz (pending → started, setup Redis)
	e.POST("/admin/quizzes/:quizId/init", deps.HTTPHandler.InitQuiz)

	// Admin: End quiz (started → ended, save final results)
	e.POST("/admin/quizzes/:quizId/end", deps.HTTPHandler.EndQuiz)

	// Admin: Question Management
	e.POST("/admin/quizzes/:quizId/questions", deps.HTTPHandler.AddQuestion)
	e.DELETE("/admin/quizzes/:quizId/questions/:questionId", deps.HTTPHandler.DeleteQuestion)
	e.GET("/admin/quizzes/:quizId/questions", deps.HTTPHandler.ListQuestions)

	// ===== User APIs (4 Core Endpoints) =====
	// 1. Join quiz
	e.POST("/quizzes/:quizId/join", deps.HTTPHandler.JoinQuiz)

	// 2. Submit answer (atomic Redis update)
	e.POST("/quizzes/:quizId/answer", deps.HTTPHandler.SubmitAnswer)

	// 3. Get leaderboard (from Redis)
	e.GET("/quizzes/:quizId/leaderboard", deps.HTTPHandler.GetLeaderboard)

	// 4. WebSocket for real-time leaderboard updates
	e.GET("/quizzes/:quizId/ws", deps.WSHandler.HandleWebSocket)

	// Health check
	e.GET("/health", deps.HTTPHandler.Health)
}
