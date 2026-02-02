# Initialize Quiz API

## Endpoint
```
POST /admin/quizzes/:quizId/init
```

## Description
Start the quiz (change status from `pending` → `started`). Cache questions to Redis and set TTL. Can restart ended quiz.

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
  "message": "quiz initialized successfully",
  "quiz_id": "quiz_a1b2c3d4e5f6",
  "status": "started",
  "start_time": "2026-02-02T10:00:00Z",
  "end_time": "2026-02-02T10:30:00Z",
  "redis_ttl_hours": 24,
  "cached_questions": 10
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

**500 Internal Server Error - No Questions**
```json
{
  "error": "quiz has no questions",
  "code": "NO_QUESTIONS"
}
```

**500 Internal Server Error - Redis Failure**
```json
{
  "error": "failed to cache questions",
  "code": "REDIS_ERROR"
}
```

## Example

```bash
curl -X POST http://localhost:8080/admin/quizzes/quiz_a1b2c3d4e5f6/init
```

## Flow
1. Validate quiz exists
2. Load questions from PostgreSQL
3. Update quiz status → `started`
4. Calculate TTL = `duration_minutes + 24 hours`
5. Cache questions to Redis
6. Set TTL for Redis keys:
   - `quiz:{id}:users`
   - `quiz:{id}:scores`
   - `quiz:{id}:leaderboard`
   - `quiz:{id}:answered`
   - `quiz:{id}:question:{qid}`
7. Return success

## Redis Keys Created
```
quiz:quiz_a1b2c3d4e5f6:users          → SET (empty)
quiz:quiz_a1b2c3d4e5f6:scores         → HASH (empty)
quiz:quiz_a1b2c3d4e5f6:leaderboard    → ZSET (empty)
quiz:quiz_a1b2c3d4e5f6:answered       → SET (empty)
quiz:quiz_a1b2c3d4e5f6:question:q_xxx → HASH (question data)
```

## TTL Strategy
- **TTL Duration**: `quiz_duration + 24 hours`
- **Auto-cleanup**: Redis automatically removes expired keys
- **Memory Policy**: allkeys-lru (evict least recently used)
- **Purpose**: Allow users to review results sau khi quiz kết thúc

## Notes
- Can restart ended quiz (reinitialize)
- TTL ensures data doesn't exist forever in Redis
- Questions cached to avoid DB reads for each request
- Multi-quiz support: Init quiz 2 does not delete quiz 1 data
