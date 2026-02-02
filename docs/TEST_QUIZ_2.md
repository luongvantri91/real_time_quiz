# API Test Workflow - RT-Quiz (Test 2)

## 🎯 Advanced Cleaning Quiz: 3 Users Complete Test

---

## **Phase 1: ADMIN - Create Quiz & Questions**

### **Step 1.1: Create New Quiz**
```powershell
curl -X POST http://localhost:8080/admin/quizzes `
  -H "Content-Type: application/json" `
  -d '{
    "title": "Home Cleaning Advanced Quiz",
    "description": "Test your advanced cleaning knowledge",
    "duration_minutes": 20,
    "created_by": "admin_tri"
  }'
```

**Expected Response:**
```json
{
  "id": "quiz_xyz789",
  "title": "Home Cleaning Advanced Quiz",
  "status": "pending",
  "created_at": "2026-02-02T11:00:00Z"
}
```

**📝 Save quiz_id for use in next steps**


---

### **Step 1.2: Add Question 1**
```powershell
$QUIZ_ID = "quiz_xyz789"  # Replace with quiz_id from Step 1.1

curl -X POST "http://localhost:8080/admin/quizzes/$QUIZ_ID/questions" `
  -H "Content-Type: application/json" `
  -d '{
    "text": "Làm thế nào để loại bỏ vết ố vàng trên gối?",
    "options": ["Giặt với nước nóng", "Ngâm với giấm trắng và baking soda", "Dùng tẩy trắng trực tiếp", "Phơi nắng"],
    "correct_answer": "B",
    "points": 10,
    "order_num": 1
  }'
```

---

### **Step 1.3: Add Question 2**
```powershell
curl -X POST "http://localhost:8080/admin/quizzes/$QUIZ_ID/questions" `
  -H "Content-Type: application/json" `
  -d '{
    "text": "Cách tốt nhất để vệ sinh lò vi sóng?",
    "options": ["Dùng hóa chất mạnh", "Đun nóng nước chanh trong lò", "Lau bằng khăn khô", "Dùng xà phòng rửa chén"],
    "correct_answer": "B",
    "points": 10,
    "order_num": 2
  }'
```

---

### **Step 1.4: Add Question 3**
```powershell
curl -X POST "http://localhost:8080/admin/quizzes/$QUIZ_ID/questions" `
  -H "Content-Type: application/json" `
  -d '{
    "text": "Bao lâu nên giặt màn cửa một lần?",
    "options": ["Mỗi tuần", "Mỗi tháng", "Mỗi 3-6 tháng", "Mỗi năm"],
    "correct_answer": "C",
    "points": 10,
    "order_num": 3
  }'
```

---

### **Step 1.5: Add Question 4**
```powershell
curl -X POST "http://localhost:8080/admin/quizzes/$QUIZ_ID/questions" `
  -H "Content-Type: application/json" `
  -d '{
    "text": "Chất nào KHÔNG nên trộn với nước tẩy (bleach)?",
    "options": ["Nước", "Giấm", "Xà phòng", "Baking soda"],
    "correct_answer": "B",
    "points": 10,
    "order_num": 4
  }'
```

---

### **Step 1.6: Add Question 5**
```powershell
curl -X POST "http://localhost:8080/admin/quizzes/$QUIZ_ID/questions" `
  -H "Content-Type: application/json" `
  -d '{
    "text": "Cách tốt nhất để làm sạch bồn cầu bị ố cứng đầu?",
    "options": ["Dùng cola đổ vào qua đêm", "Chà bằng bàn chải thép", "Dùng nước nóng sôi", "Dùng xà phòng thường"],
    "correct_answer": "A",
    "points": 10,
    "order_num": 5
  }'
```

---

