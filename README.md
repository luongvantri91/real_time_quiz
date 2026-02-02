# RT-Quiz - Real-time Quiz System

Online quiz system with real-time leaderboard using WebSocket, PostgreSQL, and Redis.

## 🎯 Features

- ✅ **Quiz Management**: Create, initialize, and end quizzes
- ✅ **Real-time Leaderboard**: WebSocket updates when users submit answers
- ✅ **Anti-Cheat**: Prevent duplicate submissions
- ✅ **Write-Through Cache**: PostgreSQL (source of truth) → Redis (cache)
- ✅ **Auto-Cleanup**: TTL-based Redis key expiration
- ✅ **Multi-Quiz Support**: Multiple quizzes run simultaneously without conflict
- ✅ **Horizontal Scaling**: Redis Pub/Sub for distributed architecture

## 🏗️ Architecture

### Tech Stack
- **Backend**: Go 1.x + Echo v4
- **Database**: PostgreSQL 16-alpine (persistent storage)
- **Cache**: Redis 7-alpine (256MB, allkeys-lru)
- **Real-time**: WebSocket (gorilla/websocket)
- **Deployment**: Docker + Docker Compose

### Architecture Pattern
```
Write-Through Cache:
User → Handler → Service → PostgreSQL (first) → Redis (cache)

Real-time Updates:
Submit Answer → Redis PUBLISH → All WebSocket Handlers → Clients
```

📖 Details: [ARCHITECTURE.md](./docs/ARCHITECTURE.md)

## 📁 Project Structure

```
rt-quiz/
├── cmd/                # Application startup
├── handlers/           # HTTP/WebSocket handlers
├── services/           # Business logic
├── repository/         # Data access interfaces
├── infrastructure/     # PostgreSQL implementations
├── client/            # Redis & WebSocket clients
├── models/            # Data structures
├── migrations/        # Database schema
└── docs/              # API documentation
```

📖 Details: [PROJECT_STRUCTURE.md](./docs/PROJECT_STRUCTURE.md)

##  API Documentation

### Admin APIs (Quiz Management)

| Method | Endpoint | Description | Doc |
|--------|----------|-------------|-----|
| POST | `/admin/quizzes` | Create quiz | [📄](./docs/create_quiz.md) |
| GET | `/admin/quizzes/:quizId` | Get quiz details | [📄](./docs/get_quiz.md) |
| GET | `/admin/quizzes` | List all quizzes | [📄](./docs/list_quizzes.md) |
| GET | `/admin/quizzes/:quizId/status` | Get quiz status | [📄](./docs/get_quiz_status.md) |
| POST | `/admin/quizzes/:quizId/init` | Initialize quiz | [📄](./docs/init_quiz.md) |
| POST | `/admin/quizzes/:quizId/end` | End quiz | [📄](./docs/end_quiz.md) |

### Admin APIs (Question Management)

| Method | Endpoint | Description | Doc |
|--------|----------|-------------|-----|
| POST | `/admin/quizzes/:quizId/questions` | Add question | [📄](./docs/add_question.md) |
| GET | `/admin/quizzes/:quizId/questions` | List questions | [📄](./docs/list_questions.md) |
| DELETE | `/admin/quizzes/:quizId/questions/:questionId` | Delete question | [📄](./docs/delete_question.md) |

### User APIs

| Method | Endpoint | Description | Doc |
|--------|----------|-------------|-----|
| POST | `/quizzes/:quizId/join` | Join quiz | [📄](./docs/join_quiz.md) |
| POST | `/quizzes/:quizId/answer` | Submit answer | [📄](./docs/submit_answer.md) |
| GET | `/quizzes/:quizId/leaderboard` | Get leaderboard | [📄](./docs/get_leaderboard.md) |
| GET | `/quizzes/:quizId/ws` | WebSocket updates | [📄](./docs/websocket.md) |

---

## 🚀 Getting Started

See [SETUP.md](./docs/SETUP.md) for installation and running instructions.

## 📖 Additional Documentation

- [ARCHITECTURE.md](./docs/ARCHITECTURE.md) - Detailed system architecture
- [PROJECT_STRUCTURE.md](./docs/PROJECT_STRUCTURE.md) - Code structure and dependencies
- [SETUP.md](./docs/SETUP.md) - Setup and running guide
- [TEST_QUIZ_1.md](./docs/TEST_QUIZ_1.md) - Test data Quiz 1
- [TEST_QUIZ_2.md](./docs/TEST_QUIZ_2.md) - Test data Quiz 2
