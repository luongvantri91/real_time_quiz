# WebSocket Real-time Updates

## Endpoint
```
GET /quizzes/:quizId/ws
```

## Description
WebSocket connection to receive realtime leaderboard updates. Uses Redis Pub/Sub to broadcast events.

## URL Parameters
- **quizId** (string): Quiz ID

## Query Parameters
- **participant_id** (string, required): Participant ID

## Connection

### Upgrade Request
```bash
# Using wscat (WebSocket CLI tool)
wscat -c "ws://localhost:8080/quizzes/quiz_a1b2c3d4e5f6/ws?participant_id=p_1234567890ab"
```

### Using JavaScript
```javascript
const ws = new WebSocket(
  'ws://localhost:8080/quizzes/quiz_a1b2c3d4e5f6/ws?participant_id=p_1234567890ab'
);

ws.onopen = () => {
  console.log('Connected to quiz leaderboard');
};

ws.onmessage = (event) => {
  const leaderboard = JSON.parse(event.data);
  console.log('Leaderboard update:', leaderboard);
};

ws.onerror = (error) => {
  console.error('WebSocket error:', error);
};

ws.onclose = () => {
  console.log('Disconnected from quiz');
};
```

## Messages

### Initial Message (on connect)
```json
{
  "quiz_id": "quiz_a1b2c3d4e5f6",
  "leaderboard": [
    {
      "participant_id": "p_1234567890ab",
      "username": "alice_dev",
      "score": 50,
      "rank": 1
    },
    {
      "participant_id": "p_0987654321yx",
      "username": "bob_coder",
      "score": 30,
      "rank": 2
    }
  ],
  "total_participants": 2,
  "timestamp": "2026-02-02T10:15:00Z"
}
```

### Update Messages (on answer submission)
```json
{
  "quiz_id": "quiz_a1b2c3d4e5f6",
  "leaderboard": [
    {
      "participant_id": "p_1234567890ab",
      "username": "alice_dev",
      "score": 60,
      "rank": 1
    },
    {
      "participant_id": "p_0987654321yx",
      "username": "bob_coder",
      "score": 30,
      "rank": 2
    }
  ],
  "total_participants": 2,
  "timestamp": "2026-02-02T10:15:05Z"
}
```

## Error Responses

### 400 Bad Request - Missing Participant ID
```
WebSocket upgrade failed: missing participant_id parameter
```

### 404 Not Found - Quiz Not Found
```
WebSocket upgrade failed: quiz not found
```

### 403 Forbidden - Participant Not Joined
```
WebSocket upgrade failed: participant not in quiz
```

## Flow

### Connection Lifecycle
```
1. HTTP GET request with Upgrade header
2. Validate participant_id and quiz_id
3. Upgrade HTTP → WebSocket
4. Send initial leaderboard state
5. Subscribe to Redis Pub/Sub channel
6. Event loop:
   - Redis message → Get leaderboard → Broadcast to client
   - Context cancelled → Close connection
7. Cleanup: Unsubscribe from Redis
```

### Architecture
```
Submit Answer API
    ↓
Update Redis scores
    ↓
PUBLISH to quiz:updates:{quiz_id}
    ↓
Redis Pub/Sub distributes to ALL subscribers
    ↓
WebSocket Handler 1    WebSocket Handler 2    WebSocket Handler 3
    ↓                        ↓                        ↓
Unicast to Client 1    Unicast to Client 2    Unicast to Client 3
```

## Concurrency Model

### Goroutine per Connection
- **1 goroutine** = 1 WebSocket connection
- **Select statement** for event multiplexing
- **Channel-based communication** (Redis → Handler → Client)

### Go Code Pattern
```go
func (h *WSHandler) HandleWebSocket(c echo.Context) error {
    // 1. Upgrade HTTP → WebSocket
    conn, _ := upgrader.Upgrade(...)
    
    // 2. Verify participant
    participantID := c.QueryParam("participant_id")
    
    // 3. Send initial leaderboard
    leaderboard, _ := h.getLeaderboard(quizID)
    conn.WriteJSON(leaderboard)
    
    // 4. Subscribe to Redis Pub/Sub
    pubsub := h.redis.Subscribe(ctx, "quiz:updates:" + quizID)
    ch := pubsub.Channel()
    
    // 5. Event loop
    for {
        select {
        case msg := <-ch:
            // Redis message received
            leaderboard, _ := h.getLeaderboard(quizID)
            conn.WriteJSON(leaderboard)
        case <-ctx.Done():
            // Context cancelled (client disconnected)
            return nil
        }
    }
}
```

## Redis Pub/Sub

### Channel Name
```
quiz:updates:{quiz_id}
```

### Message Format
```json
{
  "participant_id": "p_1234567890ab",
  "score": 60
}
```

### Publish Trigger
- Each answer submission → PUBLISH event
- All subscribers (WebSocket handlers) receive notification
- Each handler fetches updated leaderboard and broadcasts

## Scaling Considerations

### Horizontal Scaling
- **Multiple servers**: Each runs independent WebSocket handlers
- **Redis Pub/Sub**: Distributes events to ALL server instances
- **Stateless**: No shared state between servers

### Connection Limits
- **Per server**: Limited by goroutine count (thousands)
- **Total**: Distributed across all server instances
- **Load balancing**: Sticky sessions NOT required (Redis coordinates)

## Notes
- **Pattern**: Distributed unicast (NOT true broadcast hub)
- **Terminology**: "broadcast" in code = unicast to one client
- **Fault tolerance**: Client disconnect → goroutine exits, auto cleanup
- **Security**: Should validate participant_id with database
- **Performance**: Redis Pub/Sub has latency ~1-5ms
- **Memory**: 1 goroutine ≈ 2KB memory (lightweight)