### **Step 1.7: Add Question 6**
```powershell
curl -X POST "http://localhost:8080/admin/quizzes/$QUIZ_ID/questions" `
  -H "Content-Type: application/json" `
  -d '{
    "text": "Làm thế nào để khử mùi thùng rác?",
    "options": ["Rắc bột cà phê hoặc baking soda", "Phun nước hoa", "Rửa bằng nước thường", "Để ngoài nắng"],
    "correct_answer": "A",
    "points": 10,
    "order_num": 6
  }'
```

---

### **Step 1.8: Add Question 7**
```powershell
curl -X POST "http://localhost:8080/admin/quizzes/$QUIZ_ID/questions" `
  -H "Content-Type: application/json" `
  -d '{
    "text": "Vết bẩn gì khó loại bỏ nhất trên quần áo?",
    "options": ["Vết trà", "Vết mực", "Vết dầu mỡ", "Vết rượu vang đỏ"],
    "correct_answer": "C",
    "points": 10,
    "order_num": 7
  }'
```

---

### **Step 1.9: Add Question 8**
```powershell
curl -X POST "http://localhost:8080/admin/quizzes/$QUIZ_ID/questions" `
  -H "Content-Type: application/json" `
  -d '{
    "text": "Bao lâu nên thay giẻ rửa bát?",
    "options": ["Mỗi ngày", "Mỗi tuần", "Mỗi 2 tuần", "Mỗi tháng"],
    "correct_answer": "B",
    "points": 10,
    "order_num": 8
  }'
```

---

### **Step 1.10: Add Question 9**
```powershell
curl -X POST "http://localhost:8080/admin/quizzes/$QUIZ_ID/questions" `
  -H "Content-Type: application/json" `
  -d '{
    "text": "Cách loại bỏ mùi hôi trong máy giặt?",
    "options": ["Giặt không với giấm trắng", "Dùng nước hoa xịt vào", "Để cửa mở thường xuyên", "Cả A và C"],
    "correct_answer": "D",
    "points": 10,
    "order_num": 9
  }'
```

---

### **Step 1.11: Add Question 10**
```powershell
curl -X POST "http://localhost:8080/admin/quizzes/$QUIZ_ID/questions" `
  -H "Content-Type: application/json" `
  -d '{
    "text": "Vật liệu nào KHÔNG nên dùng giấm để làm sạch?",
    "options": ["Gỗ tự nhiên", "Đá granite", "Inox", "Cả A và B"],
    "correct_answer": "D",
    "points": 10,
    "order_num": 10
  }'
```

---

### **Step 1.12: Verify Questions**
```powershell
curl -X GET "http://localhost:8080/admin/quizzes/$QUIZ_ID/questions"
```

**Expected Response:**
```json
{
  "quiz_id": "quiz_xyz789",
  "count": 10,
  "questions": [
    {
      "id": 1,
      "quiz_id": "quiz_xyz789",
      "text": "Làm thế nào để loại bỏ vết ố vàng trên gối?",
      "options": ["Giặt với nước nóng", "Ngâm với giấm trắng và baking soda", "Dùng tẩy trắng trực tiếp", "Phơi nắng"],
      "correct_answer": "B",
      "points": 10,
      "order_num": 1
    }
    // ... 9 more questions
  ]
}
```

---

### **Step 1.13: Initialize Quiz (pending → started)**
```powershell
curl -X POST "http://localhost:8080/admin/quizzes/$QUIZ_ID/init"
```

**Expected Response:**
```json
{
  "id": "quiz_xyz789",
  "status": "started",
  "duration_minutes": 20,
  "started_at": "2026-02-02T11:05:00Z"
}
```

---

## **Phase 2: USER - Join & Answer**

> **📊 Database Operations in Phase 2:**
> - **Join Quiz** → INSERT into `quiz_participants` + Redis SET
> - **Submit Answer** → INSERT into `quiz_answers` + Redis ZADD (leaderboard)
> - **Real-time Updates** → WebSocket broadcasts to all connected clients

