build-frontend:
	@cd frontend && npm run build -- -m dev

run-dev: build-frontend
	@STUB_DEV=1 go run main.go
