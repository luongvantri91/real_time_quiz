# Get Quiz Status API

## Endpoint
```
GET /admin/quizzes/:quizId/status
```

## Description
Get quiz status and remaining time (if running).

## URL Parameters
- **quizId** (string): Quiz ID

## Response

### Success (200 OK) - Quiz Started
```json
{
  "quiz_id": "quiz_a1b2c3d4e5f6",
  "status": "started",
  "started_at": "2026-02-02T10:00:00Z",
  "duration_minutes": 30,
  "end_time": "2026-02-02T10:30:00Z",
  "time_remaining_seconds": 900,
  "is_active": true
}
```

### Success (200 OK) - Quiz Pending
```json
{
  "quiz_id": "quiz_a1b2c3d4e5f6",
  "status": "pending",
  "started_at": null,
  "duration_minutes": 30,
  "end_time": null,
  "time_remaining_seconds": null,
  "is_active": false
}
```

### Success (200 OK) - Quiz Ended
```json
{
  "quiz_id": "quiz_a1b2c3d4e5f6",
  "status": "ended",
  "started_at": "2026-02-02T10:00:00Z",
  "ended_at": "2026-02-02T10:30:00Z",
  "duration_minutes": 30,
  "time_remaining_seconds": 0,
  "is_active": false
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

## Example

```bash
curl http://localhost:8080/admin/quizzes/quiz_a1b2c3d4e5f6/status
```

## Flow
1. Get quiz from PostgreSQL
2. Calculate time remaining:
   - If status = `started`: `end_time - now`
   - Else: null
3. Determine `is_active`:
   - true if `started` and time remaining > 0
   - false otherwise
4. Return status info

## Time Calculation
```go
if quiz.Status == "started" && quiz.StartedAt != nil {
    endTime := quiz.StartedAt.Add(time.Duration(quiz.DurationMinutes) * time.Minute)
    timeRemaining := time.Until(endTime).Seconds()
    
    if timeRemaining < 0 {
        timeRemaining = 0
    }
}
```

## Notes
- `time_remaining_seconds` can be < 0 if quiz expired but not ended
- `is_active` = true only when quiz started and time remaining > 0
- Use case: Client polling to check quiz state
- Recommend: Combine with WebSocket for realtime updates
