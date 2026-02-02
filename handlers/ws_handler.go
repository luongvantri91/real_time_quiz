package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/you/rt-quiz/client/redis"
	"github.com/you/rt-quiz/client/ws"
	"github.com/you/rt-quiz/models"
)

// WSHandler handles WebSocket connections
type WSHandler struct {
	wsClient ws.Client
	redis    redis.Client
}

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for demo
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewWSHandler(wsClient ws.Client, redisClient redis.Client) *WSHandler {
	return &WSHandler{
		wsClient: wsClient,
		redis:    redisClient,
	}
}

// WSMessage represents a message sent via WebSocket
type WSMessage struct {
	Type        string                    `json:"type"` // "initial_state", "leaderboard_update"
	QuizID      string                    `json:"quiz_id"`
	Leaderboard []models.LeaderboardEntry `json:"leaderboard"`
	Timestamp   time.Time                 `json:"timestamp"`
}

// HandleWebSocket handles WebSocket connections at GET /quizzes/{quizId}/ws?participantId={id}
// Sends real-time leaderboard updates via Pub/Sub
func (h *WSHandler) HandleWebSocket(c echo.Context) error {
	quizID := c.Param("quizId")
	participantID := c.QueryParam("participantId")

	// ✅ Verify participantId is provided
	if participantID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "participantId query parameter is required",
		})
	}

	// ✅ Verify participant has joined this quiz (check in Redis)
	participants, err := h.redis.GetQuizParticipants(c.Request().Context(), quizID)
	if err != nil {
		log.Printf("Error fetching participants: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to verify participant",
		})
	}

	found := false
	for _, p := range participants {
		if p == participantID {
			found = true
			break
		}
	}
	if !found {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "participant not joined this quiz",
		})
	}

	// ✅ Upgrade connection
	conn, err := wsUpgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		return err
	}
	defer conn.Close()

	// ✅ Register client with service
	h.wsClient.RegisterClient(quizID, conn)
	defer h.wsClient.UnregisterClient(quizID, conn)

	// ✅ Send initial state (current leaderboard)
	leaderboard, err := h.redis.GetLeaderboard(c.Request().Context(), quizID, 100)
	if err != nil {
		log.Printf("Error getting initial leaderboard: %v", err)
		return nil
	}

	initialMsg := WSMessage{
		Type:        "initial_state",
		QuizID:      quizID,
		Leaderboard: leaderboard,
		Timestamp:   time.Now(),
	}
	if err := conn.WriteJSON(initialMsg); err != nil {
		log.Printf("Error sending initial state: %v", err)
		return nil
	}

	// ✅ Subscribe to Redis Pub/Sub for leaderboard updates
	pubsub := h.redis.SubscribeToLeaderboardEvents(c.Request().Context(), quizID)
	defer pubsub.Close()

	// ✅ Listen for Redis Pub/Sub events and push to this client
	// (Redis handles broadcasting to all subscribed handlers)
	ch := pubsub.Channel()
	for {
		select {
		case msg := <-ch:
			if msg == nil {
				return nil
			}

			// ✅ Fetch updated leaderboard from Redis
			leaderboard, err := h.redis.GetLeaderboard(c.Request().Context(), quizID, 100)
			if err != nil {
				log.Printf("Error getting leaderboard: %v", err)
				continue
			}

			// ✅ Send leaderboard update to client
			update := WSMessage{
				Type:        "leaderboard_update",
				QuizID:      quizID,
				Leaderboard: leaderboard,
				Timestamp:   time.Now(),
			}

			if err := conn.WriteJSON(update); err != nil {
				log.Printf("Error writing to WS: %v", err)
				return nil
			}

		case <-c.Request().Context().Done():
			return nil
		}
	}
}
