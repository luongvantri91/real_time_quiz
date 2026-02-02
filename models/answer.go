package models

import "time"

// AnswerRequest represents a participant's answer
type AnswerRequest struct {
	ParticipantID string `json:"participant_id"`
	QuestionID    string `json:"question_id"`
	Answer        string `json:"answer"`
}

// AnswerRecord stored in DB
type AnswerRecord struct {
	QuizID        string    `json:"quiz_id"`
	ParticipantID string    `json:"participant_id"`
	QuestionID    string    `json:"question_id"`
	Answer        string    `json:"answer"`
	IsCorrect     bool      `json:"is_correct"`
	ScoreDelta    int       `json:"score_delta"`
	SubmittedAt   time.Time `json:"submitted_at"`
}

// AnswerResponse represents the response after submitting an answer
type AnswerResponse struct {
	ParticipantID string `json:"participant_id"`
	QuestionID    string `json:"question_id"`
	IsCorrect     bool   `json:"is_correct"`
	ScoreDelta    int    `json:"score_delta"`
	CurrentScore  int    `json:"current_score"`
}
