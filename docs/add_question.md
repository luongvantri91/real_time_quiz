# Add Question API

## Endpoint
```
POST /admin/quizzes/:quizId/questions
```

## Description
Add question to quiz. Can add when quiz is in `pending` or `started` status.

## URL Parameters
- **quizId** (string): Quiz ID

## Request Body
```json
{
  "text": "Nước tẩy toilet nên để ở đâu?",
  "options": [
    "Dưới bồn rửa chén",
    "Trong tủ lạnh",
    "Trong tủ có khóa, xa tầm tay trẻ em",
    "Trên bàn ăn"
  ],
  "correct_answer": 2,
  "points": 10
}
```

### Fields
- **text** (string, required): Question text
- **options** (array of strings, required): 4 choices
- **correct_answer** (integer, required): Index of correct answer (0-3)
- **points** (integer, required): Points for correct answer (1-100)

## Response

### Success (201 Created)
```json
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
  "points": 10,
  "quiz_id": "quiz_a1b2c3d4e5f6"
}
```

### Error Responses

**400 Bad Request - Invalid Options**
```json
{
  "error": "options must contain exactly 4 choices",
  "code": "INVALID_OPTIONS"
}
```

**400 Bad Request - Invalid Correct Answer**
```json
{
  "error": "correct_answer must be between 0 and 3",
  "code": "INVALID_CORRECT_ANSWER"
}
```

**400 Bad Request - Invalid Points**
```json
{
  "error": "points must be between 1 and 100",
  "code": "INVALID_POINTS"
}
```

**404 Not Found**
```json
{
  "error": "quiz not found",
  "code": "QUIZ_NOT_FOUND"
}
```

## Example

```bash
curl -X POST http://localhost:8080/admin/quizzes/quiz_a1b2c3d4e5f6/questions \
  -H "Content-Type: application/json" \
  -d '{
    "text": "Nước tẩy toilet nên để ở đâu?",
    "options": [
      "Dưới bồn rửa chén",
      "Trong tủ lạnh",
      "Trong tủ có khóa, xa tầm tay trẻ em",
      "Trên bàn ăn"
    ],
    "correct_answer": 2,
    "points": 10
  }'
```

## Flow
1. Validate quiz exists
2. Validate question data (4 options, correct_answer 0-3, points 1-100)
3. Generate question_id
4. Save to PostgreSQL
5. Return question details

## Notes
- Question ID format: `q_` + 12 random characters
- Can add questions while quiz is running (dynamic quiz)
- Correct answer is index (0, 1, 2, or 3)
