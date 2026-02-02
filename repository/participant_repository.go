package repository

import (
	"context"
	"github.com/you/rt-quiz/models"
)

// ParticipantRepository defines DB operations for participants
type ParticipantRepository interface {
	SaveParticipant(ctx context.Context, p *models.Participant) error
}
