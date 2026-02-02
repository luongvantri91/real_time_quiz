# API Test Workflow - RT-Quiz

## 🎯 Complete Workflow: Admin → User → End Quiz

---

## **Phase 1: ADMIN - Create Quiz & Questions**

### **Step 1.1: Create New Quiz**
```powershell
curl -X POST http://localhost:8080/admin/quizzes `
  -H "Content-Type: application/json" `
  -d '{
    "title": "Cleaning App Mastery Quiz",
    "description": "Test your knowledge about cleaning best practices",
    "duration_minutes": 30,
    "created_by": "admin_tri"
  }'
```

**Expected Response:**
```json
{
  "id": "quiz_abc123",
  "title": "Cleaning App Mastery Quiz",
  "status": "pending",
  "created_at": "2026-02-02T10:00:00Z"
}
```

**📝 Save quiz_id for use in next steps**

---

### **Step 1.2: Add Question 1**
```powershell
$QUIZ_ID = "quiz_abc123"  # Replace with quiz_id from Step 1.1

curl -X POST "http://localhost:8080/admin/quizzes/$QUIZ_ID/questions" `
  -H "Content-Type: application/json" `
  -d '{
    "text": "Mua loại bộ vệ sinh nào để làm sạch kính cửa hiệu quả nhất?",
    "options": ["Chải lông lợn", "Lau kính chuyên dụng", "Giẻ thô", "Cọ nhựa"],
    "correct_answer": "B",
    "points": 10,
    "order_num": 1
  }'
```

**Expected Response:**
```json
{
  "id": 1,
  "quiz_id": "quiz_abc123",
  "text": "Mua loại bộ vệ sinh nào...",
  "options": ["Chải lông lợn", "Lau kính chuyên dụng", "Giẻ thô", "Cọ nhựa"],
  "correct_answer": "B",
  "points": 10,
  "order_num": 1
}
```

---

### **Step 1.3: Add Question 2**
```powershell
curl -X POST "http://localhost:8080/admin/quizzes/$QUIZ_ID/questions" `
  -H "Content-Type: application/json" `
  -d '{
    "text": "Bao lâu nên thay dụng cụ lau chùi?",
    "options": ["Mỗi năm một lần", "Mỗi tháng một lần", "Mỗi 2-3 tháng", "Không cần thay"],
    "correct_answer": "C",
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
    "text": "Chất nào an toàn nhất để vệ sinh sàn nhà?",
    "options": ["Axit", "Nước và xà phòng", "Cồn", "Bleach"],
    "correct_answer": "B",
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
    "text": "Thứ tự đúng khi vệ sinh phòng là gì?",
    "options": ["Lau sàn → Hút bụi → Lau bàn", "Hút bụi → Lau bàn → Lau sàn", "Lau bàn → Lau sàn → Hút bụi", "Lau sàn → Lau bàn → Hút bụi"],
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
    "text": "Nhiệt độ nước tốt nhất để lau sàn gỗ là bao nhiêu?",
    "options": ["Nước lạnh", "Nước ấm (30-40°C)", "Nước nóng (60-70°C)", "Nước sôi"],
    "correct_answer": "B",
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
    "text": "Làm thế nào để loại bỏ mùi hôi trong tủ lạnh?",
    "options": ["Dùng nước hoa", "Đặt bột baking soda", "Phun xịt khử mùi", "Để cửa mở"],
    "correct_answer": "B",
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
    "text": "Bao lâu nên vệ sinh máy hút bụi một lần?",
    "options": ["Sau mỗi lần sử dụng", "Mỗi tuần một lần", "Mỗi tháng một lần", "Mỗi 6 tháng một lần"],
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
    "text": "Cách tốt nhất để làm sạch vết bẩn trên thảm là gì?",
    "options": ["Chà mạnh bằng bàn chải", "Thấm nhẹ nhàng từ ngoài vào trong", "Dùng nước nóng đổ lên", "Để khô tự nhiên"],
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
    "text": "Nên vệ sinh nhà bếp vào thời điểm nào là hiệu quả nhất?",
    "options": ["Ngay sau khi nấu ăn", "Buổi sáng sớm", "Cuối tuần", "Khi có thời gian rảnh"],
    "correct_answer": "A",
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
    "text": "Dụng cụ nào KHÔNG nên dùng cho bề mặt inox?",
    "options": ["Khăn mềm", "Miếng bọt biển", "Bàn chải thép", "Khăn microfiber"],
    "correct_answer": "C",
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
  "quiz_id": "quiz_abc123",
  "count": 10,
  "questions": [
    {
      "id": 1,
      "quiz_id": "quiz_abc123",
      "text": "Mua loại bộ vệ sinh nào...",
      "options": ["..."],
      "correct_answer": "B",
      "points": 10,
      "order_num": 1
    },
    // ... 9 questions more
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
  "id": "quiz_abc123",
  "status": "started",
  "duration_minutes": 30,
  "started_at": "2026-02-02T10:05:00Z"
}
```