### **Step 2.1: User 1 (Alice) Join Quiz**
```powershell
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/join" `
  -H "Content-Type: application/json" `
  -d '{
    "username": "alice_dev",
    "email": "alice@example.com"
  }'
```

**Expected Response:**
```json
{
  "participant_id": "p_alice123",
  "username": "alice_dev",
  "email": "alice@example.com",
  "quiz_id": "quiz_xyz789"
}
```

**📝 Save participant_id as $P1**

**💾 Data saved:**
- **PostgreSQL `quiz_participants`:** INSERT new record
  ```sql
  INSERT INTO quiz_participants (quiz_id, participant_id, username, email, status)
  VALUES ('quiz_xyz789', 'p_alice123', 'alice_dev', 'alice@example.com', 'active');
  ```
- **Redis:** Add participant to SET
  ```
  SADD quiz:quiz_xyz789:participants p_alice123
  ZADD quiz:quiz_xyz789:leaderboard 0 p_alice123
  ```

---

### **Step 2.2: User 2 (Bob) Join Quiz**
```powershell
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/join" `
  -H "Content-Type: application/json" `
  -d '{
    "username": "bob_coder",
    "email": "bob@example.com"
  }'
```

**Expected Response:**
```json
{
  "participant_id": "p_bob456",
  "username": "bob_coder",
  "email": "bob@example.com",
  "quiz_id": "quiz_xyz789"
}
```

**📝 Save participant_id as $P2**

---

### **Step 2.3: User 3 (Charlie) Join Quiz**
```powershell
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/join" `
  -H "Content-Type: application/json" `
  -d '{
    "username": "charlie_hacker",
    "email": "charlie@example.com"
  }'
```

**Expected Response:**
```json
{
  "participant_id": "p_charlie789",
  "username": "charlie_hacker",
  "email": "charlie@example.com",
  "quiz_id": "quiz_xyz789"
}
```

**📝 Save participant_id as $P3**

---

### **Step 2.4: User 1 (Alice) - Submit All Answers (10/10 Correct)**
```powershell
$P1 = "p_alice123"  # Replace with actual participant_id from Step 2.1

# All correct answers (100 points)
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" -H "Content-Type: application/json" -d '{"participant_id": "'"$P1"'", "question_id": "1", "answer": "B"}'
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" -H "Content-Type: application/json" -d '{"participant_id": "'"$P1"'", "question_id": "2", "answer": "B"}'
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" -H "Content-Type: application/json" -d '{"participant_id": "'"$P1"'", "question_id": "3", "answer": "C"}'
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" -H "Content-Type: application/json" -d '{"participant_id": "'"$P1"'", "question_id": "4", "answer": "B"}'
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" -H "Content-Type: application/json" -d '{"participant_id": "'"$P1"'", "question_id": "5", "answer": "A"}'
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" -H "Content-Type: application/json" -d '{"participant_id": "'"$P1"'", "question_id": "6", "answer": "A"}'
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" -H "Content-Type: application/json" -d '{"participant_id": "'"$P1"'", "question_id": "7", "answer": "C"}'
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" -H "Content-Type: application/json" -d '{"participant_id": "'"$P1"'", "question_id": "8", "answer": "B"}'
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" -H "Content-Type: application/json" -d '{"participant_id": "'"$P1"'", "question_id": "9", "answer": "D"}'
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" -H "Content-Type: application/json" -d '{"participant_id": "'"$P1"'", "question_id": "10", "answer": "D"}'
```

**Expected Final Score for Alice:** 100 points (10/10 correct)

---

