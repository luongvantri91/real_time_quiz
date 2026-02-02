# RT-Quiz Project Structure

## Directory Overview

```
rt-quiz/
├── main.go                     # Application entry point
├── go.mod                      # Go dependencies
├── Dockerfile                  # Container build configuration
├── docker-compose.yml          # Multi-container orchestration
├── Makefile                    # Build and run commands
├── .env.example                # Environment variables template
│
├── cmd/                        # Application startup logic
│   ├── server_dependencies.go  # Dependency injection setup
│   ├── server_env.go           # Environment configuration
│   └── server_router.go        # HTTP route definitions
│
├── handlers/                   # HTTP/WebSocket request handlers
│   ├── http_handler.go         # User API handlers
│   ├── admin_handler.go        # Admin API handlers
│   ├── ws_handler.go           # WebSocket handler
│   └── validation.go           # Input validation logic
│
├── services/                   # Business logic layer
│   ├── quiz_service.go         # Quiz lifecycle management
│   └── participant_service.go  # Participant operations
│
├── repository/                 # Data access interfaces
│   ├── quiz_repository.go      # Quiz data operations
│   ├── participant_repository.go  # Participant data operations
│   ├── answer_repository.go    # Answer data operations
│   └── result_repository.go    # Result data operations
│
├── infrastructure/             # External service implementations
│   └── postgres/
│       └── pg_repository.go    # PostgreSQL implementations
│
├── client/                     # External client libraries
│   ├── redis/
│   │   └── redis.go            # Redis client (cache + pub/sub)
│   └── ws/
│       └── ws.go               # WebSocket utilities
│
├── models/                     # Data models and DTOs
│   ├── quiz.go                 # Quiz models
│   ├── question.go             # Question models
│   ├── participant.go          # Participant models
│   ├── answer.go               # Answer models
│   ├── result.go               # Result models
│   └── leaderboard_entry.go    # Leaderboard models
│
├── migrations/                 # Database schema
│   ├── 001_init.sql            # Initial schema
│   └── 003_admin_quiz_management.sql  # Admin features
│
├── docs/                       # API documentation
│   ├── create_quiz.md          # POST /admin/quizzes
│   ├── get_quiz.md             # GET /admin/quizzes/:quizId
│   ├── list_quizzes.md         # GET /admin/quizzes
│   ├── get_quiz_status.md      # GET /admin/quizzes/:quizId/status
│   ├── init_quiz.md            # POST /admin/quizzes/:quizId/init
│   ├── end_quiz.md             # POST /admin/quizzes/:quizId/end
│   ├── add_question.md         # POST /admin/quizzes/:quizId/questions
│   ├── delete_question.md      # DELETE /admin/quizzes/:quizId/questions/:questionId
│   ├── list_questions.md       # GET /admin/quizzes/:quizId/questions
│   ├── join_quiz.md            # POST /quizzes/:quizId/join
│   ├── submit_answer.md        # POST /quizzes/:quizId/answer
│   ├── get_leaderboard.md      # GET /quizzes/:quizId/leaderboard
│   └── websocket.md            # GET /quizzes/:quizId/ws
│
├── TEST_QUIZ_1.md              # Test data for Quiz 1
├── TEST_QUIZ_2.md              # Test data for Quiz 2
├── ARCHITECTURE.md             # System architecture documentation
├── PROJECT_STRUCTURE.md        # This file
└── README.md                   # Project overview
```

## Layer Responsibilities

### 1. Entry Point (`main.go`)
- Initialize database connections (PostgreSQL, Redis)
- Setup dependencies (repositories, services, handlers)
- Start HTTP server with routing configuration

### 2. Command Layer (`cmd/`)
- **server_dependencies.go**: Dependency Injection container
- **server_env.go**: Environment variables loading
- **server_router.go**: HTTP routes registration

### 3. Handlers Layer (`handlers/`)
**Responsibilities:**
- Parse HTTP request (JSON, params, headers)
- Validate input
- Call service layer
- Format HTTP response

**Files:**
- `http_handler.go`: User endpoints (join, submit, leaderboard)
- `admin_handler.go`: Admin endpoints (create, init, end quiz)
- `ws_handler.go`: WebSocket realtime updates
- `validation.go`: Shared validation functions

