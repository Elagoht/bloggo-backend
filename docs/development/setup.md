# Development Environment Setup

## 📋 Table of Contents
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Database Setup](#database-setup)
- [Running the Application](#running-the-application)
- [Development Workflow](#development-workflow)
- [Testing](#testing)
- [Troubleshooting](#troubleshooting)
- [IDE Configuration](#ide-configuration)
- [Docker Development](#docker-development)

## ✅ Prerequisites

### Required Software

1. **Go Programming Language** (Version 1.24.3 or later)
   - Download from [golang.org](https://golang.org/dl/)
   - Follow installation instructions for your operating system

2. **Git** (Version 2.0 or later)
   - Download from [git-scm.com](https://git-scm.com/)
   - Required for version control

3. **SQLite3** (Version 3.35 or later)
   - Usually pre-installed on most systems
   - Download from [sqlite.org](https://sqlite.org/download.html) if needed

4. **Code Editor** (Recommended)
   - Visual Studio Code with Go extension
   - GoLand (JetBrains)
   - Vim/Neovim with Go plugins

### System Requirements

- **Operating System**: Windows 10+, macOS 10.15+, or Linux (Ubuntu 18.04+)
- **RAM**: Minimum 4GB, recommended 8GB
- **Storage**: Minimum 2GB free space
- **Network**: Internet connection for dependency downloads

## 🚀 Installation

### 1. Clone the Repository

```bash
git clone https://github.com/your-org/bloggo.git
cd bloggo
```

### 2. Install Go Dependencies

```bash
# Download and install all dependencies
go mod download

# Verify dependencies
go mod verify

# Tidy up dependencies (remove unused ones)
go mod tidy
```

### 3. Build the Application

```bash
# Build for current platform
go build -o bloggo ./cli

# Or build for specific platforms
go build -o bloggo-linux ./cli          # Linux
go build -o bloggo-windows.exe ./cli    # Windows
go build -o bloggo-macos ./cli          # macOS
```

### 4. Verify Installation

```bash
# Check if the application builds successfully
go build ./...

# Run the application to verify
./bloggo --help
```

## ⚙️ Configuration

### 1. Initial Configuration

The first time you run the application, it will automatically generate a configuration file:

```bash
# Run the application to generate config
./bloggo

# This creates bloggo-config.json with random secrets
```

### 2. Configuration File Structure

```json
{
  "port": 8723,
  "JWTSecret": "random-32-character-secret-key",
  "accessTokenDuration": 900,
  "refreshTokenDuration": 604800,
  "geminiApiKey": "",
  "trustedFrontendKey": "random-32-character-frontend-key"
}
```

### 3. Environment-Specific Configuration

For development, you can override configuration using environment variables:

```bash
# Linux/macOS
export BLOGGO_PORT=8080
export BLOGGO_JWT_SECRET="your-secret-key"

# Windows
set BLOGGO_PORT=8080
set BLOGGO_JWT_SECRET="your-secret-key"
```

### 4. Development Configuration

Create a development configuration file:

```json
{
  "port": 8080,
  "JWTSecret": "dev-secret-key-not-for-production",
  "accessTokenDuration": 3600,
  "refreshTokenDuration": 86400,
  "geminiApiKey": "your-gemini-api-key",
  "trustedFrontendKey": "dev-frontend-key"
}
```

## 🗄️ Database Setup

### 1. Automatic Database Initialization

The application automatically creates and initializes the SQLite database:

```bash
# The database file will be created automatically
./bloggo

# Database file location
ls -la bloggo.sqlite
```

### 2. Manual Database Setup (Optional)

If you want to set up the database manually:

```bash
# Install SQLite3 command line tool
sqlite3 --version

# Create database
sqlite3 bloggo.sqlite "VACUUM;"

# Verify database schema
sqlite3 bloggo.sqlite ".schema"
```

### 3. Database Migration

If you need to run database migrations manually:

```bash
# Go into the database directory
cd internal/db

# Run migrations (this is done automatically on startup)
go run *.go
```

### 4. Database Backup and Restore

```bash
# Backup database
cp bloggo.sqlite backup/bloggo_$(date +%Y%m%d_%H%M%S).sqlite

# Restore database
cp backup/bloggo_20251004_120000.sqlite bloggo.sqlite
```

## 🏃 Running the Application

### 1. Development Mode

```bash
# Run the application
./bloggo

# Or use go run for development
go run ./cli/main.go
```

The application will start on `http://localhost:8723` (or your configured port).

### 2. Development Server with Hot Reload

For development with automatic reloading, use air:

```bash
# Install air (if not already installed)
go install github.com/cosmtrek/air@latest

# Create air configuration file
air init

# Run with hot reload
air
```

### 3. Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### 4. Building for Production

```bash
# Build optimized binary
go build -ldflags="-s -w" -o bloggo ./cli

# Or build with version information
go build -ldflags="-X main.version=1.0.0" -o bloggo ./cli
```

## 🔄 Development Workflow

### 1. Code Organization

```
bloggo/
├── cli/                    # Application entry point
├── internal/              # Internal packages
│   ├── app/              # Application core
│   ├── config/           # Configuration management
│   ├── db/               # Database layer
│   ├── middleware/       # HTTP middleware
│   ├── module/           # Feature modules
│   └── utils/            # Utility packages
├── docs/                 # Documentation
├── pkg/                  # Public packages (if any)
├── scripts/              # Build and utility scripts
├── uploads/              # File uploads directory
├── go.mod                # Go module definition
├── go.sum                # Go module checksums
├── Makefile              # Build automation
└── README.md            # Project documentation
```

### 2. Typical Development Workflow

```bash
# 1. Create a new feature branch
git checkout -b feature/new-feature

# 2. Make your changes
# Edit files...

# 3. Run tests
go test ./...

# 4. Run application to test manually
go run ./cli/main.go

# 5. Build to ensure no compilation errors
go build ./...

# 6. Run static analysis
go vet ./...
gofmt -s -w .

# 7. Commit changes
git add .
git commit -m "feat: add new feature"

# 8. Push and create pull request
git push origin feature/new-feature
```

### 3. Code Quality Checks

```bash
# Run all quality checks
make check

# Or individual checks
go vet ./...              # Static analysis
gofmt -d .               # Format checking
golint ./...              # Linting (if installed)
go test ./...             # Unit tests
```

### 4. Adding New Features

1. **Create a new module**:
```bash
mkdir internal/module/new_feature
touch internal/module/new_feature/{handler,service,repository,models}.go
```

2. **Implement the module interface**:
```go
type NewFeatureModule struct{}

func NewModule() *NewFeatureModule {
    return &NewFeatureModule{}
}

func (m *NewFeatureModule) RegisterModule(router *chi.Mux) {
    router.Route("/api/new-feature", m.registerRoutes)
}
```

3. **Register the module**:
```go
// In cli/main.go
modules := []module.Module{
    // ... existing modules
    new_feature.NewModule(),
}
```

## 🧪 Testing

### 1. Unit Tests

```bash
# Run all unit tests
go test ./...

# Run tests for a specific package
go test ./internal/module/post

# Run tests with verbose output
go test -v ./...

# Run specific test
go test -run TestCreatePost ./internal/module/post
```

### 2. Integration Tests

```bash
# Run integration tests
go test -tags=integration ./...

# Run tests with test database
export BLOGGO_TEST_DB="test_bloggo.sqlite"
go test ./...
```

### 3. Test Database Setup

```go
// In your test files
func setupTestDB(t *testing.T) *sql.DB {
    db, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        t.Fatal(err)
    }

    // Initialize test schema
    InitializeTables(db)
    SeedDatabase(db)

    return db
}
```

### 4. API Testing

```bash
# Run API tests
go test ./api/...

# Run tests with coverage
go test -coverprofile=api_coverage.out ./api/...
```

### 5. Benchmark Tests

```bash
# Run benchmarks
go test -bench=. ./...

# Run specific benchmark
go test -bench=BenchmarkCreatePost ./internal/module/post
```

### 6. Test Examples

```go
func TestCreatePost(t *testing.T) {
    // Setup
    db := setupTestDB(t)
    service := post.NewService(post.NewRepository(db))

    // Test
    post := &models.Post{
        Title:   "Test Post",
        Content: "Test content",
        AuthorID: 1,
    }

    err := service.Create(post)

    // Assert
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }

    if post.ID == 0 {
        t.Error("Expected post ID to be set")
    }
}
```

## 🔧 Troubleshooting

### Common Issues

#### 1. Go Installation Issues

**Problem**: `go: command not found`

**Solution**:
```bash
# Check if Go is installed
go version

# If not installed, download and install Go
# Add Go bin directory to PATH
export PATH=$PATH:/usr/local/go/bin
```

#### 2. Permission Issues

**Problem**: `permission denied` when accessing files

**Solution**:
```bash
# Check file permissions
ls -la bloggo.sqlite

# Fix permissions (if needed)
chmod 644 bloggo.sqlite
chmod 755 uploads/
```

#### 3. Port Already in Use

**Problem**: `address already in use` error

**Solution**:
```bash
# Find process using the port
lsof -i :8723  # Linux/macOS
netstat -ano | findstr :8723  # Windows

# Kill the process
kill -9 <PID>  # Linux/macOS
taskkill /PID <PID> /F  # Windows

# Or use a different port
export BLOGGO_PORT=8080
```

#### 4. Database Issues

**Problem**: `database is locked` error

**Solution**:
```bash
# Check if another instance is running
ps aux | grep bloggo

# Kill existing instances
killall bloggo

# Remove lock files
rm -f bloggo.sqlite-wal bloggo.sqlite-shm
```

#### 5. Dependency Issues

**Problem**: `module not found` errors

**Solution**:
```bash
# Clean module cache
go clean -modcache

# Re-download dependencies
go mod download

# Rebuild
go build ./...
```

### Debug Mode

Enable debug logging for troubleshooting:

```bash
# Set debug environment variable
export BLOGGO_DEBUG=true
export BLOGGO_LOG_LEVEL=debug

# Run application
./bloggo
```

### Performance Profiling

Enable CPU and memory profiling:

```bash
# Enable profiling
export BLOGGO_PROFILE_CPU=true
export BLOGGO_PROFILE_MEMORY=true

# Run application
./bloggo

# Profile files will be created:
# cpu.prof
# mem.prof
```

Analyze profiles:

```bash
# CPU profile
go tool pprof cpu.prof

# Memory profile
go tool pprof mem.prof
```

## 💻 IDE Configuration

### Visual Studio Code

1. **Install Go Extension**:
   - Open VS Code
   - Go to Extensions
   - Search for "Go" and install the official Go extension

2. **Configure Go Extension**:
   ```json
   {
     "go.useLanguageServer": true,
     "go.formatTool": "goimports",
     "go.lintTool": "golangci-lint",
     "go.testFlags": ["-v"],
     "go.coverOnSave": true,
     "go.coverageDecorator": {
       "type": "gutter",
       "coveredHighlightColor": "rgba(64,128,64,0.5)",
       "uncoveredHighlightColor": "rgba(128,64,64,0.25)"
     }
   }
   ```

3. **Install Go Tools**:
   ```bash
   # In VS Code command palette (Ctrl+Shift+P)
   # Run "Go: Install/Update Tools"
   # Select all tools and install
   ```

4. **Debug Configuration**:
   ```json
   {
     "version": "0.2.0",
     "configurations": [
       {
         "name": "Launch Bloggo",
         "type": "go",
         "request": "launch",
         "mode": "auto",
         "program": "${workspaceFolder}/cli",
         "env": {
           "BLOGGO_PORT": "8080"
         }
       }
     ]
   }
   ```

### GoLand (JetBrains)

1. **Import Project**:
   - Open GoLand
   - File → Open → Select bloggo directory
   - Go module should be detected automatically

2. **Configuration**:
   - Go to Settings → Go
   - Set GOPATH and GOROOT if not detected
   - Enable code formatting on save
   - Configure run configurations

3. **Run Configuration**:
   - Edit Configurations → Add → Go Build
   - File: `cli/main.go`
   - Working directory: Project root
   - Environment variables: `BLOGGO_PORT=8080`

### Vim/Neovim

1. **Install vim-go plugin**:
   ```vim
   Plug 'fatih/vim-go'
   ```

2. **Configuration**:
   ```vim
   let g:go_fmt_command = "goimports"
   let g:go_fmt_autosave = 1
   let g:go_autodetect_gopath = 1
   let g:go_test_timeout = '10s'
   let g:go_test_show_name = 1
   ```

## 🐳 Docker Development

### 1. Dockerfile

```dockerfile
# Multi-stage build
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o bloggo ./cli

FROM alpine:latest
RUN apk --no-cache add ca-certificates sqlite
WORKDIR /root/

COPY --from=builder /app/bloggo .
COPY --from=builder /app/internal/db/initialize.queries.go ./queries/

EXPOSE 8723
CMD ["./bloggo"]
```

### 2. Docker Compose

```yaml
version: '3.8'

services:
  bloggo:
    build: .
    ports:
      - "8723:8723"
    volumes:
      - ./uploads:/root/uploads
      - ./bloggo.sqlite:/root/bloggo.sqlite
    environment:
      - BLOGGO_PORT=8723
      - BLOGGO_DEBUG=true
    restart: unless-stopped

  # Optional: Database admin interface
  sqlite-browser:
    image: coleifer/sqlite-web
    ports:
      - "8080:8080"
    volumes:
      - ./bloggo.sqlite:/data/bloggo.db
    command: sqlite_web -H 0.0.0.0 /data/bloggo.db
```

### 3. Development Commands

```bash
# Build Docker image
docker build -t bloggo .

# Run with Docker Compose
docker-compose up -d

# View logs
docker-compose logs -f bloggo

# Stop containers
docker-compose down

# Execute commands in container
docker-compose exec bloggo sh
```

### 4. Docker Development Workflow

```bash
# 1. Build and run
docker-compose up -d

# 2. View logs
docker-compose logs -f

# 3. Test API
curl http://localhost:8723/health

# 4. Access database browser
open http://localhost:8080

# 5. Stop and clean up
docker-compose down -v
```

---

**Document Version**: 1.0.0
**Last Updated**: October 4, 2025
**Author**: Bloggo Development Team
**Reviewers**: Development Infrastructure Committee