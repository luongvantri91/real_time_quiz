# RT-Quiz System Architecture

## Overview

RT-Quiz is a real-time online quiz system with Clean Architecture, using PostgreSQL as the source of truth and Redis as a cache layer for high performance.

## Architecture Pattern

### Write-Through Cache Pattern
```
User Request → Handler → Service → Database (PostgreSQL) → Cache (Redis)
                                         ↓
                                   Success/Failure
```

**Principles:**
- PostgreSQL is the source of truth (accurate data, persistent)
- Redis is the cache layer (high speed, realtime updates)
- All write operations: **DB first, then Redis**
- Redis failure → warning log, does not fail request

### Layer Architecture

```
┌─────────────────────────────────────┐
│         Presentation Layer          │
│  (handlers/, HTTP + WebSocket)      │
└─────────────────────────────────────┘
              ↓
┌─────────────────────────────────────┐
│         Business Logic Layer        │
│  (services/, quiz + participant)    │
└─────────────────────────────────────┘
              ↓
┌─────────────────────────────────────┐
│         Data Access Layer           │
│  (repository/, infrastructure/)     │
└─────────────────────────────────────┘
              ↓
┌──────────────────┬──────────────────┐
│   PostgreSQL     │      Redis       │
│  (Source of      │  (Cache +        │
│   Truth)         │   Pub/Sub)       │
└──────────────────┴──────────────────┘
```

## Components

### 1. Handlers Layer
- **http_handler.go**: User APIs (join, submit, leaderboard)
- **admin_handler.go**: Admin APIs (create, init, end quiz)
- **ws_handler.go**: WebSocket realtime updates
- **validation.go**: Request validation logic

### 2. Services Layer
- **quiz_service.go**: Quiz lifecycle (create, init, end)
- **participant_service.go**: Participant logic (join, submit answers)

### 3. Repository Layer
- **quiz_repository.go**: Quiz CRUD operations
- **participant_repository.go**: Participant data access
- **answer_repository.go**: Answer storage
- **result_repository.go**: Final results management

### 4. Infrastructure Layer
- **postgres/pg_repository.go**: PostgreSQL implementation
- **redis/redis.go**: Redis operations (cache, pub/sub)
- **ws/ws.go**: WebSocket utilities

## Data Flow

### 1. Quiz Creation Flow
```
Admin → POST /admin/quizzes
  → CreateQuiz handler
  → Save to PostgreSQL (status: pending)
  → Return quiz_id
```

### 2. Quiz Initialization Flow
```
Admin → POST /admin/quizzes/{id}/init
  → InitQuiz handler
  → Update DB (status: pending → started)
  → Calculate TTL (duration + 24h)
  → Cache questions to Redis
  → Set TTL for all Redis keys
  → Return success
```

### 3. User Join Flow
```
User → POST /quizzes/{id}/join
  → JoinQuiz handler
  → Validate quiz status
  → Save to PostgreSQL (quiz_participants)
  → Add to Redis SET (quiz:{id}:users)
  → Initialize score in Redis HASH
  → Return participant_id
```

### 4. Submit Answer Flow
```
User → POST /quizzes/{id}/answer
  → SubmitAnswer handler
  → Check anti-cheat (Redis SET)
  → Validate answer correctness
  → Save to PostgreSQL (quiz_answers) [SYNCHRONOUS]
  → Update score in Redis HASH
  → Update leaderboard in Redis ZSET
  → Publish event to Redis Pub/Sub
  → Return result
```

### 5. Real-time Updates Flow
```
User → GET /quizzes/{id}/ws
  → HandleWebSocket handler
  → Upgrade HTTP → WebSocket
  → Send initial leaderboard
  → Subscribe to Redis Pub/Sub
  → Event loop:
    - Redis message → Broadcast to client
    - Context done → Close connection
```

### 6. Quiz End Flow
```
Admin → POST /admin/quizzes/{id}/end
  → EndQuiz handler
  → Update DB (status: started → ended)
  → Get leaderboard from Redis
  → Save final results to PostgreSQL (results table)
  → Return final leaderboard
```

