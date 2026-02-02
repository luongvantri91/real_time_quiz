package models

import "time"

// QuizResult represents final quiz result (for persistence)
type QuizResult struct {
	ID            int       `json:"id"`
	QuizID        string    `json:"quiz_id"`
	ParticipantID string    `json:"participant_id"`
	FinalScore    int       `json:"final_score"` // Maps to DB column 'score'
	Rank          int       `json:"rank"`
	CompletedAt   time.Time `json:"completed_at"`
	CreatedAt     time.Time `json:"created_at"`
}
