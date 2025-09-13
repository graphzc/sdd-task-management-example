generate:
	@echo "Generating code..."

	wiresetgen generate

	go fmt ./cmd/api/di

	wire ./...

	@echo "Code generation complete."

fmt:
	@echo "Formatting code..."

	go fmt ./...

	@echo "Code formatting complete."

start:
	@echo "Starting the application..."

	go run ./cmd/api/main.go