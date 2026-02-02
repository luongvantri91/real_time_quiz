package services

import (
	"context"
	"fmt"
	"time"

	"github.com/you/rt-quiz/client/redis"
	"github.com/you/rt-quiz/models"
	"github.com/you/rt-quiz/repository"
)

// ParticipantService handles participant-related business logic
type ParticipantService struct {
	participantRepo repository.ParticipantRepository
	answerRepo      repository.AnswerRepository
	redisClient     redis.Client
}

// NewParticipantService constructs a ParticipantService
func NewParticipantService(pr repository.ParticipantRepository, ar repository.AnswerRepository, rc redis.Client) *ParticipantService {
	return &ParticipantService{
		participantRepo: pr,
		answerRepo:      ar,
		redisClient:     rc,
	}
}

// JoinQuiz adds a participant to a quiz - DB first (source of truth), then Redis (cache)
func (s *ParticipantService) JoinQuiz(ctx context.Context, quizID, participantID, username, email string) error {
	// 1️⃣ Save to DB FIRST (source of truth, persistent storage)
	participant := &models.Participant{
		ID:       participantID,
		QuizID:   quizID,
		Username: username,
		Email:    email,
	}

	if err := s.participantRepo.SaveParticipant(ctx, participant); err != nil {
		return fmt.Errorf("failed to save participant to DB: %w", err)
	}

	// 2️⃣ Add to Redis (cache layer) - only if DB succeeded
	if err := s.redisClient.AddParticipant(ctx, quizID, participantID); err != nil {
		// Log warning but don't fail - Redis is cache, DB already has data
		fmt.Printf("WARNING: failed to add participant to Redis cache: %v\n", err)
	}

	// 3️⃣ Initialize score in Redis
	if err := s.redisClient.InitializeParticipantScore(ctx, quizID, participantID); err != nil {
		fmt.Printf("WARNING: failed to initialize score in Redis: %v\n", err)
	}

	return nil
}

// SubmitAnswer handles answer submission - DB first (persistent), then Redis (cache + broadcast)
func (s *ParticipantService) SubmitAnswer(ctx context.Context, quizID, participantID, questionID, answer string, isCorrect bool, scoreDelta int) (int, error) {
	// 1️⃣ Save to DB FIRST (source of truth) - SYNCHRONOUS to ensure data consistency
	record := &models.AnswerRecord{
		QuizID:        quizID,
		ParticipantID: participantID,
		QuestionID:    questionID,
		Answer:        answer,
		IsCorrect:     isCorrect,
		ScoreDelta:    scoreDelta,
		SubmittedAt:   time.Now(),
	}

	if err := s.answerRepo.SaveAnswer(ctx, record); err != nil {
		return 0, fmt.Errorf("failed to save answer to DB: %w", err)
	}

	// 2️⃣ Mark as answered in Redis (anti-cheat) - only if DB succeeded
	if err := s.redisClient.MarkAnswered(ctx, quizID, participantID, questionID); err != nil {
		// Log warning but continue - DB already has record
		fmt.Printf("WARNING: failed to mark answered in Redis: %v\n", err)
	}

	// 3️⃣ Update score in Redis (cache layer)
	newScore, err := s.redisClient.UpdateScoreAtomic(ctx, quizID, participantID, scoreDelta)
	if err != nil {
		// Log warning but don't fail - DB already saved
		fmt.Printf("WARNING: failed to update score in Redis: %v\n", err)
		// Return scoreDelta as newScore since Redis update failed
		newScore = scoreDelta
	}

	// 4️⃣ Publish leaderboard update (real-time feature, non-critical)
	leaderboard, _ := s.redisClient.GetLeaderboard(ctx, quizID, 100)
	_ = s.redisClient.PublishLeaderboardUpdate(ctx, quizID, leaderboard)

	return newScore, nil
}

// GetLeaderboard retrieves the current leaderboard from Redis
func (s *ParticipantService) GetLeaderboard(ctx context.Context, quizID string, limit int) ([]models.LeaderboardEntry, error) {
	return s.redisClient.GetLeaderboard(ctx, quizID, limit)
}
