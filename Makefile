db-for-tests:
	@docker compose -f ./ci/docker-compose-tests.yaml up -d

clear-db-for-tests:
	@docker compose -f ./ci/docker-compose-tests.yaml down -v

static-analysis:
	@CONTEXT=$(PWD) docker compose -f ./ci/docker-compose.yaml up --build static-analysis

build:
	@CONTEXT=$(PWD) docker compose -f ./ci/docker-compose.yaml up build-bin
