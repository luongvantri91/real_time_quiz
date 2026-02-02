# Get Leaderboard API

## Endpoint
```
GET /quizzes/:quizId/leaderboard
```

## Description
Get leaderboard from Redis (realtime scores). Sorted by score descending.

## URL Parameters
- **quizId** (string): Quiz ID

## Response

### Success (200 OK)
```json
{
  "quiz_id": "quiz_a1b2c3d4e5f6",
  "leaderboard": [
    {
      "participant_id": "p_1234567890ab",
      "username": "alice_dev",
      "score": 100,
      "rank": 1
    },
    {
      "participant_id": "p_0987654321yx",
      "username": "bob_coder",
      "score": 70,
      "rank": 2
    },
    {
      "participant_id": "p_abcdef123456",
      "username": "charlie_hacker",
      "score": 50,
      "rank": 3
    }
  ],
  "total_participants": 3
}
```

### Error Responses

**404 Not Found**
```json
{
  "error": "quiz not found",
  "code": "QUIZ_NOT_FOUND"
}
```

**500 Internal Server Error**
```json
{
  "error": "failed to get leaderboard",
  "code": "REDIS_ERROR"
}
```

## Example

```bash
curl http://localhost:8080/quizzes/quiz_a1b2c3d4e5f6/leaderboard
```

## Flow
1. Validate quiz exists
2. Get sorted scores from Redis ZSET
3. Get usernames from PostgreSQL for each participant_id
4. Assign ranks (1-indexed)
5. Return leaderboard

## Redis Operations
```redis
# Get top scores (descending order)
ZREVRANGE quiz:quiz_a1b2c3d4e5f6:leaderboard 0 -1 WITHSCORES

# Example output:
# 1) "p_1234567890ab"
# 2) "100"
# 3) "p_0987654321yx"
# 4) "70"
# 5) "p_abcdef123456"
# 6) "50"
```

## Database Operations
```sql
-- Get usernames for participants
SELECT id, username 
FROM quiz_participants 
WHERE id IN ($1, $2, $3)
```

## Ranking Logic
- **Descending order**: Highest score = Rank 1
- **Tie handling**: Same score → same rank (hoặc timestamp-based)
- **0 to -1**: Get all participants (không limit)

## Notes
- Data source: **Redis ZSET** (realtime, cached scores)
- Username lookup: **PostgreSQL** (persistent data)
- Use case: Display current standings during quiz
- WebSocket subscribers receive updates automatically
- Leaderboard can have 0 participants (quiz just initialized)
