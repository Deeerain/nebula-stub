bin=app
static_dir=static
docker_compose_file=${PWD}/docker-compose.dev.yaml

build-styles:
	@mkdir -p ${PWD}/${static_dir}
	@npm run sass:build
	@cp -r ${PWD}/assets/*/**.svg ${PWD}/${static}

build-app:
	@go build -o ${PWD}/${bin}

clean:
	@rm -rf ${PWD}/${static_dir}
	@rm -rf ${PWD}/${bin}

run: build-styles build-app
	@${PWD}/${bin}

up-build:
	@docker compose -f ${docker_compose_file} up --build