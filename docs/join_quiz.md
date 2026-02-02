# Join Quiz API

## Endpoint
```
POST /quizzes/:quizId/join
```

## Description
User joins quiz. Save to PostgreSQL first, then cache to Redis.

## URL Parameters
- **quizId** (string): Quiz ID

## Request Body
```json
{
  "username": "alice_dev"
}
```

### Fields
- **username** (string, required): Participant name

## Response

### Success (200 OK)
```json
{
  "participant_id": "p_1234567890ab",
  "quiz_id": "quiz_a1b2c3d4e5f6",
  "username": "alice_dev",
  "message": "joined quiz successfully"
}
```

### Error Responses

**400 Bad Request - Missing Username**
```json
{
  "error": "username is required",
  "code": "MISSING_USERNAME"
}
```

**404 Not Found**
```json
{
  "error": "quiz not found",
  "code": "QUIZ_NOT_FOUND"
}
```

**400 Bad Request - Quiz Not Started**
```json
{
  "error": "quiz is not started yet",
  "code": "QUIZ_NOT_STARTED"
}
```

**500 Internal Server Error**
```json
{
  "error": "failed to save participant",
  "code": "DATABASE_ERROR"
}
```

## Example

```bash
curl -X POST http://localhost:8080/quizzes/quiz_a1b2c3d4e5f6/join \
  -H "Content-Type: application/json" \
  -d '{"username": "alice_dev"}'
```

## Flow (Write-Through Cache)
1. Validate quiz exists and status = `started`
2. Validate username
3. Generate participant_id
4. **Save to PostgreSQL** (quiz_participants table) ✅ First
5. **Add to Redis SET** (quiz:{id}:users) ⚠️ Warning if fails
6. **Initialize score** in Redis HASH (quiz:{id}:scores) = 0
7. Return participant_id

## Database Operations
```sql
-- PostgreSQL (Source of Truth)
INSERT INTO quiz_participants (id, quiz_id, username, joined_at)
VALUES ($1, $2, $3, NOW())
```

## Redis Operations
```redis
# Add to users set
SADD quiz:quiz_a1b2c3d4e5f6:users p_1234567890ab

# Initialize score
HSET quiz:quiz_a1b2c3d4e5f6:scores p_1234567890ab 0

# Add to leaderboard
ZADD quiz:quiz_a1b2c3d4e5f6:leaderboard 0 p_1234567890ab
```

## Notes
- Participant ID format: `p_` + 12 random characters
- **DB is source of truth**: PostgreSQL save first, Redis after
- Redis failure → log warning, does not fail request
- Username not unique (can be duplicate)
- Can join when quiz is already started