### 4. Services Layer (`services/`)
**Responsibilities:**
- Business logic implementation
- Orchestrate multiple repositories
- Coordinate DB + Redis operations
- Error handling and logging

**Files:**
- `quiz_service.go`: Quiz lifecycle (create, init, end)
- `participant_service.go`: Join quiz, submit answers

### 5. Repository Layer (`repository/`)
**Responsibilities:**
- Define data access interfaces
- Abstract database operations
- No business logic, only data operations

**Files:**
- `quiz_repository.go`: Quiz CRUD
- `participant_repository.go`: Participant CRUD
- `answer_repository.go`: Answer storage
- `result_repository.go`: Result storage

### 6. Infrastructure Layer (`infrastructure/`)
**Responsibilities:**
- Implement repository interfaces
- Concrete database implementations
- Database-specific code

**Files:**
- `postgres/pg_repository.go`: PostgreSQL implementations

### 7. Client Layer (`client/`)
**Responsibilities:**
- External service clients
- Cache operations (Redis)
- WebSocket utilities

**Files:**
- `redis/redis.go`: Redis cache + pub/sub
- `ws/ws.go`: WebSocket helpers

### 8. Models Layer (`models/`)
**Responsibilities:**
- Data structures (entities, DTOs)
- Request/Response models
- Database models

**Files:**
- `quiz.go`, `question.go`, `participant.go`, `answer.go`, `result.go`, `leaderboard_entry.go`

### 9. Migrations (`migrations/`)
**Responsibilities:**
- Database schema definitions
- SQL migration files

**Files:**
- `001_init.sql`: Core tables
- `003_admin_quiz_management.sql`: Admin features

### 10. Documentation (`docs/`)
**Responsibilities:**
- API endpoint documentation
- Request/Response examples
- Integration guides

## Data Flow Example

```
HTTP Request
    ↓
handlers/admin_handler.go (CreateQuiz)
    ↓
services/quiz_service.go (CreateQuiz method)
    ↓
repository/quiz_repository.go (SaveQuiz interface)
    ↓
infrastructure/postgres/pg_repository.go (SaveQuiz implementation)
    ↓
PostgreSQL Database
```

## Dependencies

### External Packages
- `github.com/labstack/echo/v4` - HTTP framework
- `github.com/lib/pq` - PostgreSQL driver
- `github.com/go-redis/redis/v9` - Redis client
- `github.com/gorilla/websocket` - WebSocket
- `github.com/google/uuid` - UUID generation

### Internal Dependencies Flow
```
main.go
  ├── cmd/server_dependencies.go
  │     ├── infrastructure/postgres
  │     ├── client/redis
  │     ├── repository/*
  │     ├── services/*
  │     └── handlers/*
  └── cmd/server_router.go
```

## Configuration Files

### `docker-compose.yml`
- PostgreSQL service (port 5432)
- Redis service (port 6379, 256MB limit)
- rt-quiz-app service (port 8080)
- Health checks and dependencies

### `Dockerfile`
- Multi-stage build (build + runtime)
- Go 1.x base image
- Expose port 8080

### `Makefile`
- `make build`: Build binary
- `make run`: Run with docker-compose
- `make test`: Run tests
- `make clean`: Clean build artifacts

### `.env.example`
```
DB_HOST=postgres
DB_PORT=5432
DB_USER=tri
DB_PASSWORD=123456
DB_NAME=rt_quiz
REDIS_HOST=redis
REDIS_PORT=6379
SERVER_PORT=8080
```

## Testing Files

- **TEST_QUIZ_1.md**: Curl commands for Quiz 1 (basic cleaning questions)
- **TEST_QUIZ_2.md**: Curl commands for Quiz 2 (advanced cleaning questions)

## Build & Run

```bash
# Start all services
docker-compose up --build

# Access application
http://localhost:8080

# Run migrations (auto on startup)
# Check logs: docker-compose logs rt-quiz-app
```

## Code Organization Principles

1. **Separation of Concerns**: Each layer has its own responsibility
2. **Dependency Rule**: Dependencies only point inward (handlers → services → repositories)
3. **Interface-Based Design**: Repository interfaces allow easy testing
4. **Clean Architecture**: Business logic independent of frameworks
5. **Single Responsibility**: Each file has 1 clear responsibility
