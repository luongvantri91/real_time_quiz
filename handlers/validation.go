package handlers

import (
	"context"
	"fmt"
	"time"
)

// ValidationError represents a validation error
type ValidationError struct {
	Code    string
	Message string
}

func (v *ValidationError) Error() string {
	return v.Message
}

// ValidateQuizExists checks if a quiz exists in DB
func (h *HTTPHandler) ValidateQuizExists(ctx context.Context, quizID string) (bool, error) {
	// Query DB to check if quiz exists
	result, err := h.quizService.GetQuizByID(ctx, quizID)
	if err != nil {
		return false, fmt.Errorf("failed to check quiz existence: %w", err)
	}
	return result != nil, nil
}

// ValidateQuizStatus checks if quiz has the expected status
func (h *HTTPHandler) ValidateQuizStatus(ctx context.Context, quizID string, expectedStatus string) error {
	quiz, err := h.quizService.GetQuizByID(ctx, quizID)
	if err != nil {
		return fmt.Errorf("failed to fetch quiz: %w", err)
	}
	if quiz == nil {
		return &ValidationError{Code: "QUIZ_NOT_FOUND", Message: "Quiz does not exist"}
	}
	if quiz.Status != expectedStatus {
		return &ValidationError{
			Code:    "INVALID_QUIZ_STATUS",
			Message: fmt.Sprintf("Quiz status is '%s', expected '%s'", quiz.Status, expectedStatus),
		}
	}
	return nil
}

// ValidateQuizNotExpired checks if a quiz session has not exceeded duration
func (h *HTTPHandler) ValidateQuizNotExpired(ctx context.Context, quizID string) error {
	quiz, err := h.quizService.GetQuizByID(ctx, quizID)
	if err != nil {
		return fmt.Errorf("failed to fetch quiz: %w", err)
	}
	if quiz == nil {
		return &ValidationError{Code: "QUIZ_NOT_FOUND", Message: "Quiz does not exist"}
	}

	if quiz.Status != "started" {
		return nil // Not started or already ended
	}

	if quiz.StartedAt == nil {
		return &ValidationError{Code: "QUIZ_NOT_STARTED", Message: "Quiz has not been started"}
	}

	duration := time.Duration(quiz.DurationMinutes) * time.Minute
	if time.Since(*quiz.StartedAt) > duration {
		return &ValidationError{
			Code:    "QUIZ_SESSION_EXPIRED",
			Message: fmt.Sprintf("Quiz session expired after %d minutes", quiz.DurationMinutes),
		}
	}

	return nil
}

// ValidateAnswerSubmission checks if answer can be submitted for a quiz
func (h *HTTPHandler) ValidateAnswerSubmission(ctx context.Context, quizID string, participantID string, questionID string) error {
	// Check quiz status
	if err := h.ValidateQuizStatus(ctx, quizID, "started"); err != nil {
		if ve, ok := err.(*ValidationError); ok && ve.Code == "INVALID_QUIZ_STATUS" {
			return &ValidationError{
				Code:    "QUIZ_ENDED",
				Message: "Quiz session has ended. No more answers accepted.",
			}
		}
		return err
	}

	// Check quiz not expired
	if err := h.ValidateQuizNotExpired(ctx, quizID); err != nil {
		if ve, ok := err.(*ValidationError); ok && ve.Code == "QUIZ_SESSION_EXPIRED" {
			return &ValidationError{
				Code:    "QUIZ_EXPIRED",
				Message: "Quiz session expired.",
			}
		}
		return err
	}

	// ✅ Check idempotency: has participant already answered this question?
	answered, err := h.redis.HasAnswered(ctx, quizID, participantID, questionID)
	if err != nil {
		return fmt.Errorf("failed to check answer status: %w", err)
	}
	if answered {
		return &ValidationError{
			Code:    "ALREADY_ANSWERED",
			Message: fmt.Sprintf("You have already answered question %s", questionID),
		}
	}

	return nil
}

// ValidateJoinQuiz checks if participant can join a quiz
func (h *HTTPHandler) ValidateJoinQuiz(ctx context.Context, quizID string, participantID string) error {
	// Check quiz exists and status is 'started'
	if err := h.ValidateQuizStatus(ctx, quizID, "started"); err != nil {
		if ve, ok := err.(*ValidationError); ok && ve.Code == "INVALID_QUIZ_STATUS" {
			return &ValidationError{
				Code:    "QUIZ_NOT_ACTIVE",
				Message: "Quiz is not active. Only 'started' quizzes accept new participants.",
			}
		}
		return err
	}

	// Check quiz not expired
	if err := h.ValidateQuizNotExpired(ctx, quizID); err != nil {
		if ve, ok := err.(*ValidationError); ok && ve.Code == "QUIZ_SESSION_EXPIRED" {
			return &ValidationError{
				Code:    "QUIZ_EXPIRED",
				Message: "Quiz session expired.",
			}
		}
		return err
	}

	// Check participant hasn't joined already
	participants, err := h.redis.GetQuizParticipants(ctx, quizID)
	if err != nil {
		return fmt.Errorf("failed to fetch participants: %w", err)
	}

	for _, p := range participants {
		if p == participantID {
			return &ValidationError{
				Code:    "ALREADY_JOINED",
				Message: "Participant has already joined this quiz",
			}
		}
	}

	return nil
}

// ValidateInitQuiz checks if quiz can be initialized
func (h *HTTPHandler) ValidateInitQuiz(ctx context.Context, quizID string) error {
	// Check quiz exists and status is 'pending'
	if err := h.ValidateQuizStatus(ctx, quizID, "pending"); err != nil {
		if ve, ok := err.(*ValidationError); ok && ve.Code == "INVALID_QUIZ_STATUS" {
			return &ValidationError{
				Code:    "INVALID_STATE",
				Message: "Only 'pending' quizzes can be initialized. Current status must be changed first.",
			}
		}
		return err
	}

	return nil
}

// ValidateEndQuiz checks if quiz can be ended (must exist and be in 'started' status)
func (h *HTTPHandler) ValidateEndQuiz(ctx context.Context, quizID string) error {
	if quizID == "" {
		return &ValidationError{Code: "INVALID_INPUT", Message: "Quiz ID is required"}
	}

	// ✅ Check quiz exists
	quiz, err := h.quizService.GetQuizByID(ctx, quizID)
	if err != nil {
		return fmt.Errorf("failed to fetch quiz: %w", err)
	}
	if quiz == nil {
		return &ValidationError{Code: "QUIZ_NOT_FOUND", Message: "Quiz does not exist"}
	}

	// ✅ Check status is 'started'
	if quiz.Status != "started" {
		return &ValidationError{
			Code:    "INVALID_STATE",
			Message: fmt.Sprintf("Only 'started' quizzes can be ended. Current status: '%s'", quiz.Status),
		}
	}

	return nil
}
