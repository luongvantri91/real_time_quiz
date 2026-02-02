package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/you/rt-quiz/models"
)

// Client defines the interface for Redis operations
type Client interface {
	// Participant and score operations
	AddParticipant(ctx context.Context, quizID, participantID string) error
	InitializeParticipantScore(ctx context.Context, quizID, participantID string) error
	UpdateScoreAtomic(ctx context.Context, quizID, participantID string, scoreDelta int) (int, error)
	GetLeaderboard(ctx context.Context, quizID string, limit int) ([]models.LeaderboardEntry, error)
	PublishLeaderboardUpdate(ctx context.Context, quizID string, leaderboard []models.LeaderboardEntry) error
	GetQuizParticipants(ctx context.Context, quizID string) ([]string, error)
	GetScore(ctx context.Context, quizID, participantID string) (int, error)

	// Question cache operations
	CacheQuestions(ctx context.Context, quizID string, questions []*models.QuizQuestion, ttl time.Duration) error
	GetCachedQuestion(ctx context.Context, quizID string, questionID int) (*models.QuizQuestion, error)

	// Answer idempotency (anti-cheat)
	HasAnswered(ctx context.Context, quizID, participantID, questionID string) (bool, error)
	MarkAnswered(ctx context.Context, quizID, participantID, questionID string) error

	// Subscription
	SubscribeToLeaderboardEvents(ctx context.Context, quizID string) *redis.PubSub
	SubscribeToQuizEvents(ctx context.Context, quizID string) (<-chan string, error)

	// Lifecycle & Cleanup
	SetQuizTTL(ctx context.Context, quizID string, ttl time.Duration) error
	Close() error
}

// RedisClient implements the Client interface directly
type RedisClient struct {
	client          *redis.Client
	updateScriptSHA string
}

// NewRedisClient creates a new Redis client from URL
func NewRedisClient(redisURL string) (Client, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("invalid redis url: %w", err)
	}

	client := redis.NewClient(opt)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	log.Println("✓ Redis connected")

	rc := &RedisClient{client: client}

	// Load Lua script for atomic score update
	if err := rc.loadUpdateScript(ctx); err != nil {
		return nil, err
	}

	return rc, nil
}

// NewRedisClientFromConn creates a new Redis client from existing connection
func NewRedisClientFromConn(client *redis.Client) (Client, error) {
	rc := &RedisClient{client: client}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rc.loadUpdateScript(ctx); err != nil {
		return nil, err
	}
	return rc, nil
}

// loadUpdateScript loads the Lua script for atomic score updates
func (rc *RedisClient) loadUpdateScript(ctx context.Context) error {
	// Lua script to atomically update score and leaderboard
	script := `
	local key_users = 'quiz:' .. KEYS[1] .. ':users'
	local key_scores = 'quiz:' .. KEYS[1] .. ':scores'
	local key_leaderboard = 'quiz:' .. KEYS[1] .. ':leaderboard'
	local key_channel = 'quiz:' .. KEYS[1] .. ':events'
	
	local participant_id = ARGV[1]
	local score_delta = tonumber(ARGV[2])
	
	-- 1. Check if participant exists
	if redis.call('SISMEMBER', key_users, participant_id) == 0 then
		return {err = "Participant not found"}
	end
	
	-- 2. Increment score in HASH
	local new_score = redis.call('HINCRBY', key_scores, participant_id, score_delta)
	
	-- 3. Update leaderboard ZSET
	redis.call('ZADD', key_leaderboard, new_score, participant_id)
	
	-- 4. Publish event (for WebSocket broadcast)
	redis.call('PUBLISH', key_channel, 'leaderboard_updated')
	
	return new_score
	`

	sha, err := rc.client.ScriptLoad(ctx, script).Result()
	if err != nil {
		return fmt.Errorf("failed to load lua script: %w", err)
	}

	rc.updateScriptSHA = sha
	return nil
}

// ===== Participant Management =====

func (rc *RedisClient) AddParticipant(ctx context.Context, quizID, participantID string) error {
	key := fmt.Sprintf("quiz:%s:users", quizID)
	return rc.client.SAdd(ctx, key, participantID).Err()
}

func (rc *RedisClient) InitializeParticipantScore(ctx context.Context, quizID, participantID string) error {
	key := fmt.Sprintf("quiz:%s:scores", quizID)
	return rc.client.HSetNX(ctx, key, participantID, 0).Err()
}

func (rc *RedisClient) GetQuizParticipants(ctx context.Context, quizID string) ([]string, error) {
	key := fmt.Sprintf("quiz:%s:users", quizID)
	return rc.client.SMembers(ctx, key).Result()
}

// ===== Scoring & Leaderboard =====

func (rc *RedisClient) UpdateScoreAtomic(ctx context.Context, quizID, participantID string, scoreDelta int) (int, error) {
	result, err := rc.client.EvalSha(ctx, rc.updateScriptSHA, []string{quizID}, participantID, scoreDelta).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to update score: %w", err)
	}

	newScore, ok := result.(int64)
	if !ok {
		return 0, fmt.Errorf("unexpected result type from lua script")
	}

	return int(newScore), nil
}

func (rc *RedisClient) GetScore(ctx context.Context, quizID, participantID string) (int, error) {
	key := fmt.Sprintf("quiz:%s:scores", quizID)
	score, err := rc.client.HGet(ctx, key, participantID).Int()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}
	return score, nil
}