---

## **Phase 2: USER - Join & Answer**

> **📊 Database Operations in Phase 2:**
> - **Join Quiz** → INSERT into `quiz_participants` + Redis SET
> - **Submit Answer** → INSERT into `quiz_answers` + Redis ZADD (leaderboard)
> - **Anti-cheat** → Redis SET checks `quiz:quiz_id:answered`

### **Step 2.1: User 1 Join Quiz**
```powershell
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/join" `
  -H "Content-Type: application/json" `
  -d '{
    "username": "user_nguyen_van_a",
    "email": "nguyenvana@example.com"
  }'
```

**Expected Response:**
```json
{
  "participant_id": "p_12ab34cd",
  "username": "user_nguyen_van_a",
  "email": "nguyenvana@example.com",
  "quiz_id": "quiz_abc123"
}
```

**📝 Save participant_id**

**💾 Data saved:**
- **PostgreSQL `quiz_participants`:** INSERT new record
  ```sql
  INSERT INTO quiz_participants (quiz_id, participant_id, username, email, status)
  VALUES ('quiz_abc123', 'p_12ab34cd', 'user_nguyen_van_a', 'nguyenvana@example.com', 'active');
  ```
- **Redis:** Add participant to SET
  ```
  SADD quiz:quiz_abc123:participants p_12ab34cd
  ZADD quiz:quiz_abc123:leaderboard 0 p_12ab34cd
  ```

---

### **Step 2.2: User 2 Join Quiz**
```powershell
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/join" `
  -H "Content-Type: application/json" `
  -d '{
    "username": "user_tran_thi_b",
    "email": "tranthib@example.com"
  }'
```

**Expected Response:**
```json
{
  "participant_id": "p_56ef78gh",
  "username": "user_tran_thi_b",
  "email": "tranthib@example.com",
  "quiz_id": "quiz_abc123"
}
```

---

### **Step 2.3: User 3 Join Quiz**
```powershell
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/join" `
  -H "Content-Type: application/json" `
  -d '{
    "username": "user_le_van_c",
    "email": "levanc@example.com"
  }'
```

---

### **Step 2.4: User 1 Submit Answer Q1 (Correct)**
```powershell
$PARTICIPANT_1 = "p_12ab34cd"  # From Step 2.1

curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" `
  -H "Content-Type: application/json" `
  -d '{
    "participant_id": "'"$PARTICIPANT_1"'",
    "question_id": "1",
    "answer": "B"
  }'
```

**Expected Response:**
```json
{
  "participant_id": "p_12ab34cd",
  "question_id": "1",
  "is_correct": true,
  "score_delta": 10,
  "current_score": 10
}
```

**💾 Data saved:**
- **PostgreSQL `quiz_answers`:** INSERT answer record
  ```sql
  INSERT INTO quiz_answers (quiz_id, participant_id, question_id, answer, is_correct, score_delta)
  VALUES ('quiz_abc123', 'p_12ab34cd', 1, 'B', true, 10);
  ```
- **Redis Leaderboard:** Update score
  ```
  ZADD quiz:quiz_abc123:leaderboard 10 p_12ab34cd
  SADD quiz:quiz_abc123:answered p_12ab34cd:1
  ```
- **WebSocket:** Broadcast leaderboard update to all clients

---

### **Step 2.5: User 1 Submit Answer Q2 (Wrong)**
```powershell
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" `
  -H "Content-Type: application/json" `
  -d '{
    "participant_id": "'"$PARTICIPANT_1"'",
    "question_id": "2",
    "answer": "A"
  }'
```

**Expected Response:**
```json
{
  "participant_id": "p_12ab34cd",
  "question_id": "2",
  "is_correct": false,
  "score_delta": 0,
  "current_score": 10
}
```

---

### **Step 2.6: User 1 Submit Answer Q3 (Correct)**
```powershell
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" `
  -H "Content-Type: application/json" `
  -d '{
    "participant_id": "'"$PARTICIPANT_1"'",
    "question_id": "3",
    "answer": "B"
  }'
```

**Expected Response:**
```json
{
  "participant_id": "p_12ab34cd",
  "question_id": "3",
  "is_correct": true,
  "score_delta": 10,
  "current_score": 20
}
```

---

