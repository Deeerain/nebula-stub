build-frontend:
	@cd frontend && npm run build -- -m dev

run-dev: build-frontend
	@go run main.go