## Redis Data Structures

### Key Naming Convention
```
quiz:{quiz_id}:users          → SET (participant IDs)
quiz:{quiz_id}:scores         → HASH (participant_id → score)
quiz:{quiz_id}:leaderboard    → ZSET (sorted by score)
quiz:{quiz_id}:answered       → SET (participant_id:question_id)
quiz:{quiz_id}:question:{qid} → HASH (question details)
quiz:updates:{quiz_id}        → Pub/Sub channel
```

### TTL Strategy
- **Default TTL**: `quiz_duration + 24 hours`
- **Applied during**: InitQuiz operation
- **Auto-cleanup**: Redis automatically removes expired keys
- **Memory Policy**: allkeys-lru (evict least recently used)
- **Max Memory**: 256MB

## Database Schema

### Core Tables
- **quizzes**: Quiz metadata (title, duration, status)
- **questions**: Quiz questions and options
- **quiz_participants**: Participants who joined
- **quiz_answers**: All submitted answers
- **results**: Final quiz results (saved at end)

### Status Workflow
```
pending → started → ended
   ↓         ↓        ↓
Create    Init     End
```

## Concurrency Model

### WebSocket Connections
- **1 Goroutine per WebSocket connection**
- **Channel-based communication** (Redis Pub/Sub → Handler → Client)
- **Select statement** for event multiplexing

### Redis Pub/Sub Pattern
```
Submit Answer → Redis PUBLISH
       ↓
Multiple Subscribers (WebSocket Handlers)
       ↓
Unicast to Individual Clients
```

**Note**: This is "distributed unicast", NOT a true broadcast hub. Each handler has 1 connection with 1 client.

## Scaling Considerations

### Horizontal Scaling
- **Stateless handlers**: Can run multiple instances
- **Redis Pub/Sub**: Distribute notifications across servers
- **PostgreSQL connection pooling**: Limit connections per instance

### Performance Optimization
- **Redis cache**: Reduce DB reads for leaderboard
- **WebSocket**: Persistent connections, no HTTP overhead
- **TTL auto-cleanup**: No manual cleanup jobs needed

## Security Considerations

### Anti-Cheat
- Redis SET `quiz:{id}:answered` tracks submitted questions
- Prevent duplicate submissions per participant

### Input Validation
- Required fields validation
- Duration limits (5-180 minutes)
- Quiz status checks before operations

### Error Handling
- Redis failures → log warnings, continue operation
- Database failures → return error to client
- WebSocket errors → close connection gracefully

## Monitoring & Observability

### Key Metrics
- **Quiz status transitions**: pending → started → ended
- **Redis hit rate**: Cache effectiveness
- **WebSocket connections**: Active realtime users
- **Database query latency**: Performance bottlenecks

### Logging
- **Service layer**: Business logic events
- **Handler layer**: Request/response logs
- **Redis failures**: Warning logs for cache misses

## Deployment

### Docker Compose
```
┌─────────────────┐
│   PostgreSQL    │ (Port 5432, Health check)
└─────────────────┘
┌─────────────────┐
│     Redis       │ (Port 6379, Health check)
└─────────────────┘
┌─────────────────┐
│   rt-quiz-app   │ (Port 8080)
│  (depends on    │
│   DB + Redis)   │
└─────────────────┘
```

### Environment Variables
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`
- `REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD`
- `SERVER_PORT` (default: 8080)

## Technology Stack

- **Language**: Go 1.x
- **Web Framework**: Echo v4
- **Database**: PostgreSQL 16-alpine
- **Cache**: Redis 7-alpine (256MB, allkeys-lru)
- **WebSocket**: gorilla/websocket
- **Container**: Docker + Docker Compose

## Design Principles

1. **Single Source of Truth**: PostgreSQL is the primary database
2. **Cache Aside**: Redis is cache, not primary storage
3. **Fail Safe**: Redis failure does not crash the system
4. **Idempotency**: Submit answer for already answered question → reject
5. **Clean Architecture**: Separation of concerns across layers
6. **Real-time First**: WebSocket for instant updates
