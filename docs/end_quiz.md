# End Quiz API

## Endpoint
```
POST /admin/quizzes/:quizId/end
```

## Description
End quiz (change status from `started` → `ended`). Save final results from Redis to PostgreSQL.

## URL Parameters
- **quizId** (string): Quiz ID

## Request Body
```json
{}
```
(Empty body)

## Response

### Success (200 OK)
```json
{
  "message": "quiz ended successfully",
  "quiz_id": "quiz_a1b2c3d4e5f6",
  "status": "ended",
  "ended_at": "2026-02-02T10:30:00Z",
  "final_results": [
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

**400 Bad Request - Already Ended**
```json
{
  "error": "quiz is not started",
  "code": "QUIZ_NOT_STARTED"
}
```

**500 Internal Server Error**
```json
{
  "error": "failed to save results",
  "code": "DATABASE_ERROR"
}
```

## Example

```bash
curl -X POST http://localhost:8080/admin/quizzes/quiz_a1b2c3d4e5f6/end
```

## Flow (Write-Through Cache)
1. Validate quiz exists and status = `started`
2. Update quiz status → `ended` in PostgreSQL
3. Get final leaderboard from Redis ZSET
4. **Save results to PostgreSQL (synchronous loop)** ✅
   - Loop through leaderboard entries
   - Insert into `results` table
5. Get usernames from PostgreSQL
6. Return final leaderboard

## Database Operations
```sql
-- Update quiz status
UPDATE quizzes 
SET status = 'ended', ended_at = NOW() 
WHERE id = $1

-- Save final results (per participant)
INSERT INTO results (
  id, quiz_id, participant_id, score, rank, completed_at
)
VALUES ($1, $2, $3, $4, $5, NOW())
```

## Redis Operations
```redis
# Get final leaderboard
ZREVRANGE quiz:quiz_a1b2c3d4e5f6:leaderboard 0 -1 WITHSCORES

# Note: Redis data still exists (TTL will expire later)
# No explicit delete of Redis keys
```

## Result Storage
- **Table**: `results`
- **Columns**: 
  - `id` (primary key)
  - `quiz_id` (foreign key)
  - `participant_id` (foreign key)
  - `score` (integer)
  - `rank` (integer, nullable)
  - `completed_at` (timestamp)

## Notes
- **Synchronous result saving**: Loop runs in main thread, not async
- **Context**: Use request context (no goroutine so it's safe)
- **Idempotent**: End quiz multiple times → update results
- Redis data NOT deleted immediately (TTL expires after duration + 24h)
- Can reinitialize quiz after ended (restart)
- Final leaderboard = snapshot at end time
