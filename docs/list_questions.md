# List Questions API

## Endpoint
```
GET /admin/quizzes/:quizId/questions
```

## Description
Get list of all questions for a quiz.

## URL Parameters
- **quizId** (string): Quiz ID

## Response

### Success (200 OK)
```json
{
  "quiz_id": "quiz_a1b2c3d4e5f6",
  "questions": [
    {
      "question_id": "q_1234567890ab",
      "text": "Nước tẩy toilet nên để ở đâu?",
      "options": [
        "Dưới bồn rửa chén",
        "Trong tủ lạnh",
        "Trong tủ có khóa, xa tầm tay trẻ em",
        "Trên bàn ăn"
      ],
      "correct_answer": 2,
      "points": 10
    },
    {
      "question_id": "q_0987654321yx",
      "text": "Bao lâu nên thay bọt biển rửa chén?",
      "options": [
        "1 tháng một lần",
        "6 tháng một lần",
        "1 tuần một lần",
        "1 năm một lần"
      ],
      "correct_answer": 2,
      "points": 10
    }
  ],
  "total": 2
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
  "error": "failed to list questions",
  "code": "DATABASE_ERROR"
}
```

## Example

```bash
curl http://localhost:8080/admin/quizzes/quiz_a1b2c3d4e5f6/questions
```

## Flow
1. Validate quiz exists
2. Query PostgreSQL for all questions
3. Parse options JSON
4. Return question list

## Database Operations
```sql
-- Get all questions for quiz
SELECT 
  id, quiz_id, text, options, correct_answer, points
FROM questions 
WHERE quiz_id = $1
ORDER BY id ASC
```

## Notes
- **Security**: This API exposes `correct_answer` → only for admin use
- User API should not access this endpoint
- Empty list if quiz has no questions
- Options returned as array (parsed from JSON column)
