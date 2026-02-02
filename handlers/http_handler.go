package handlers

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/you/rt-quiz/client/redis"
	"github.com/you/rt-quiz/models"
	"github.com/you/rt-quiz/services"
)

// HTTPHandler handles HTTP requests
type HTTPHandler struct {
	redis              redis.Client
	quizService        *services.QuizService
	participantService *services.ParticipantService
	quizzes            map[string]*QuizData // In-memory quiz definitions
}

// QuizData represents quiz questions
type QuizData struct {
	ID        string
	Questions map[string]models.QuizQuestion
}

func NewHTTPHandler(redisClient redis.Client) *HTTPHandler {
	handler := &HTTPHandler{
		redis:   redisClient,
		quizzes: make(map[string]*QuizData),
	}

	// Load sample quiz for demo
	handler.initializeSampleQuiz()
	return handler
}

func NewHTTPHandlerWithServices(qs *services.QuizService, ps *services.ParticipantService, redisClient redis.Client) *HTTPHandler {
	handler := &HTTPHandler{
		redis:              redisClient,
		quizService:        qs,
		participantService: ps,
		quizzes:            make(map[string]*QuizData),
	}
	handler.initializeSampleQuiz()
	return handler
}

// initializeSampleQuiz initializes a sample quiz
func (h *HTTPHandler) initializeSampleQuiz() {
	quiz1 := &QuizData{
		ID: "quiz1",
		Questions: map[string]models.QuizQuestion{
			"q1": {
				ID:            1, // INT from DB
				Text:          "What is 2+2?",
				Options:       []string{"A: 1", "B: 4", "C: 5", "D: 3"},
				CorrectAnswer: "B",
				Points:        10,
			},
			"q2": {
				ID:            2, // INT from DB
				Text:          "What is the capital of France?",
				Options:       []string{"A: London", "B: Paris", "C: Berlin", "D: Madrid"},
				CorrectAnswer: "B",
				Points:        10,
			},
		},
	}
	h.quizzes["quiz1"] = quiz1
}

// JoinQuiz handles POST /quizzes/{quizId}/join
func (h *HTTPHandler) JoinQuiz(c echo.Context) error {
	quizID := c.Param("quizId")

	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if req.Username == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "username required"})
	}

	// Generate participant ID
	participantID := "p_" + uuid.New().String()[:8]

	ctx := c.Request().Context()

	// Delegate to service
	err := h.participantService.JoinQuiz(ctx, quizID, participantID, req.Username, req.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to join quiz"})
	}

	return c.JSON(http.StatusOK, models.Participant{
		ID:       participantID,
		Username: req.Username,
		Email:    req.Email,
		QuizID:   quizID,
	})
}

// SubmitAnswer handles POST /quizzes/{quizId}/answer
func (h *HTTPHandler) SubmitAnswer(c echo.Context) error {
	quizID := c.Param("quizId")

	var req models.AnswerRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	// Validation: ensure answer allowed (quiz started, not expired, participant joined, not answered yet)
	if err := h.ValidateAnswerSubmission(c.Request().Context(), quizID, req.ParticipantID, req.QuestionID); err != nil {
		if ve, ok := err.(*ValidationError); ok {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": ve.Message})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "validation error"})
	}

	ctx := c.Request().Context()

	// Try to get question from Redis cache first (fast path)
	questionID := 0
	fmt.Sscanf(req.QuestionID, "%d", &questionID)

	question, err := h.redis.GetCachedQuestion(ctx, quizID, questionID)
	if err != nil {
		// Fallback to database if Redis cache miss
		questions, dbErr := h.quizService.GetQuestions(ctx, quizID)
		if dbErr != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get questions"})
		}
		if len(questions) == 0 {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "quiz has no questions"})
		}

		// Find the specific question
		for _, q := range questions {
			if q.QuizID == quizID && q.ID == questionID {
				question = q
				break
			}
		}
		if question == nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "question not found"})
		}
	}

	// Check if answer is correct
	isCorrect := req.Answer == question.CorrectAnswer
	scoreDelta := 0
	if isCorrect {
		scoreDelta = question.Points
	}

	// ctx := c.Request().Context()

	// Delegate to service (orchestrates: update score, mark answered, broadcast leaderboard, save to DB async)
	newScore, err := h.participantService.SubmitAnswer(ctx, quizID, req.ParticipantID, req.QuestionID, req.Answer, isCorrect, scoreDelta)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to submit answer"})
	}

	return c.JSON(http.StatusOK, models.AnswerResponse{
		ParticipantID: req.ParticipantID,
		QuestionID:    req.QuestionID,
		IsCorrect:     isCorrect,
		ScoreDelta:    scoreDelta,
		CurrentScore:  newScore,
	})
}

// GetLeaderboard handles GET /quizzes/{quizId}/leaderboard
func (h *HTTPHandler) GetLeaderboard(c echo.Context) error {
	quizID := c.Param("quizId")

	leaderboard, err := h.participantService.GetLeaderboard(c.Request().Context(), quizID, 100)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get leaderboard"})
	}

	return c.JSON(http.StatusOK, leaderboard)
}

// Health handles GET /health
func (h *HTTPHandler) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}
