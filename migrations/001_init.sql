-- Migration 001: Create tables for RT-Quiz

CREATE TABLE IF NOT EXISTS quizzes (
    id TEXT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending, started, ended
    duration_minutes INT NOT NULL DEFAULT 30,
    created_by VARCHAR(255),
    started_at TIMESTAMP,
    ended_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS questions (
    id SERIAL PRIMARY KEY,
    quiz_id TEXT NOT NULL REFERENCES quizzes(id) ON DELETE CASCADE,
    text TEXT NOT NULL,
    options JSONB NOT NULL,  -- ["A", "B", "C", "D"]
    correct_answer VARCHAR(1) NOT NULL, -- A, B, C, D
    points INT NOT NULL DEFAULT 10,
    order_num INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_questions_quiz_id ON questions(quiz_id);

CREATE TABLE IF NOT EXISTS results (
    id SERIAL PRIMARY KEY,
    quiz_id TEXT NOT NULL REFERENCES quizzes(id) ON DELETE CASCADE,
    participant_id VARCHAR(255) NOT NULL,
    score INT NOT NULL DEFAULT 0,
    rank INT,
    completed_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (quiz_id, participant_id)
);

CREATE INDEX idx_results_quiz_id ON results(quiz_id);
CREATE INDEX idx_results_participant_id ON results(participant_id);
