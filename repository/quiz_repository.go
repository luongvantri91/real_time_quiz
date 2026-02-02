package repository

import (
	"context"
	"time"

	"github.com/you/rt-quiz/models"
)

// QuizRepository defines DB operations for quizzes metadata and lifecycle
type QuizRepository interface {
	CreateQuiz(ctx context.Context, q *models.Quiz) error
	GetQuizByID(ctx context.Context, quizID string) (*models.Quiz, error)
	ListAllQuizzes(ctx context.Context) ([]*models.Quiz, error)
	UpdateQuizStatus(ctx context.Context, quizID string, status string, timestamp *time.Time) error
	GetAllStartedQuizzes(ctx context.Context) ([]*models.Quiz, error)

	// Question management
	AddQuestion(ctx context.Context, q *models.QuizQuestion) error
	DeleteQuestion(ctx context.Context, quizID, questionID string) error
	GetQuestion(ctx context.Context, quizID, questionID string) (*models.QuizQuestion, error)
	GetQuestionsByQuizID(ctx context.Context, quizID string) ([]*models.QuizQuestion, error)
	CountQuestionsByQuizID(ctx context.Context, quizID string) (int, error)
}
