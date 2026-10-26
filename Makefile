build-frontend:
	@cd frontend && npm run build -- -m dev

run-dev: build-frontend
	@set -a; [ -f .env ] && . ./.env; set +a; STUB_DEV=1 go run main.go
