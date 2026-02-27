.PHONY: run build test test-coverage lint fmt clean tidy docker-build

# Variables
APP_NAME=landing-cms-api
BUILD_DIR=./bin
MAIN_PATH=./cmd/api/main.go

# Run the application
run:
	go run $(MAIN_PATH)

# Run with hot reload (requires air: go install github.com/air-verse/air@latest)
dev:
	air

# Build the binary
build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)

# Run tests
test:
	go test -v -race -coverprofile=coverage.out ./...

# Run tests with coverage report
test-coverage:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Lint (requires golangci-lint)
lint:
	golangci-lint run ./...

# Format code
fmt:
	gofmt -s -w .
	goimports -w .

# Tidy dependencies
tidy:
	go mod tidy
	go mod verify

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Download dependencies
deps:
	go mod download

# Build Docker image
docker-build:
	docker build -t $(APP_NAME):latest .

# Run database migrations
migrate:
	@echo "Run migrations manually with psql:"
	@echo "psql \$$DATABASE_URL -f ../../scripts/migrations/001_create_users.sql"
	@echo "psql \$$DATABASE_URL -f ../../scripts/migrations/002_create_sites.sql"
	@echo "psql \$$DATABASE_URL -f ../../scripts/migrations/003_create_pages.sql"
	@echo "psql \$$DATABASE_URL -f ../../scripts/migrations/004_create_content.sql"
	@echo "psql \$$DATABASE_URL -f ../../scripts/migrations/005_create_media.sql"
	@echo "psql \$$DATABASE_URL -f ../../scripts/migrations/006_create_navigation.sql"
	@echo "psql \$$DATABASE_URL -f ../../scripts/migrations/007_create_components.sql"
	@echo "psql \$$DATABASE_URL -f ../../scripts/migrations/008_create_audit.sql"

# Generate mock files (requires mockery)
mocks:
	mockery --all --dir=./internal/repository --output=./internal/repository/mocks
	mockery --all --dir=./internal/service --output=./internal/service/mocks

# Security scan (requires gosec)
security:
	gosec ./...

# Check for vulnerabilities
vuln:
	govulncheck ./...
