# Delete Question API

## Endpoint
```
DELETE /admin/quizzes/:quizId/questions/:questionId
```

## Description
Delete question from quiz. Can delete when quiz is in `pending` or `started` status.

## URL Parameters
- **quizId** (string): Quiz ID
- **questionId** (string): Question ID to delete

## Response

### Success (200 OK)
```json
{
  "message": "question deleted successfully",
  "question_id": "q_1234567890ab",
  "quiz_id": "quiz_a1b2c3d4e5f6"
}
```

### Error Responses

**404 Not Found - Quiz Not Found**
```json
{
  "error": "quiz not found",
  "code": "QUIZ_NOT_FOUND"
}
```

**404 Not Found - Question Not Found**
```json
{
  "error": "question not found",
  "code": "QUESTION_NOT_FOUND"
}
```

**500 Internal Server Error**
```json
{
  "error": "failed to delete question",
  "code": "DATABASE_ERROR"
}
```

## Example

```bash
curl -X DELETE http://localhost:8080/admin/quizzes/quiz_a1b2c3d4e5f6/questions/q_1234567890ab
```

## Flow
1. Validate quiz exists
2. Validate question exists và belongs to quiz
3. Delete from PostgreSQL
4. (Optional) Remove from Redis cache if quiz started
5. Return success

## Database Operations
```sql
-- Delete question
DELETE FROM questions 
WHERE id = $1 AND quiz_id = $2
```

## Redis Operations (if quiz started)
```redis
# Remove cached question
DEL quiz:quiz_a1b2c3d4e5f6:question:q_1234567890ab
```

## Notes
- Can delete question even when quiz is running
- Answers already submitted for this question still exist in DB
- Deleting question does not affect calculated scores
- Recommend: Do not delete question when quiz is active (confusing for users)
