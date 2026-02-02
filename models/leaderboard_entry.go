package models

// LeaderboardEntry represents a single entry in the leaderboard
type LeaderboardEntry struct {
	Rank          int    `json:"rank"`
	ParticipantID string `json:"participant_id"`
	Score         int    `json:"score"`
}

// LeaderboardUpdate is sent via WebSocket when leaderboard changes
type LeaderboardUpdate struct {
	Type        string             `json:"type"`
	QuizID      string             `json:"quiz_id"`
	Leaderboard []LeaderboardEntry `json:"leaderboard"`
}