### **Step 2.5: User 2 (Bob) - Submit All Answers (7/10 Correct)**
```powershell
$P2 = "p_bob456"  # Replace with actual participant_id from Step 2.2

# Q1✓ Q2✓ Q3✗ Q4✓ Q5✗ Q6✗ Q7✓ Q8✓ Q9✓ Q10✗ (70 points)
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" -H "Content-Type: application/json" -d '{"participant_id": "'"$P2"'", "question_id": "1", "answer": "B"}'
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" -H "Content-Type: application/json" -d '{"participant_id": "'"$P2"'", "question_id": "2", "answer": "B"}'
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" -H "Content-Type: application/json" -d '{"participant_id": "'"$P2"'", "question_id": "3", "answer": "B"}'
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" -H "Content-Type: application/json" -d '{"participant_id": "'"$P2"'", "question_id": "4", "answer": "B"}'
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" -H "Content-Type: application/json" -d '{"participant_id": "'"$P2"'", "question_id": "5", "answer": "C"}'
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" -H "Content-Type: application/json" -d '{"participant_id": "'"$P2"'", "question_id": "6", "answer": "D"}'
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" -H "Content-Type: application/json" -d '{"participant_id": "'"$P2"'", "question_id": "7", "answer": "C"}'
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" -H "Content-Type: application/json" -d '{"participant_id": "'"$P2"'", "question_id": "8", "answer": "B"}'
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" -H "Content-Type: application/json" -d '{"participant_id": "'"$P2"'", "question_id": "9", "answer": "D"}'
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" -H "Content-Type: application/json" -d '{"participant_id": "'"$P2"'", "question_id": "10", "answer": "A"}'
```

**Expected Final Score for Bob:** 70 points (7/10 correct - wrong Q3, Q5, Q6, Q10)

---

### **Step 2.6: User 3 (Charlie) - Submit All Answers (5/10 Correct)**
```powershell
$P3 = "p_charlie789"  # Replace with actual participant_id from Step 2.3

# Q1✓ Q2✗ Q3✗ Q4✓ Q5✗ Q6✓ Q7✓ Q8✓ Q9✗ Q10✓ (60 points - corrected)
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" -H "Content-Type: application/json" -d '{"participant_id": "'"$P3"'", "question_id": "1", "answer": "B"}'
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" -H "Content-Type: application/json" -d '{"participant_id": "'"$P3"'", "question_id": "2", "answer": "A"}'
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" -H "Content-Type: application/json" -d '{"participant_id": "'"$P3"'", "question_id": "3", "answer": "A"}'
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" -H "Content-Type: application/json" -d '{"participant_id": "'"$P3"'", "question_id": "4", "answer": "B"}'
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" -H "Content-Type: application/json" -d '{"participant_id": "'"$P3"'", "question_id": "5", "answer": "B"}'
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" -H "Content-Type: application/json" -d '{"participant_id": "'"$P3"'", "question_id": "6", "answer": "A"}'
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" -H "Content-Type: application/json" -d '{"participant_id": "'"$P3"'", "question_id": "7", "answer": "C"}'
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" -H "Content-Type: application/json" -d '{"participant_id": "'"$P3"'", "question_id": "8", "answer": "B"}'
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" -H "Content-Type: application/json" -d '{"participant_id": "'"$P3"'", "question_id": "9", "answer": "C"}'
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" -H "Content-Type: application/json" -d '{"participant_id": "'"$P3"'", "question_id": "10", "answer": "D"}'
```

**Expected Final Score for Charlie:** 60 points (6/10 correct - wrong Q2, Q3, Q5, Q9)

---

### **Step 2.7: Get Current Leaderboard**
```powershell
curl -X GET "http://localhost:8080/quizzes/$QUIZ_ID/leaderboard"
```

**Expected Response:**
```json
[
  {
    "rank": 1,
    "participant_id": "p_alice123",
    "username": "alice_dev",
    "score": 100
  },
  {
    "rank": 2,
    "participant_id": "p_bob456",
    "username": "bob_coder",
    "score": 70
  },
  {
    "rank": 3,
    "participant_id": "p_charlie789",
    "username": "charlie_hacker",
    "score": 60
  }
]
```

---

## **Phase 3: ADMIN - End Quiz & View Results**

