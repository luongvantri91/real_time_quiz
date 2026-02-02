package repository

import (
	"context"
	"github.com/you/rt-quiz/models"
)

// AnswerRepository defines DB operations for answers
type AnswerRepository interface {
	SaveAnswer(ctx context.Context, a *models.AnswerRecord) error
}
