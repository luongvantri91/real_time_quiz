package models

import "time"

// QuizQuestion represents a quiz question
type QuizQuestion struct {
	ID            int       `json:"id"` // INT SERIAL from DB
	QuizID        string    `json:"quiz_id"`
	Text          string    `json:"text"`
	Options       []string  `json:"options"`
	CorrectAnswer string    `json:"correct_answer"`
	Points        int       `json:"points"`
	OrderNum      int       `json:"order_num"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// AddQuestionRequest is the request body for adding a question
type AddQuestionRequest struct {
	Text          string   `json:"text" validate:"required,min=5,max=1000"`
	Options       []string `json:"options" validate:"required,len=4"` // exactly 4 options
	CorrectAnswer string   `json:"correct_answer" validate:"required,oneof=A B C D"`
	Points        int      `json:"points" validate:"required,min=1,max=100"`
	OrderNum      int      `json:"order_num" validate:"required,min=0"`
}

// DeleteQuestionRequest is the request body for deleting a question
type DeleteQuestionRequest struct {
	QuestionID string `json:"question_id" validate:"required"`
}
