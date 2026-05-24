build_folder=${PWD}/cmd/server/
static_dir=${build_folder}/static
docker_compose_file=${PWD}/docker-compose.dev.yaml

build-styles:
	@mkdir -p ${PWD}/${static_dir}
	@npm run sass:build
	@cp -r ${PWD}/assets/*/**.svg ${static_dir}

build-app:
	@go build ${build_folder}

clean:
	@rm -rf ${PWD}/${build_folder}/${static_dir}
	@rm -rf ${static_dir}

run: build-styles build-app
	@${PWD}/server

up-build:
	@docker compose -f ${docker_compose_file} up --build