package handlers

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/you/rt-quiz/models"
)

// CreateQuiz handles POST /admin/quizzes
// Creates a new quiz in pending status
func (h *HTTPHandler) CreateQuiz(c echo.Context) error {
	var req models.CreateQuizRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid request format",
			Code:    "INVALID_REQUEST",
			Details: err.Error(),
		})
	}

	// Validate required fields
	if req.Title == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "title is required",
			Code:  "MISSING_FIELD",
		})
	}

	if req.DurationMinutes < 5 || req.DurationMinutes > 180 {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "duration_minutes must be between 5 and 180",
			Code:  "INVALID_DURATION",
		})
	}

	// Generate quiz ID
	quizID := "quiz_" + uuid.New().String()[:12]

	ctx := c.Request().Context()

	// Create quiz in database
	quiz := &models.Quiz{
		ID:              quizID,
		Title:           req.Title,
		Description:     req.Description,
		Status:          "pending",
		DurationMinutes: req.DurationMinutes,
		CreatedBy:       req.CreatedBy,
		CreatedAt:       time.Now(),
	}

	if err := h.quizService.CreateQuiz(ctx, quiz); err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "failed to create quiz",
			Code:    "DB_ERROR",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, models.CreateQuizResponse{
		ID:        quiz.ID,
		Title:     quiz.Title,
		Status:    quiz.Status,
		CreatedAt: quiz.CreatedAt,
	})
}

// GetQuiz handles GET /admin/quizzes/{quizId}
// Retrieves quiz details
func (h *HTTPHandler) GetQuiz(c echo.Context) error {
	quizID := c.Param("quizId")
	ctx := c.Request().Context()

	quiz, err := h.quizService.GetQuizByID(ctx, quizID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "failed to fetch quiz",
			Code:  "DB_ERROR",
		})
	}

	if quiz == nil {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "quiz not found",
			Code:  "NOT_FOUND",
		})
	}

	return c.JSON(http.StatusOK, quiz)
}

// ListQuizzes handles GET /admin/quizzes
// Lists all quizzes
func (h *HTTPHandler) ListQuizzes(c echo.Context) error {
	ctx := c.Request().Context()

	quizzes, err := h.quizService.ListAllQuizzes(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "failed to list quizzes",
			Code:  "DB_ERROR",
		})
	}

	if quizzes == nil {
		quizzes = []*models.Quiz{} // Empty slice instead of nil
	}

	return c.JSON(http.StatusOK, quizzes)
}

// InitQuiz handles POST /admin/quizzes/{quizId}/init
// Initializes a quiz: pending → started, setup Redis namespace
func (h *HTTPHandler) InitQuiz(c echo.Context) error {
	quizID := c.Param("quizId")
	ctx := c.Request().Context()

	// Validate quiz can be initialized (status = pending)
	if err := h.ValidateInitQuiz(ctx, quizID); err != nil {
		if ve, ok := err.(*ValidationError); ok {
			return c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error: ve.Message,
				Code:  ve.Code,
			})
		}
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
			Code:  "VALIDATION_ERROR",
		})
	}

	// Delegate to service (orchestrates DB update + Redis init)
	quiz, err := h.quizService.InitQuiz(ctx, quizID)
	if err != nil || quiz == nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "failed to initialize quiz",
			Code:  "SERVICE_ERROR",
		})
	}

	return c.JSON(http.StatusOK, models.QuizStatusResponse{
		ID:              quiz.ID,
		Status:          quiz.Status,
		DurationMinutes: quiz.DurationMinutes,
		StartedAt:       quiz.StartedAt,
	})
}

// EndQuiz handles POST /admin/quizzes/{quizId}/end
// Ends a quiz: started → ended, save final results to DB
func (h *HTTPHandler) EndQuiz(c echo.Context) error {
	quizID := c.Param("quizId")
	ctx := c.Request().Context()

	// Validate quiz can be ended (status = started)
	if err := h.ValidateEndQuiz(ctx, quizID); err != nil {
		if ve, ok := err.(*ValidationError); ok {
			return c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error: ve.Message,
				Code:  ve.Code,
			})
		}
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
			Code:  "VALIDATION_ERROR",
		})
	}

	// Delegate to service (orchestrates: get leaderboard, save results, update status, publish event)
	quiz, err := h.quizService.EndQuiz(ctx, quizID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "failed to end quiz",
			Code:  "SERVICE_ERROR",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"quiz_id":  quizID,
		"status":   quiz.Status,
		"ended_at": quiz.EndedAt,
		"message":  "Quiz has been ended. Redis leaderboard kept for 24 hours for user tracking.",
	})
}