### **Step 2.7: User 2 Submit Answers (All Correct - 10/10)**
```powershell
$PARTICIPANT_2 = "p_56ef78gh"

# Q1-Q10 - All Correct (100 points)
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" -H "Content-Type: application/json" -d '{"participant_id": "'"$PARTICIPANT_2"'", "question_id": "1", "answer": "B"}'
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" -H "Content-Type: application/json" -d '{"participant_id": "'"$PARTICIPANT_2"'", "question_id": "2", "answer": "C"}'
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" -H "Content-Type: application/json" -d '{"participant_id": "'"$PARTICIPANT_2"'", "question_id": "3", "answer": "B"}'
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" -H "Content-Type: application/json" -d '{"participant_id": "'"$PARTICIPANT_2"'", "question_id": "4", "answer": "B"}'
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" -H "Content-Type: application/json" -d '{"participant_id": "'"$PARTICIPANT_2"'", "question_id": "5", "answer": "B"}'
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" -H "Content-Type: application/json" -d '{"participant_id": "'"$PARTICIPANT_2"'", "question_id": "6", "answer": "B"}'
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" -H "Content-Type: application/json" -d '{"participant_id": "'"$PARTICIPANT_2"'", "question_id": "7", "answer": "C"}'
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" -H "Content-Type: application/json" -d '{"participant_id": "'"$PARTICIPANT_2"'", "question_id": "8", "answer": "B"}'
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" -H "Content-Type: application/json" -d '{"participant_id": "'"$PARTICIPANT_2"'", "question_id": "9", "answer": "A"}'
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" -H "Content-Type: application/json" -d '{"participant_id": "'"$PARTICIPANT_2"'", "question_id": "10", "answer": "C"}'
```

**Expected Final Score for User 2:** 100 points

---

### **Step 2.8: User 3 Submit Answers (Mixed)**
```powershell
$PARTICIPANT_3 = "p_90ij12kl"

# Q1 - Correct
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" `
  -H "Content-Type: application/json" `
  -d '{
    "participant_id": "'"$PARTICIPANT_3"'",
    "question_id": "1",
    "answer": "B"
  }'

# Q2 - Wrong
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" `
  -H "Content-Type: application/json" ` - 5/10)**
```powershell
$PARTICIPANT_3 = "p_90ij12kl"

# Q1 - Correct, Q2 - Wrong, Q3 - Wrong, Q4 - Correct, Q5 - Wrong
# Q6 - Correct, Q7 - Correct, Q8 - Wrong, Q9 - Correct, Q10 - Wrong
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" -H "Content-Type: application/json" -d '{"participant_id": "'"$PARTICIPANT_3"'", "question_id": "1", "answer": "B"}'
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" -H "Content-Type: application/json" -d '{"participant_id": "'"$PARTICIPANT_3"'", "question_id": "2", "answer": "D"}'
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" -H "Content-Type: application/json" -d '{"participant_id": "'"$PARTICIPANT_3"'", "question_id": "3", "answer": "A"}'
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" -H "Content-Type: application/json" -d '{"participant_id": "'"$PARTICIPANT_3"'", "question_id": "4", "answer": "B"}'
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" -H "Content-Type: application/json" -d '{"participant_id": "'"$PARTICIPANT_3"'", "question_id": "5", "answer": "A"}'
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" -H "Content-Type: application/json" -d '{"participant_id": "'"$PARTICIPANT_3"'", "question_id": "6", "answer": "B"}'
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" -H "Content-Type: application/json" -d '{"participant_id": "'"$PARTICIPANT_3"'", "question_id": "7", "answer": "C"}'
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" -H "Content-Type: application/json" -d '{"participant_id": "'"$PARTICIPANT_3"'", "question_id": "8", "answer": "A"}'
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" -H "Content-Type: application/json" -d '{"participant_id": "'"$PARTICIPANT_3"'", "question_id": "9", "answer": "A"}'
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" -H "Content-Type: application/json" -d '{"participant_id": "'"$PARTICIPANT_3"'", "question_id": "10", "answer": "A"}'
```

**Expected Final Score for User 3:** 50 points (5 correct × 10 points)
[
  {
    "rank": 2,
    "participant_id": "p_12ab34cd",
    "score": 20
  },
  {
    "rank": 3,
    "participant_id": "p_90ij12kl",
    "score": 10
  }
]
```

---

### **Step 2.10: Test Anti-Cheat (Duplicate Answer)**
```powershell
# User 1 tries to answer Q1 again
curl -X POST "http://localhost:8080/quizzes/$QUIZ_ID/answer" `
  -H "Content-Type: application/json" `
  -d '{
    "participant_id": "'"$PARTICIPANT_1"'",
    "question_id": "1",
    "answer": "B"
  }'
```

