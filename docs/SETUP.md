# RT-Quiz Setup Guide

Complete guide to install, configure, and run the RT-Quiz application.

## 📋 Prerequisites

- **Docker** (version 20.10+)
- **Docker Compose** (version 1.29+)

Optional for local development:
- **Go** (version 1.21+)
- **Make**

## 🚀 Quick Start

### 1. Clone Repository

```bash
git clone <repository-url>
cd rt-quiz
```

### 2. Configure Environment

Copy the example environment file:

```bash
cp .env.example .env
```

Edit `.env` if needed (default values work for Docker setup):

```bash
DB_HOST=postgres
DB_PORT=5432
DB_USER=tri
DB_PASSWORD=123456
DB_NAME=rt_quiz
REDIS_HOST=redis
REDIS_PORT=6379
SERVER_PORT=8080
```

### 3. Start All Services

```bash
docker-compose up --build
```

This will start:
- **PostgreSQL** (port 5432)
- **Redis** (port 6379)
- **RT-Quiz App** (port 8080)

### 4. Verify Installation

```bash
# Check health endpoint
curl http://localhost:8080/health

# Expected response: {"status":"ok"}
```

## 🔧 Development Setup

### Running Locally (without Docker)

#### 1. Start Dependencies Only

```bash
docker-compose up postgres redis -d
```

#### 2. Update Environment Variables

Create `.env` for local development:

```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=tri
DB_PASSWORD=123456
DB_NAME=rt_quiz
REDIS_HOST=localhost
REDIS_PORT=6379
SERVER_PORT=8080
```

#### 3. Install Go Dependencies

```bash
go mod download
```

#### 4. Run Application

```bash
# Using Go directly
go run main.go

# Or using Make
make run

# Or build and run binary
make build
./rt-quiz
```

### Database Migrations

Migrations run automatically on application startup. The system will:
1. Connect to PostgreSQL
2. Execute migration files in order:
   - `migrations/001_init.sql` - Core tables
   - `migrations/003_admin_quiz_management.sql` - Admin features

To run migrations manually:

```bash
docker exec -it rt-quiz-postgres psql -U tri -d rt_quiz -f /migrations/001_init.sql
docker exec -it rt-quiz-postgres psql -U tri -d rt_quiz -f /migrations/003_admin_quiz_management.sql
```

## 🧪 Testing

### Test Data

Use the provided test data files to quickly test the system:

#### Test Quiz 1 (Basic Cleaning)
```bash
# Follow commands in TEST_QUIZ_1.md
# Includes: Create quiz, add questions, init, join, submit answers, end
```

#### Test Quiz 2 (Advanced Cleaning)
```bash
# Follow commands in TEST_QUIZ_2.md
# Includes: 10 advanced questions with 3 test users
```

### Example Workflow

#### 1. Create Quiz

```bash
curl -X POST http://localhost:8080/admin/quizzes \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Quiz",
    "description": "Sample quiz for testing",
    "duration_minutes": 30
  }'
```

Response:
```json
{
  "quiz_id": "quiz_a1b2c3d4e5f6",
  "title": "Test Quiz",
  "status": "pending"
}
```

#### 2. Add Questions

```bash
QUIZ_ID="quiz_a1b2c3d4e5f6"

curl -X POST http://localhost:8080/admin/quizzes/$QUIZ_ID/questions \
  -H "Content-Type: application/json" \
  -d '{
    "text": "What is 2+2?",
    "options": ["3", "4", "5", "6"],
    "correct_answer": 1,
    "points": 10
  }'
```

#### 3. Initialize Quiz

```bash
curl -X POST http://localhost:8080/admin/quizzes/$QUIZ_ID/init
```

#### 4. Join Quiz (User)

```bash
curl -X POST http://localhost:8080/quizzes/$QUIZ_ID/join \
  -H "Content-Type: application/json" \
  -d '{"username": "alice"}'
```

Response:
```json
{
  "participant_id": "p_1234567890ab",
  "message": "joined quiz successfully"
}
```

#### 5. Submit Answer

```bash
PARTICIPANT_ID="p_1234567890ab"
QUESTION_ID="q_0987654321yx"

curl -X POST http://localhost:8080/quizzes/$QUIZ_ID/answer \
  -H "Content-Type: application/json" \
  -d '{
    "participant_id": "'$PARTICIPANT_ID'",
    "question_id": "'$QUESTION_ID'",
    "answer_index": 1
  }'
```

#### 6. Get Leaderboard

```bash
curl http://localhost:8080/quizzes/$QUIZ_ID/leaderboard
```

#### 7. Connect to WebSocket

Using `wscat` (install with `npm install -g wscat`):

```bash
wscat -c "ws://localhost:8080/quizzes/$QUIZ_ID/ws?participant_id=$PARTICIPANT_ID"
```

Or using JavaScript:

```javascript
const ws = new WebSocket(
  'ws://localhost:8080/quizzes/quiz_a1b2c3d4e5f6/ws?participant_id=p_1234567890ab'
);

ws.onopen = () => console.log('Connected');
ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log('Leaderboard update:', data);
};
```

#### 8. End Quiz

```bash
curl -X POST http://localhost:8080/admin/quizzes/$QUIZ_ID/end
```

## 📊 Monitoring & Debugging

