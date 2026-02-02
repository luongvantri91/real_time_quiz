# Create Quiz API

## Endpoint
```
POST /admin/quizzes
```

## Description
Create a new quiz with `pending` status. Admin needs to add questions before initializing the quiz.

## Request Body
```json
{
  "title": "Home Cleaning Quiz",
  "description": "Test your knowledge about hygiene and cleaning",
  "duration_minutes": 30
}
```

### Fields
- **title** (string, required): Quiz title
- **description** (string, optional): Quiz description
- **duration_minutes** (integer, required): Quiz duration (5-180 minutes)

## Response

### Success (201 Created)
```json
{
  "quiz_id": "quiz_a1b2c3d4e5f6",
  "title": "Home Cleaning Quiz",
  "description": "Test your knowledge about hygiene and cleaning",
  "status": "pending",
  "duration_minutes": 30,
  "created_at": "2026-02-02T10:00:00Z"
}
```

### Error Responses

**400 Bad Request - Missing Title**
```json
{
  "error": "title is required",
  "code": "MISSING_FIELD"
}
```

**400 Bad Request - Invalid Duration**
```json
{
  "error": "duration_minutes must be between 5 and 180",
  "code": "INVALID_DURATION"
}
```

**500 Internal Server Error**
```json
{
  "error": "failed to create quiz",
  "code": "DATABASE_ERROR",
  "details": "connection refused"
}
```

## Example

```bash
curl -X POST http://localhost:8080/admin/quizzes \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Home Cleaning Quiz",
    "description": "Test your knowledge about hygiene and cleaning",
    "duration_minutes": 30
  }'
```

## Flow
1. Validate request body
2. Generate unique quiz_id
3. Save quiz to PostgreSQL (status: pending)
4. Return quiz details

## Notes
- Quiz ID format: `quiz_` + 12 random characters
- After creating quiz, need to add questions before init
- Quiz can be initialized multiple times (restart quiz)
