package models

import "time"

// Participant represents a quiz participant
type Participant struct {
	ID                    string    `json:"participant_id"`
	Username              string    `json:"username"`
	Email                 string    `json:"email"`
	QuizID                string    `json:"quiz_id"`
	JoinedAt              time.Time `json:"joined_at"`
	Status                string    `json:"status,omitempty"` // active, inactive
	ParticipantExternalID string    `json:"-"`
}