**Expected Response:**
```json
{
  "error": "participant has already answered this question"
}
```

**💾 Anti-cheat mechanism:**
- **Redis Check:** `SISMEMBER quiz:quiz_abc123:answered p_12ab34cd:1` → returns 1 (already exists)
- **Result:** API rejects duplicate submission, does not INSERT into database

---

## **Phase 3: ADMIN - End Quiz & View Results**

> **📊 Database Operations in Phase 3:**
> - **End Quiz** → UPDATE `quizzes` SET status='ended' + INSERT into `results` table
> - **Calculate Final Ranks** → Get from Redis leaderboard, save to `results`
> - **Snapshot** → INSERT into `quiz_results_snapshot` (if needed)

**Expected Re100
  },
  {
    "rank": 2,
    "participant_id": "p_90ij12kl",
    "score": 50
  },
  {
    "rank": 3,
    "participant_id": "p_12ab34cd",
    "score": 2
### **Step 3.1: End Quiz (started → ended)**
```powershell
curl -X POST "http://localhost:8080/admin/quizzes/$QUIZ_ID/end"
```

**Expected Response:**
```json
{
  "

**💾 Data saved:**
- **PostgreSQL `quizzes`:** UPDATE status
  ```sql
  UPDATE quizzes SET status='ended', ended_at=NOW() WHERE id='quiz_abc123';
  ```
- **PostgreSQL `results`:** INSERT final rankings
  ```sql
  -- Get from Redis leaderboard and save permanently
  INSERT INTO results (quiz_id, participant_id, score, rank)
  VALUES 
    ('quiz_abc123', 'p_56ef78gh', 100, 1),
    ('quiz_abc123', 'p_90ij12kl', 50, 2),
    ('quiz_abc123', 'p_12ab34cd', 20, 3);
  ```
- **Redis:** Keep data for 24h for user trackingquiz_id": "quiz_abc123",
  "status": "ended",
  "ended_at": "2026-02-02T10:30:00Z",
  "message": "Quiz has been ended. Redis leaderboard kept for 24 hours for user tracking."
}
```

---

### **Step 3.2: Get Quiz Status**
```powershell
curl -X GET "http://localhost:8080/admin/quizzes/$QUIZ_ID/status"
```

**Expected Response:**
```json
{
  "id": "quiz_abc123",
  "status": "ended",
  "duration_minutes": 30,
  "started_at": "2026-02-02T10:05:00Z",
  "ended_at": "2026-02-02T10:30:00Z",
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
  "id": "quiz_abc123",
  "title": "Cleaning App Mastery Quiz",
  "description": "Test your knowledge...",
  "status": "ended",
  "duration_minutes": 30,
  "created_by": "admin_tri",
  "created_at": "2026-02-02T10:00:00Z",
  "started_at": "2026-02-02T10:05:00Z",
  "ended_at": "2026-02-02T10:30:00Z"
}
```

---

### **Step 3.4: List All Quizzes**
```powershell
curl -X GET "http://localhost:8080/admin/quizzes"
```

**Expected Response:**
```json
[
  {
    "id": "quiz_abc123",
    "title": "Cleaning App Mastery Quiz",
    "status": "ended",
    "created_at": "2026-02-02T10:00:00Z"
  }
  // ... other quizzes
]
```

---

## **Phase 4: WebSocket Testing (Real-time Leaderboard)**

> **🔌 WebSocket Endpoint:**  
> `GET /quizzes/{quizId}/ws?participantId={participantId}`
>
> **Authentication:** Query parameter `participantId` (required)  
> **Verification:** Participant must have joined quiz before connecting

### **Step 4.1: Connect WebSocket**

**Using Browser Console (JavaScript):**
```javascript
// Use browser console or WebSocket client
const ws = new WebSocket('ws://localhost:8080/quizzes/quiz_abc123/ws?participantId=p_12ab34cd');

ws.onopen = () => {
  console.log('WebSocket connected');
};

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log('Received:', data);
  
  // Initial state or leaderboard update
  if (data.type === 'initial_state') {
    console.log('Initial leaderboard:', data.leaderboard);
  } else if (data.type === 'leaderboard_update') {
    console.log('Updated leaderboard:', data.leaderboard);
  }
};

ws.onerror = (error) => {
  console.error('WebSocket error:', error);
};
```

**Using Postman:**
1. New Request → WebSocket Request
2. URL: `ws://localhost:8080/quizzes/quiz_abc123/ws?participantId=p_12ab34cd`
3. Click "Connect"
4. Xem messages trong Console tab

---

### **Step 4.2: Expected WebSocket Messages**

**Message 1 - Initial State (Ngay khi connect):**
```json
{
  "type": "initial_state",
  "quiz_id": "quiz_abc123",
  "leaderboard": [
    {
      "rank": 1,
      "participant_id": "p_56ef78gh",
      "score": 100
    },
    {
      "rank": 2,
      "participant_id": "p_90ij12kl",
      "score": 50
    },
    {
      "rank": 3,
      "participant_id": "p_12ab34cd",
      "score": 20
    }
  ],
  "timestamp": "2026-02-02T10:15:00Z"
}
```

**Message 2 - Leaderboard Update (Khi có user submit answer):**
```json
{
  "type": "leaderboard_update",
  "quiz_id": "quiz_abc123",
  "leaderboard": [
    {
      "rank": 1,
      "participant_id": "p_56ef78gh",
      "score": 110
    },
    {
      "rank": 2,
      "participant_id": "p_90ij12kl",
      "score": 50
    },
    {
      "rank": 3,
      "participant_id": "p_12ab34cd",
      "score": 20
    }
  ],
  "timestamp": "2026-02-02T10:16:30Z"
}
```

**Message 3 - Quiz Ended:**
```json
{
  "type": "quiz_ended",
  "quiz_id": "quiz_abc123",
  "message": "Quiz has ended",
  "timestamp": "2026-02-02T10:30:00Z"
}
```

---

### **Step 4.3: WebSocket Error Scenarios**

**Error 1: participantId missing**
```bash
ws://localhost:8080/quizzes/quiz_abc123/ws
# Response: 400 Bad Request
# {"error": "participantId query parameter is required"}
```

**Error 2: Participant chưa join quiz**
```bash
ws://localhost:8080/quizzes/quiz_abc123/ws?participantId=p_fake_id
# Response: 400 Bad Request
# {"error": "participant not joined this quiz"}
```

**Error 3: Quiz không tồn tại**
```bash
ws://localhost:8080/quizzes/quiz_invalid/ws?participantId=p_12ab34cd
# WebSocket connect nhưng không nhận được initial_state
```

---

### **Step 4.4: Real-time Testing Flow**

**Terminal 1 - WebSocket Client 1:**
```powershell
wscat -c "ws://localhost:8080/quizzes/quiz_abc123/ws?participantId=p_12ab34cd"
# Wait for initial_state message
```

**Terminal 2 - WebSocket Client 2:**
```powershell
wscat -c "ws://localhost:8080/quizzes/quiz_abc123/ws?participantId=p_56ef78gh"
# Wait for initial_state message
```

**Terminal 3 - Submit Answer:**
```powershell
curl -X POST "http://localhost:8080/quizzes/quiz_abc123/answer" `
  -H "Content-Type: application/json" `
  -d '{"participant_id": "p_12ab34cd", "question_id": "4", "answer": "B"}'
```

**✅ Expected Result:**
- Terminal 1 & 2 both receive `leaderboard_update` message
- Leaderboard automatically updates with new scores
- No need for refresh or polling

**Expected Messages:**
1. **On connect:** `initial_state` with current leaderboard
2. **On answer submit:** `leaderboard_update` with new leaderboard


## **🔧 Database GUI Tools**

### **PostgreSQL - HeidiSQL** ✅
Successfully connected with configuration:
- **Host:** localhost
- **Port:** 5432
- **User:** tri
- **Password:** 123456
- **Database:** rt_quiz

### **Redis - GUI Options**

#### **Option 1: RedisInsight (Recommended - Official Tool)**
1. **Download:** https://redis.io/insight/
2. **Install and open RedisInsight**
3. **Add Database:**
   - Host: `localhost`
   - Port: `6379`
   - Name: `RT-Quiz Redis`
4. **Features:**
   - ✅ Browse keys with tree view
   - ✅ View/Edit data real-time
   - ✅ Monitor performance
   - ✅ Run Redis commands


## **🎯 Success Criteria**

✅ **Phase 1:** Quiz created successfully with 10 questions, status = "started"  
✅ **Phase 2:** 3 users join, submit answers (100, 50, 20 points), scores updated real-time  
✅ **Anti-cheat:** Duplicate answer rejected  
✅ **Leaderboard:** Correct order (100, 50, 20)  
✅ **Phase 3:** Quiz ended successfully, results saved to DB  
✅ **WebSocket:** Received initial_state and leaderboard_update  
✅ **Database:** All data persisted correctly (10 questions, 30 answers total)  

---

