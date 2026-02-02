package repository

import (
	"context"
	"github.com/you/rt-quiz/models"
)

// ResultRepository defines DB operations for final results
type ResultRepository interface {
	SaveResult(ctx context.Context, r *models.QuizResult) error
	GetResult(ctx context.Context, quizID, participantID string) (*models.QuizResult, error)
	GetFinalResults(ctx context.Context, quizID string) ([]models.QuizResult, error)
}
