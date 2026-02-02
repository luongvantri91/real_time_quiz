package postgres

import (
	"context"
	"database/sql"

	// "fmt"
	"time"

	"github.com/you/rt-quiz/models"
	"github.com/you/rt-quiz/repository"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// Ensure interfaces
var _ repository.ParticipantRepository = (*PostgresRepository)(nil)
var _ repository.AnswerRepository = (*PostgresRepository)(nil)
var _ repository.ResultRepository = (*PostgresRepository)(nil)
var _ repository.QuizRepository = (*PostgresRepository)(nil)

func (p *PostgresRepository) SaveParticipant(ctx context.Context, part *models.Participant) error {
	query := `INSERT INTO quiz_participants (quiz_id, participant_id, username, email, joined_at, status) VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := p.db.ExecContext(ctx, query, part.QuizID, part.ID, part.Username, part.Email, time.Now(), "active")
	return err
}

func (p *PostgresRepository) SaveAnswer(ctx context.Context, a *models.AnswerRecord) error {
	query := `INSERT INTO quiz_answers (quiz_id, participant_id, question_id, answer, is_correct, score_delta, submitted_at) VALUES ($1,$2,$3,$4,$5,$6,$7)`
	_, err := p.db.ExecContext(ctx, query, a.QuizID, a.ParticipantID, a.QuestionID, a.Answer, a.IsCorrect, a.ScoreDelta, time.Now())
	return err
}

func (p *PostgresRepository) SaveResult(ctx context.Context, r *models.QuizResult) error {
	query := `INSERT INTO results (quiz_id, participant_id, score, rank, completed_at) VALUES ($1,$2,$3,$4,$5) ON CONFLICT (quiz_id, participant_id) DO UPDATE SET score = $3, rank = $4, completed_at = $5`
	_, err := p.db.ExecContext(ctx, query, r.QuizID, r.ParticipantID, r.FinalScore, r.Rank, r.CompletedAt)
	return err
}

func (p *PostgresRepository) GetResult(ctx context.Context, quizID, participantID string) (*models.QuizResult, error) {
	query := `SELECT quiz_id, participant_id, score, rank, completed_at FROM results WHERE quiz_id = $1 AND participant_id = $2`
	var r models.QuizResult
	var completedAt sql.NullTime
	err := p.db.QueryRowContext(ctx, query, quizID, participantID).Scan(&r.QuizID, &r.ParticipantID, &r.FinalScore, &r.Rank, &completedAt)
	if err != nil {
		return nil, err
	}
	if completedAt.Valid {
		r.CompletedAt = completedAt.Time
	}
	return &r, nil
}

func (p *PostgresRepository) GetFinalResults(ctx context.Context, quizID string) ([]models.QuizResult, error) {
	query := `SELECT quiz_id, participant_id, score, rank, completed_at FROM results WHERE quiz_id = $1 ORDER BY score DESC, completed_at ASC`
	rows, err := p.db.QueryContext(ctx, query, quizID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []models.QuizResult
	for rows.Next() {
		var r models.QuizResult
		var completedAt sql.NullTime
		if err := rows.Scan(&r.QuizID, &r.ParticipantID, &r.FinalScore, &r.Rank, &completedAt); err != nil {
			return nil, err
		}
		if completedAt.Valid {
			r.CompletedAt = completedAt.Time
		}
		results = append(results, r)
	}
	return results, rows.Err()
}

// Optional: helper to snapshot leaderboard JSON
func (p *PostgresRepository) SaveLeaderboardSnapshot(ctx context.Context, quizID string, jsonData string) error {
	query := `INSERT INTO leaderboard_snapshots (quiz_id, snapshot_at, data) VALUES ($1, $2, $3)`
	_, err := p.db.ExecContext(ctx, query, quizID, time.Now(), jsonData)
	return err
}

// Quiz management methods
func (p *PostgresRepository) CreateQuiz(ctx context.Context, q *models.Quiz) error {
	query := `INSERT INTO quizzes (id, title, description, status, duration_minutes, created_by, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7)`
	_, err := p.db.ExecContext(ctx, query, q.ID, q.Title, q.Description, q.Status, q.DurationMinutes, q.CreatedBy, q.CreatedAt)
	return err
}

func (p *PostgresRepository) GetQuizByID(ctx context.Context, quizID string) (*models.Quiz, error) {
	query := `SELECT id, title, description, status, duration_minutes, created_by, created_at, started_at, ended_at FROM quizzes WHERE id = $1`
	var q models.Quiz
	var startedAt sql.NullTime
	var endedAt sql.NullTime
	err := p.db.QueryRowContext(ctx, query, quizID).Scan(&q.ID, &q.Title, &q.Description, &q.Status, &q.DurationMinutes, &q.CreatedBy, &q.CreatedAt, &startedAt, &endedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if startedAt.Valid {
		q.StartedAt = &startedAt.Time
	}
	if endedAt.Valid {
		q.EndedAt = &endedAt.Time
	}
	return &q, nil
}

func (p *PostgresRepository) ListAllQuizzes(ctx context.Context) ([]*models.Quiz, error) {
	query := `SELECT id, title, description, status, duration_minutes, created_by, created_at, started_at, ended_at FROM quizzes ORDER BY created_at DESC`
	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var quizzes []*models.Quiz
	for rows.Next() {
		var q models.Quiz
		var startedAt sql.NullTime
		var endedAt sql.NullTime
		if err := rows.Scan(&q.ID, &q.Title, &q.Description, &q.Status, &q.DurationMinutes, &q.CreatedBy, &q.CreatedAt, &startedAt, &endedAt); err != nil {
			return nil, err
		}
		if startedAt.Valid {
			q.StartedAt = &startedAt.Time
		}
		if endedAt.Valid {
			q.EndedAt = &endedAt.Time
		}
		quizzes = append(quizzes, &q)
	}
	return quizzes, rows.Err()
}

func (p *PostgresRepository) UpdateQuizStatus(ctx context.Context, quizID string, status string, timestamp *time.Time) error {
	var query string
	switch status {
	case "started":
		query = `UPDATE quizzes SET status = $1, started_at = $2 WHERE id = $3`
		_, err := p.db.ExecContext(ctx, query, status, timestamp, quizID)
		return err
	case "ended":
		query = `UPDATE quizzes SET status = $1, ended_at = $2 WHERE id = $3`
		_, err := p.db.ExecContext(ctx, query, status, timestamp, quizID)
		return err
	default:
		query = `UPDATE quizzes SET status = $1 WHERE id = $2`
		_, err := p.db.ExecContext(ctx, query, status, quizID)
		return err
	}
}

func (p *PostgresRepository) GetAllStartedQuizzes(ctx context.Context) ([]*models.Quiz, error) {
	query := `SELECT id, title, description, status, duration_minutes, created_by, created_at, started_at, ended_at FROM quizzes WHERE status = 'started'`
	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var quizzes []*models.Quiz
	for rows.Next() {
		var q models.Quiz
		var startedAt sql.NullTime
		var endedAt sql.NullTime
		if err := rows.Scan(&q.ID, &q.Title, &q.Description, &q.Status, &q.DurationMinutes, &q.CreatedBy, &q.CreatedAt, &startedAt, &endedAt); err != nil {
			return nil, err
		}
		if startedAt.Valid {
			q.StartedAt = &startedAt.Time
		}
		if endedAt.Valid {
			q.EndedAt = &endedAt.Time
		}
		quizzes = append(quizzes, &q)
	}
	return quizzes, rows.Err()
}

// Question management methods
func (p *PostgresRepository) AddQuestion(ctx context.Context, q *models.QuizQuestion) error {
	query := `INSERT INTO questions (quiz_id, text, options, correct_answer, points, order_num) 
	          VALUES ($1, $2, $3::jsonb, $4, $5, $6) RETURNING id`
	optionsJSON := `["` + q.Options[0] + `","` + q.Options[1] + `","` + q.Options[2] + `","` + q.Options[3] + `"]`
	err := p.db.QueryRowContext(ctx, query, q.QuizID, q.Text, optionsJSON, q.CorrectAnswer, q.Points, q.OrderNum).Scan(&q.ID)
	return err
}

func (p *PostgresRepository) DeleteQuestion(ctx context.Context, quizID, questionID string) error {
	query := `DELETE FROM questions WHERE quiz_id = $1 AND id = $2`
	_, err := p.db.ExecContext(ctx, query, quizID, questionID)
	return err
}

func (p *PostgresRepository) GetQuestion(ctx context.Context, quizID, questionID string) (*models.QuizQuestion, error) {
	query := `SELECT id, text, options, correct_answer, points, order_num FROM questions 
	          WHERE quiz_id = $1 AND id = $2`
	var q models.QuizQuestion
	var optionsJSON string
	err := p.db.QueryRowContext(ctx, query, quizID, questionID).Scan(&q.ID, &q.Text, &optionsJSON, &q.CorrectAnswer, &q.Points, &q.OrderNum)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	q.QuizID = quizID
	// Parse JSON options - simple parsing
	q.Options = parseOptions(optionsJSON)
	return &q, nil
}

func (p *PostgresRepository) GetQuestionsByQuizID(ctx context.Context, quizID string) ([]*models.QuizQuestion, error) {
	query := `SELECT id, text, options, correct_answer, points, order_num FROM questions 
	          WHERE quiz_id = $1 ORDER BY order_num ASC`
	rows, err := p.db.QueryContext(ctx, query, quizID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var questions []*models.QuizQuestion
	for rows.Next() {
		var q models.QuizQuestion
		var optionsJSON string
		if err := rows.Scan(&q.ID, &q.Text, &optionsJSON, &q.CorrectAnswer, &q.Points, &q.OrderNum); err != nil {
			return nil, err
		}
		q.QuizID = quizID
		q.Options = parseOptions(optionsJSON)
		questions = append(questions, &q)
	}
	return questions, rows.Err()
}

func (p *PostgresRepository) CountQuestionsByQuizID(ctx context.Context, quizID string) (int, error) {
	query := `SELECT COUNT(*) FROM questions WHERE quiz_id = $1`
	var count int
	err := p.db.QueryRowContext(ctx, query, quizID).Scan(&count)
	return count, err
}

// Simple helper to parse JSON options
func parseOptions(jsonStr string) []string {
	// Remove brackets and quotes, split by comma
	jsonStr = jsonStr[1 : len(jsonStr)-1]
	var options []string
	for _, opt := range splitOptions(jsonStr) {
		options = append(options, opt)
	}
	return options
}

func splitOptions(s string) []string {
	var result []string
	var current string
	inQuote := false
	for _, c := range s {
		if c == '"' {
			inQuote = !inQuote
		} else if c == ',' && !inQuote {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else if c != ' ' {
			current += string(c)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

// Simple health check
func (p *PostgresRepository) Ping(ctx context.Context) error {
	return p.db.PingContext(ctx)
}
