# Submit Answer API

## Endpoint
```
POST /quizzes/:quizId/answer
```

## Description
Submit answer. Check anti-cheat, validate correctness, save to DB, update Redis, broadcast realtime.

## URL Parameters
- **quizId** (string): Quiz ID

## Request Body
```json
{
  "participant_id": "p_1234567890ab",
  "question_id": "q_0987654321yx",
  "answer_index": 2
}
```

### Fields
- **participant_id** (string, required): Participant ID
- **question_id** (string, required): Question ID
- **answer_index** (integer, required): Answer index (0-3)

## Response

### Success (200 OK)
```json
{
  "correct": true,
  "points_earned": 10,
  "total_score": 50,
  "message": "answer submitted successfully"
}
```

### Error Responses

**400 Bad Request - Invalid Answer Index**
```json
{
  "error": "answer_index must be between 0 and 3",
  "code": "INVALID_ANSWER_INDEX"
}
```

**404 Not Found - Question Not Found**
```json
{
  "error": "question not found",
  "code": "QUESTION_NOT_FOUND"
}
```

**409 Conflict - Already Answered**
```json
{
  "error": "question already answered",
  "code": "ALREADY_ANSWERED"
}
```

**500 Internal Server Error**
```json
{
  "error": "failed to save answer",
  "code": "DATABASE_ERROR"
}
```

## Example

```bash
curl -X POST http://localhost:8080/quizzes/quiz_a1b2c3d4e5f6/answer \
  -H "Content-Type: application/json" \
  -d '{
    "participant_id": "p_1234567890ab",
    "question_id": "q_0987654321yx",
    "answer_index": 2
  }'
```

## Flow (Write-Through Cache)
1. **Anti-cheat check**: Check Redis SET `quiz:{id}:answered`
   - If exists `p_xxx:q_yyy` → return 409 Conflict
2. **Get question**: From Redis cache (fallback to DB)
3. **Validate answer**: Compare với correct_answer
4. **Calculate points**: 
   - Correct → points_earned = question.points
   - Wrong → points_earned = 0
5. **Save to PostgreSQL** (quiz_answers table) ✅ **SYNCHRONOUS**
6. **Update Redis HASH** (quiz:{id}:scores) ⚠️ Warning if fails
7. **Update Redis ZSET** (quiz:{id}:leaderboard)
8. **Mark as answered** in Redis SET (quiz:{id}:answered)
9. **Publish to Redis Pub/Sub** (quiz:updates:{id})
10. Return result

## Database Operations
```sql
-- PostgreSQL (Source of Truth)
INSERT INTO quiz_answers (
  id, quiz_id, participant_id, question_id, 
  answer_index, is_correct, points_earned, answered_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
```

## Redis Operations
```redis
# Check if already answered (anti-cheat)
SISMEMBER quiz:quiz_a1b2c3d4e5f6:answered p_xxx:q_yyy

# Mark as answered
SADD quiz:quiz_a1b2c3d4e5f6:answered p_xxx:q_yyy

# Update score
HINCRBY quiz:quiz_a1b2c3d4e5f6:scores p_1234567890ab 10

# Update leaderboard
ZINCRBY quiz:quiz_a1b2c3d4e5f6:leaderboard 10 p_1234567890ab

# Publish realtime update
PUBLISH quiz:updates:quiz_a1b2c3d4e5f6 '{"participant_id":"p_xxx","score":50}'
```

## Anti-Cheat Mechanism
- Redis SET `quiz:{id}:answered` tracks `participant_id:question_id`
- Prevents duplicate submissions
- Checked BEFORE database save

## Realtime Broadcast
- Publish to Redis Pub/Sub channel: `quiz:updates:{quiz_id}`
- All WebSocket handlers receive update
- Each handler unicasts to its connected client
- Result: All clients see leaderboard update instantly

## Notes
- **Synchronous DB save**: No goroutine, ensures data saved
- **Context**: Use `context.Background()` with 5s timeout
- **Idempotent**: Submit duplicate question → 409 Conflict
- Answer index must be 0-3 (4 options)
- Points earned = 0 if wrong