### View Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f rt-quiz-app
docker-compose logs -f postgres
docker-compose logs -f redis
```

### PostgreSQL Database

```bash
# Connect to PostgreSQL
docker exec -it rt-quiz-postgres psql -U tri -d rt_quiz

# List tables
\dt

# View quizzes
SELECT * FROM quizzes;

# View participants
SELECT * FROM quiz_participants;

# View answers
SELECT * FROM quiz_answers;

# View results
SELECT * FROM results;

# Exit
\q
```

### Redis Cache

```bash
# Connect to Redis CLI
docker exec -it rt-quiz-redis redis-cli

# View all quiz keys
KEYS quiz:*

# View leaderboard for a quiz
ZREVRANGE quiz:quiz_a1b2c3d4e5f6:leaderboard 0 -1 WITHSCORES

# View scores
HGETALL quiz:quiz_a1b2c3d4e5f6:scores

# View participants
SMEMBERS quiz:quiz_a1b2c3d4e5f6:users

# Check key TTL
TTL quiz:quiz_a1b2c3d4e5f6:leaderboard

# Exit
exit
```

### Health Check

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{"status":"ok"}
```

## 🐛 Troubleshooting

### Database Connection Failed

**Problem:** Application can't connect to PostgreSQL

**Solution:**
```bash
# Check PostgreSQL is running
docker-compose ps postgres

# Check PostgreSQL logs
docker-compose logs postgres

# Restart PostgreSQL
docker-compose restart postgres

# Verify connection manually
docker exec -it rt-quiz-postgres psql -U tri -d rt_quiz -c "SELECT 1"
```

### Redis Connection Failed

**Problem:** Application can't connect to Redis

**Solution:**
```bash
# Check Redis is running
docker-compose ps redis

# Test Redis connection
docker exec -it rt-quiz-redis redis-cli PING

# Should return: PONG

# Check Redis logs
docker-compose logs redis

# Restart Redis
docker-compose restart redis
```

### Port Already in Use

**Problem:** Port 8080, 5432, or 6379 already in use

**Solution:**

Option 1 - Stop conflicting services:
```bash
# Find process using port
lsof -i :8080  # macOS/Linux
netstat -ano | findstr :8080  # Windows

# Kill the process or stop the service
```

Option 2 - Change ports in `docker-compose.yml`:
```yaml
services:
  rt-quiz-app:
    ports:
      - "8081:8080"  # Change host port
```

### WebSocket Connection Failed

**Problem:** Can't establish WebSocket connection

**Checklist:**
1. ✅ Verify participant has joined the quiz
2. ✅ Check quiz status is `started` (not `pending` or `ended`)
3. ✅ Ensure correct WebSocket URL format
4. ✅ Check application logs for errors

**Test WebSocket:**
```bash
# Check if WebSocket endpoint is accessible
curl -i -N \
  -H "Connection: Upgrade" \
  -H "Upgrade: websocket" \
  -H "Sec-WebSocket-Key: test" \
  -H "Sec-WebSocket-Version: 13" \
  http://localhost:8080/quizzes/quiz_xxx/ws?participant_id=p_xxx
```

### Migration Errors

**Problem:** Database migrations fail

**Solution:**
```bash
# Check migration files exist
ls -la migrations/

# Reset database (WARNING: deletes all data)
docker-compose down -v
docker-compose up --build

# Or manually drop and recreate
docker exec -it rt-quiz-postgres psql -U tri -c "DROP DATABASE rt_quiz"
docker exec -it rt-quiz-postgres psql -U tri -c "CREATE DATABASE rt_quiz"
docker-compose restart rt-quiz-app
```

## 🛠️ Build Commands

### Using Make

```bash
# Build binary
make build

# Run application
make run

# Run with Docker Compose
make docker-up

# Stop Docker services
make docker-down

# Clean build artifacts
make clean

# Run tests
make test
```

### Manual Commands

```bash
# Build
go build -o rt-quiz main.go

# Run
./rt-quiz

# Build for production (optimized)
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o rt-quiz main.go

# Run tests
go test ./... -v

# Format code
go fmt ./...

# Lint code
golangci-lint run
```

## 🔐 Security Notes

### Production Deployment

Before deploying to production:

1. **Change default passwords:**
   ```bash
   DB_PASSWORD=<strong-password>
   ```

2. **Use environment-specific configs:**
   - Separate `.env` files for dev/staging/production
   - Use secret management (AWS Secrets Manager, HashiCorp Vault)

3. **Enable HTTPS/WSS:**
   - Configure TLS certificates
   - Use reverse proxy (nginx, traefik)

4. **Restrict database access:**
   - Don't expose PostgreSQL/Redis ports publicly
   - Use firewall rules

5. **Add authentication:**
   - Implement JWT or session-based auth
   - Protect admin endpoints

## 📞 Support

For issues or questions:
- Check [ARCHITECTURE.md](./ARCHITECTURE.md) for system design
- Check [PROJECT_STRUCTURE.md](./PROJECT_STRUCTURE.md) for code structure
- Review API docs in `/docs/` folder
- Check application logs: `docker-compose logs -f rt-quiz-app`

## 🎉 Next Steps

After successful setup:
1. Explore [API Documentation](./docs/)
2. Try test workflows in [TEST_QUIZ_1.md](./TEST_QUIZ_1.md)
3. Read [ARCHITECTURE.md](./ARCHITECTURE.md) to understand the system
4. Start building your quiz features!