> **📊 Database Operations in Phase 3:**
> - **End Quiz** → UPDATE `quizzes` SET status='ended' + INSERT into `results` table
> - **Calculate Final Ranks** → Get from Redis leaderboard, save to `results`
> - **Persistence** → Final results saved permanently to PostgreSQL

### **Step 3.1: End Quiz (started → ended)**
```powershell
curl -X POST "http://localhost:8080/admin/quizzes/$QUIZ_ID/end"
```

**Expected Response:**
```json
{
  "quiz_id": "quiz_xyz789",
  "status": "ended",
  "ended_at": "2026-02-02T11:25:00Z",
  "message": "Quiz has been ended. Redis leaderboard kept for 24 hours for user tracking."
}
```

**💾 Data saved:**
- **PostgreSQL `quizzes`:** UPDATE status
  ```sql
  UPDATE quizzes SET status='ended', ended_at=NOW() WHERE id='quiz_xyz789';
  ```
- **PostgreSQL `results`:** INSERT final rankings
  ```sql
  INSERT INTO results (quiz_id, participant_id, score, rank)
  VALUES 
    ('quiz_xyz789', 'p_alice123', 100, 1),
    ('quiz_xyz789', 'p_bob456', 70, 2),
    ('quiz_xyz789', 'p_charlie789', 60, 3);
  ```

---

### **Step 3.2: Get Quiz Status**
```powershell
curl -X GET "http://localhost:8080/admin/quizzes/$QUIZ_ID/status"
```

**Expected Response:**
```json
{
  "id": "quiz_xyz789",
  "status": "ended",
  "duration_minutes": 20,
  "started_at": "2026-02-02T11:05:00Z",
  "ended_at": "2026-02-02T11:25:00Z",
  "time_remaining_seconds": 0
}
```

---

### **Step 3.3: Get Final Quiz Details**
```powershell
curl -X GET "http://localhost:8080/admin/quizzes/$QUIZ_ID"
```

**Expected Response:**
```json
{
  "id": "quiz_xyz789",
  "title": "Home Cleaning Advanced Quiz",
  "description": "Test your advanced cleaning knowledge",
  "status": "ended",
  "duration_minutes": 20,
  "created_by": "admin_tri",
  "created_at": "2026-02-02T11:00:00Z",
  "started_at": "2026-02-02T11:05:00Z",
  "ended_at": "2026-02-02T11:25:00Z"
}
```

---

## **Phase 4: WebSocket Testing (Real-time Leaderboard)**

> **🔌 WebSocket Endpoint:** `GET /quizzes/{quizId}/ws?participantId={participantId}`

### **Step 4.1: Connect WebSocket**

**Using Postman:**
1. New Request → WebSocket Request
2. URL: `ws://localhost:8080/quizzes/quiz_xyz789/ws?participantId=p_alice123`
3. Click "Connect"
4. View messages in Console tab

### **Step 4.2: Expected WebSocket Messages**

**Message 1 - Initial State (On connect):**
```json
{
  "type": "initial_state",
  "quiz_id": "quiz_xyz789",
  "leaderboard": [
    {
      "rank": 1,
      "participant_id": "p_alice123",
      "score": 100
    },
    {
      "rank": 2,
      "participant_id": "p_bob456",
      "score": 70
    },
    {
      "rank": 3,
      "participant_id": "p_charlie789",
      "score": 60
    }
  ],
  "timestamp": "2026-02-02T11:15:00Z"
}
```

**Message 2 - Leaderboard Update (When user submits answer):**
```json
{
  "type": "leaderboard_update",
  "quiz_id": "quiz_xyz789",
  "leaderboard": [
    {
      "rank": 1,
      "participant_id": "p_alice123",
      "score": 100
    },
    {
      "rank": 2,
      "participant_id": "p_bob456",
      "score": 80
    },
    {
      "rank": 3,
      "participant_id": "p_charlie789",
      "score": 60
    }
  ],
  "timestamp": "2026-02-02T11:16:30Z"
}
```

