# List Quizzes API

## Endpoint
```
GET /admin/quizzes
```

## Description
Get list of all quizzes in the system.

## Query Parameters
- **status** (string, optional): Filter by status (`pending`, `started`, `ended`)
- **limit** (integer, optional): Number of quizzes to return (default: 50)
- **offset** (integer, optional): Offset for pagination (default: 0)

## Response

### Success (200 OK)
```json
{
  "quizzes": [
    {
      "quiz_id": "quiz_a1b2c3d4e5f6",
      "title": "Quiz về Dọn Dẹp Nhà Cửa",
      "status": "started",
      "duration_minutes": 30,
      "created_at": "2026-02-02T10:00:00Z",
      "total_questions": 10,
      "total_participants": 15
    },
    {
      "quiz_id": "quiz_xyz123456789",
      "title": "Quiz về An Toàn Thực Phẩm",
      "status": "ended",
      "duration_minutes": 45,
      "created_at": "2026-02-01T14:00:00Z",
      "total_questions": 20,
      "total_participants": 50
    }
  ],
  "total": 2,
  "limit": 50,
  "offset": 0
}
```

### Error Responses

**500 Internal Server Error**
```json
{
  "error": "failed to list quizzes",
  "code": "DATABASE_ERROR"
}
```

## Examples

### Get all quizzes
```bash
curl http://localhost:8080/admin/quizzes
```

### Filter by status
```bash
curl "http://localhost:8080/admin/quizzes?status=started"
```

### Pagination
```bash
curl "http://localhost:8080/admin/quizzes?limit=10&offset=20"
```

## Flow
1. Parse query parameters (status, limit, offset)
2. Query PostgreSQL with filters
3. Count total questions per quiz
4. Count total participants per quiz
5. Return quiz list

## Database Operations
```sql
-- Get quizzes with optional status filter
SELECT * FROM quizzes 
WHERE ($1::text IS NULL OR status = $1)
ORDER BY created_at DESC 
LIMIT $2 OFFSET $3

-- Count total
SELECT COUNT(*) FROM quizzes 
WHERE ($1::text IS NULL OR status = $1)
```

## Notes
- Default sort: `created_at DESC` (newest first)
- `total_questions` and `total_participants` computed per quiz
- Empty list if no quizzes exist
