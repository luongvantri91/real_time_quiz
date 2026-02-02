package ws

import (
	"context"

	"github.com/gorilla/websocket"
	"github.com/you/rt-quiz/models"
)

// Client defines the interface for WebSocket operations
type Client interface {
	// Handle WebSocket connection for a quiz
	HandleWebSocket(ctx context.Context, quizID, participantID string, conn interface{}) error

	// Broadcast updates to all connected clients of a quiz
	BroadcastLeaderboardUpdate(ctx context.Context, quizID string, leaderboard []models.LeaderboardEntry) error
	BroadcastQuizEnded(ctx context.Context, quizID string) error
	BroadcastEvent(ctx context.Context, quizID string, eventName string, eventData interface{}) error

	// Unsubscribe a participant from quiz updates
	UnsubscribeParticipant(quizID, participantID string)

	// Register/Unregister client connections
	RegisterClient(quizID string, conn *websocket.Conn)
	UnregisterClient(quizID string, conn *websocket.Conn)

	// Broadcasting internal methods
	BroadcastToClients(quizID string, update models.LeaderboardUpdate)
}

// WebSocketServer implements the Client interface directly
type WebSocketServer struct {
	clients   map[string]map[*websocket.Conn]bool // quizId -> connections
	broadcast chan interface{}
}

// NewWebSocketServer creates a new WebSocket server
func NewWebSocketServer() Client {
	return &WebSocketServer{
		clients:   make(map[string]map[*websocket.Conn]bool),
		broadcast: make(chan interface{}, 100),
	}
}

// RegisterClient registers a WebSocket connection for a quiz
func (ws *WebSocketServer) RegisterClient(quizID string, conn *websocket.Conn) {
	if _, ok := ws.clients[quizID]; !ok {
		ws.clients[quizID] = make(map[*websocket.Conn]bool)
	}
	ws.clients[quizID][conn] = true
}

// UnregisterClient removes a WebSocket connection
func (ws *WebSocketServer) UnregisterClient(quizID string, conn *websocket.Conn) {
	if clients, ok := ws.clients[quizID]; ok {
		delete(clients, conn)
	}
}

// HandleWebSocket handles WebSocket connection (no-op here, handled by handler)
func (ws *WebSocketServer) HandleWebSocket(ctx context.Context, quizID, participantID string, conn interface{}) error {
	return nil
}

// BroadcastLeaderboardUpdate sends leaderboard update to all connected clients
func (ws *WebSocketServer) BroadcastLeaderboardUpdate(ctx context.Context, quizID string, leaderboard []models.LeaderboardEntry) error {
	update := models.LeaderboardUpdate{
		Type:        "leaderboard_update",
		QuizID:      quizID,
		Leaderboard: leaderboard,
	}
	ws.BroadcastToClients(quizID, update)
	return nil
}

// BroadcastQuizEnded sends quiz_ended event to all connected clients
func (ws *WebSocketServer) BroadcastQuizEnded(ctx context.Context, quizID string) error {
	update := models.LeaderboardUpdate{
		Type:   "quiz_ended",
		QuizID: quizID,
	}
	ws.BroadcastToClients(quizID, update)
	return nil
}

// BroadcastEvent sends generic event to all connected clients
func (ws *WebSocketServer) BroadcastEvent(ctx context.Context, quizID string, eventName string, eventData interface{}) error {
	update := models.LeaderboardUpdate{
		Type:   eventName,
		QuizID: quizID,
	}
	ws.BroadcastToClients(quizID, update)
	return nil
}

// BroadcastToClients sends update to all connected clients for a quiz
func (ws *WebSocketServer) BroadcastToClients(quizID string, update models.LeaderboardUpdate) {
	if clients, ok := ws.clients[quizID]; ok {
		for client := range clients {
			_ = client.WriteJSON(update)
		}
	}
}

// UnsubscribeParticipant removes a participant (no-op - handled via connection close)
func (ws *WebSocketServer) UnsubscribeParticipant(quizID, participantID string) {
	// No-op - participant unregistration is handled via connection close
}
