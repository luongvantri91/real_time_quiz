package models

import "time"

// Quiz represents a quiz instance for admin management
type Quiz struct {
	ID              string     `json:"id"`
	Title           string     `json:"title"`
	Description     string     `json:"description,omitempty"`
	Status          string     `json:"status"` // pending, started, ended
	DurationMinutes int        `json:"duration_minutes"`
	CreatedBy       string     `json:"created_by,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	StartedAt       *time.Time `json:"started_at,omitempty"`
	EndedAt         *time.Time `json:"ended_at,omitempty"`
}

// CreateQuizRequest is the request body for creating a new quiz
type CreateQuizRequest struct {
	Title           string `json:"title" validate:"required,min=3,max=200"`
	Description     string `json:"description" validate:"max=1000"`
	DurationMinutes int    `json:"duration_minutes" validate:"required,min=5,max=180"`
	CreatedBy       string `json:"created_by" validate:"required,min=1"`
}

// InitQuizRequest is the request body for initializing quiz in Redis
type InitQuizRequest struct {
	QuizID string `json:"quiz_id" validate:"required"`
}

// EndQuizRequest is the request body for ending a quiz
type EndQuizRequest struct {
	QuizID string `json:"quiz_id" validate:"required"`
}

// CreateQuizResponse is the response after creating a quiz
type CreateQuizResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// QuizStatusResponse represents the current status of a quiz
type QuizStatusResponse struct {
	ID              string     `json:"id"`
	Status          string     `json:"status"`
	DurationMinutes int        `json:"duration_minutes"`
	StartedAt       *time.Time `json:"started_at,omitempty"`
	EndedAt         *time.Time `json:"ended_at,omitempty"`
	TimeRemaining   int        `json:"time_remaining_seconds,omitempty"` // 0 if ended
}

// ErrorResponse represents a standard error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}
