build-styles:
	@mkdir -p ./static
	@npm run sass:build

build-app:
	@go build -o server

clean:
	@rm -rf ./static

run: build-styles build-app
	@./server