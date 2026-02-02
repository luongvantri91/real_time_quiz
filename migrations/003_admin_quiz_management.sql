-- Migration: Add quiz results and questions tables for admin management and quiz lifecycle
-- Purpose: Store quiz results snapshot and detailed question data

-- Index for querying started quizzes (for auto-end cron)
CREATE INDEX IF NOT EXISTS idx_quizzes_status ON quizzes(status);
CREATE INDEX IF NOT EXISTS idx_quizzes_started_at ON quizzes(started_at) WHERE status = 'started';

-- Table for quiz_participants (cold-path persistence)
CREATE TABLE IF NOT EXISTS quiz_participants (
    id SERIAL PRIMARY KEY,
    quiz_id TEXT NOT NULL REFERENCES quizzes(id) ON DELETE CASCADE,
    participant_id TEXT NOT NULL,
    username TEXT,
    email TEXT,
    joined_at TIMESTAMP NOT NULL DEFAULT NOW(),
    status TEXT DEFAULT 'active'
);

CREATE INDEX IF NOT EXISTS idx_quiz_participants_quiz_id ON quiz_participants(quiz_id);
CREATE INDEX IF NOT EXISTS idx_quiz_participants_participant_id ON quiz_participants(participant_id);

-- Table for quiz_answers (cold-path persistence)
CREATE TABLE IF NOT EXISTS quiz_answers (
    id SERIAL PRIMARY KEY,
    quiz_id TEXT NOT NULL REFERENCES quizzes(id) ON DELETE CASCADE,
    participant_id TEXT NOT NULL,
    question_id INT NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    answer TEXT,
    is_correct BOOLEAN,
    score_delta INT,
    submitted_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_quiz_answers_quiz_id ON quiz_answers(quiz_id);
CREATE INDEX IF NOT EXISTS idx_quiz_answers_participant ON quiz_answers(participant_id);

