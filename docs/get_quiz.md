# Get Quiz API

## Endpoint
```
GET /admin/quizzes/:quizId
```

## Description
Get detailed information of a specific quiz.

## URL Parameters
- **quizId** (string): Quiz ID

## Response

### Success (200 OK)
```json
{
  "quiz_id": "quiz_a1b2c3d4e5f6",
  "title": "Quiz về Dọn Dẹp Nhà Cửa",
  "description": "Kiểm tra kiến thức về vệ sinh và dọn dẹp",
  "status": "started",
  "duration_minutes": 30,
  "created_at": "2026-02-02T10:00:00Z",
  "started_at": "2026-02-02T10:05:00Z",
  "ended_at": null,
  "total_questions": 10,
  "total_participants": 15
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
  "error": "failed to get quiz",
  "code": "DATABASE_ERROR"
}
```

## Example

```bash
curl http://localhost:8080/admin/quizzes/quiz_a1b2c3d4e5f6
```

## Flow
1. Validate quiz_id
2. Get quiz from PostgreSQL
3. Count total questions
4. Count total participants
5. Return quiz details

## Database Operations
```sql
-- Get quiz
SELECT * FROM quizzes WHERE id = $1

-- Count questions
SELECT COUNT(*) FROM questions WHERE quiz_id = $1

-- Count participants
SELECT COUNT(*) FROM quiz_participants WHERE quiz_id = $1
```

## Status Values
- **pending**: Quiz created, chưa init
- **started**: Quiz đang chạy
- **ended**: Quiz đã kết thúc

## Notes
- `ended_at` = null if quiz not ended yet
- `started_at` = null if quiz not initialized yet
- `total_questions` and `total_participants` are computed fields