// ===== Answer Idempotency (Anti-Cheat) =====

func (rc *RedisClient) HasAnswered(ctx context.Context, quizID, participantID, questionID string) (bool, error) {
	key := fmt.Sprintf("quiz:%s:answered", quizID)
	member := fmt.Sprintf("%s:%s", participantID, questionID)
	result, err := rc.client.SIsMember(ctx, key, member).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}
	return result, nil
}

func (rc *RedisClient) MarkAnswered(ctx context.Context, quizID, participantID, questionID string) error {
	key := fmt.Sprintf("quiz:%s:answered", quizID)
	member := fmt.Sprintf("%s:%s", participantID, questionID)
	return rc.client.SAdd(ctx, key, member).Err()
}

// ===== Leaderboard =====

func (rc *RedisClient) GetLeaderboard(ctx context.Context, quizID string, limit int) ([]models.LeaderboardEntry, error) {
	key := fmt.Sprintf("quiz:%s:leaderboard", quizID)

	results, err := rc.client.ZRevRangeByScoreWithScores(ctx, key, &redis.ZRangeBy{
		Min:   "-inf",
		Max:   "+inf",
		Count: int64(limit),
	}).Result()
	if err != nil {
		return nil, err
	}

	entries := make([]models.LeaderboardEntry, len(results))
	for i, z := range results {
		entries[i] = models.LeaderboardEntry{
			Rank:          i + 1,
			ParticipantID: z.Member.(string),
			Score:         int(z.Score),
		}
	}

	return entries, nil
}

func (rc *RedisClient) PublishLeaderboardUpdate(ctx context.Context, quizID string, leaderboard []models.LeaderboardEntry) error {
	channel := fmt.Sprintf("quiz:%s:events", quizID)
	return rc.client.Publish(ctx, channel, "leaderboard_updated").Err()
}

// ===== Question Cache Operations =====

// CacheQuestions stores quiz questions in Redis with TTL
func (rc *RedisClient) CacheQuestions(ctx context.Context, quizID string, questions []*models.QuizQuestion, ttl time.Duration) error {
	pipe := rc.client.Pipeline()

	for _, q := range questions {
		key := fmt.Sprintf("quiz:%s:question:%d", quizID, q.ID)
		data := map[string]interface{}{
			"id":             q.ID,
			"quiz_id":        q.QuizID,
			"text":           q.Text,
			"correct_answer": q.CorrectAnswer,
			"points":         q.Points,
			"order_num":      q.OrderNum,
		}

		// Store as hash
		pipe.HSet(ctx, key, data)
		pipe.Expire(ctx, key, ttl)

		// Store options as list
		optionsKey := fmt.Sprintf("quiz:%s:question:%d:options", quizID, q.ID)
		for _, opt := range q.Options {
			pipe.RPush(ctx, optionsKey, opt)
		}
		pipe.Expire(ctx, optionsKey, ttl)
	}

	_, err := pipe.Exec(ctx)
	return err
}

// GetCachedQuestion retrieves a question from Redis cache
func (rc *RedisClient) GetCachedQuestion(ctx context.Context, quizID string, questionID int) (*models.QuizQuestion, error) {
	key := fmt.Sprintf("quiz:%s:question:%d", quizID, questionID)

	// Get question hash
	data, err := rc.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("question not found in cache")
	}

	// Get options list
	optionsKey := fmt.Sprintf("quiz:%s:question:%d:options", quizID, questionID)
	options, err := rc.client.LRange(ctx, optionsKey, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	// Parse data
	q := &models.QuizQuestion{
		QuizID:        data["quiz_id"],
		Text:          data["text"],
		Options:       options,
		CorrectAnswer: data["correct_answer"],
	}

	// Parse int fields
	fmt.Sscanf(data["id"], "%d", &q.ID)
	fmt.Sscanf(data["points"], "%d", &q.Points)
	fmt.Sscanf(data["order_num"], "%d", &q.OrderNum)

	return q, nil
}

// ===== Events & Pub/Sub =====

func (rc *RedisClient) SubscribeToLeaderboardEvents(ctx context.Context, quizID string) *redis.PubSub {
	channel := fmt.Sprintf("quiz:%s:events", quizID)
	return rc.client.Subscribe(ctx, channel)
}

func (rc *RedisClient) SubscribeToQuizEvents(ctx context.Context, quizID string) (<-chan string, error) {
	ch := make(chan string, 100)
	return ch, nil
}

// ===== Lifecycle & Cleanup =====

// SetQuizTTL sets expiration time for all quiz-related keys
func (rc *RedisClient) SetQuizTTL(ctx context.Context, quizID string, ttl time.Duration) error {
	keys := []string{
		fmt.Sprintf("quiz:%s:users", quizID),
		fmt.Sprintf("quiz:%s:scores", quizID),
		fmt.Sprintf("quiz:%s:leaderboard", quizID),
		fmt.Sprintf("quiz:%s:answered", quizID),
	}

	pipe := rc.client.Pipeline()
	for _, key := range keys {
		pipe.Expire(ctx, key, ttl)
	}

	_, err := pipe.Exec(ctx)
	return err
}

func (rc *RedisClient) Close() error {
	return rc.client.Close()
}