// GetQuizStatus handles GET /admin/quizzes/{quizId}/status
// Returns current status and time remaining (if started)
func (h *HTTPHandler) GetQuizStatus(c echo.Context) error {
	quizID := c.Param("quizId")
	ctx := c.Request().Context()

	quiz, err := h.quizService.GetQuizByID(ctx, quizID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "failed to fetch quiz",
			Code:  "DB_ERROR",
		})
	}

	if quiz == nil {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "quiz not found",
			Code:  "NOT_FOUND",
		})
	}

	resp := models.QuizStatusResponse{
		ID:              quiz.ID,
		Status:          quiz.Status,
		DurationMinutes: quiz.DurationMinutes,
		StartedAt:       quiz.StartedAt,
		EndedAt:         quiz.EndedAt,
	}

	// Calculate remaining time if quiz is started
	if quiz.Status == "started" && quiz.StartedAt != nil {
		duration := time.Duration(quiz.DurationMinutes) * time.Minute
		elapsed := time.Since(*quiz.StartedAt)
		remaining := int((duration - elapsed).Seconds())
		if remaining < 0 {
			remaining = 0
		}
		resp.TimeRemaining = remaining
	}

	return c.JSON(http.StatusOK, resp)
}

// AddQuestion handles POST /admin/quizzes/{quizId}/questions
// Adds a new question to a pending quiz
func (h *HTTPHandler) AddQuestion(c echo.Context) error {
	quizID := c.Param("quizId")
	ctx := c.Request().Context()

	var req models.AddQuestionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid request format",
			Code:    "INVALID_REQUEST",
			Details: err.Error(),
		})
	}

	// Validation
	if len(req.Options) != 4 {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "options must have exactly 4 items (A, B, C, D)",
			Code:  "INVALID_OPTIONS",
		})
	}

	if req.Text == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "question text is required",
			Code:  "MISSING_TEXT",
		})
	}

	// ✅ Check quiz exists and status is 'pending'
	if err := h.ValidateQuizStatus(ctx, quizID, "pending"); err != nil {
		if ve, ok := err.(*ValidationError); ok && ve.Code == "INVALID_QUIZ_STATUS" {
			return c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error: "can only add questions to pending quizzes",
				Code:  "INVALID_STATE",
			})
		}
		if ve, ok := err.(*ValidationError); ok && ve.Code == "QUIZ_NOT_FOUND" {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: "quiz not found",
				Code:  "NOT_FOUND",
			})
		}
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
			Code:  "VALIDATION_ERROR",
		})
	}

	question := &models.QuizQuestion{
		// ID auto-generated by SERIAL
		QuizID:        quizID,
		Text:          req.Text,
		Options:       req.Options,
		CorrectAnswer: req.CorrectAnswer,
		Points:        req.Points,
		OrderNum:      req.OrderNum,
	}

	if err := h.quizService.AddQuestion(ctx, question); err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "failed to add question",
			Code:    "DB_ERROR",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, question)
}

// DeleteQuestion handles DELETE /admin/quizzes/{quizId}/questions/{questionId}
// Deletes a question from a pending quiz
func (h *HTTPHandler) DeleteQuestion(c echo.Context) error {
	quizID := c.Param("quizId")
	questionID := c.Param("questionId")
	ctx := c.Request().Context()

	// ✅ Check quiz exists and status is 'pending'
	if err := h.ValidateQuizStatus(ctx, quizID, "pending"); err != nil {
		if ve, ok := err.(*ValidationError); ok && ve.Code == "INVALID_QUIZ_STATUS" {
			return c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error: "can only modify questions in pending quizzes",
				Code:  "INVALID_STATE",
			})
		}
		return c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "quiz not found",
			Code:  "NOT_FOUND",
		})
	}

	if err := h.quizService.DeleteQuestion(ctx, quizID, questionID); err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "failed to delete question",
			Code:    "DB_ERROR",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "question deleted successfully",
	})
}

// ListQuestions handles GET /admin/quizzes/{quizId}/questions
// Lists all questions for a quiz
func (h *HTTPHandler) ListQuestions(c echo.Context) error {
	quizID := c.Param("quizId")
	ctx := c.Request().Context()

	// Check quiz exists
	quiz, err := h.quizService.GetQuizByID(ctx, quizID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "failed to fetch quiz",
			Code:  "DB_ERROR",
		})
	}
	if quiz == nil {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "quiz not found",
			Code:  "NOT_FOUND",
		})
	}

	questions, err := h.quizService.GetQuestions(ctx, quizID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "failed to fetch questions",
			Code:    "DB_ERROR",
			Details: err.Error(),
		})
	}

	if questions == nil {
		questions = []*models.QuizQuestion{}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"quiz_id":   quizID,
		"count":     len(questions),
		"questions": questions,
	})
}