**Message 3 - Quiz Ended:**
```json
{
  "type": "quiz_ended",
  "quiz_id": "quiz_xyz789",
  "message": "Quiz has ended",
  "timestamp": "2026-02-02T11:25:00Z"
}
```

---

### **Step 4.3: WebSocket Error Scenarios**

**Error 1: participantId missing**
```bash
ws://localhost:8080/quizzes/quiz_xyz789/ws
# Response: 400 Bad Request
# {"error": "participantId query parameter is required"}
```

**Error 2: Participant not joined quiz**
```bash
ws://localhost:8080/quizzes/quiz_xyz789/ws?participantId=p_fake_id
# Response: 400 Bad Request
# {"error": "participant not joined this quiz"}
```

**Error 3: Quiz doesn't exist**
```bash
ws://localhost:8080/quizzes/quiz_invalid/ws?participantId=p_alice123
# WebSocket connects but no initial_state received
```

---

### **Step 4.4: Real-time Testing Flow**

**Terminal 1 - WebSocket Client 1:**
```powershell
wscat -c "ws://localhost:8080/quizzes/quiz_xyz789/ws?participantId=p_alice123"
# Wait for initial_state message
```

**Terminal 2 - WebSocket Client 2:**
```powershell
wscat -c "ws://localhost:8080/quizzes/quiz_xyz789/ws?participantId=p_bob456"
# Wait for initial_state message
```

**Terminal 3 - Submit Answer:**
```powershell
curl -X POST "http://localhost:8080/quizzes/quiz_xyz789/answer" `
  -H "Content-Type: application/json" `
  -d '{"participant_id": "p_bob456", "question_id": "5", "answer": "A"}'
```

**✅ Expected Result:**
- Terminal 1 & 2 both receive `leaderboard_update` message
- Leaderboard automatically updates with new scores
- No need for refresh or polling

**Expected Messages:**
1. **On connect:** `initial_state` with current leaderboard
2. **On answer submit:** `leaderboard_update` with new leaderboard

---

## **📊 Expected Final State**

### **Database - quiz_participants:**
| id | quiz_id | participant_id | username | status |
|----|---------|----------------|----------|--------|
| 1 | quiz_xyz789 | p_alice123 | alice_dev | active |
| 2 | quiz_xyz789 | p_bob456 | bob_coder | active |
| 3 | quiz_xyz789 | p_charlie789 | charlie_hacker | active |

### **Database - results:**
| id | quiz_id | participant_id | score | rank |
|----|---------|----------------|-------|------|
| 1 | quiz_xyz789 | p_alice123 | 100 | 1 |
| 2 | quiz_xyz789 | p_bob456 | 70 | 2 |
| 3 | quiz_xyz789 | p_charlie789 | 60 | 3 |

### **Redis Keys:**
- `quiz:quiz_xyz789:participants` → SET {p_alice123, p_bob456, p_charlie789}
- `quiz:quiz_xyz789:leaderboard` → ZSET {p_alice123:100, p_bob456:70, p_charlie789:60}
- **TTL:** ~86400 seconds (24 hours)

---


## **🎯 Success Criteria**

✅ **Phase 1:** Quiz created successfully with 10 questions, status = "started"  
✅ **Phase 2:** 3 users join, submit 30 total answers (Alice: 100pts, Bob: 70pts, Charlie: 60pts)  
✅ **Leaderboard:** Correct order (Alice 1st, Bob 2nd, Charlie 3rd)  
✅ **Phase 3:** Quiz ended successfully, results saved to DB  
✅ **WebSocket:** Real-time leaderboard updates working  
✅ **Database:** All data persisted correctly (3 participants, 30 answers, 3 results)

**🚀 Ready to Test!** Run each step in order to test the complete advanced workflow.
