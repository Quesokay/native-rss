dev:
	@echo "Compiling templates..."
	@templ generate
	@echo "Compiling Tailwind CSS..."
	@echo "Starting server..."
	@go run cmd/server/main.go