include ./build/common/Makefile

DEV_ENV_PATH=build/dev
DOCKER_DEV_ENV_PATH=$(DEV_ENV_PATH)/docker


.PHONY: integration-tests
integration-tests: ## Run go test with integration flags
	@HASURA_AUTH_BEARER=$(shell make dev-jwt) \
	 TEST_S3_ACCESS_KEY=$(shell make dev-s3-access-key) \
	 TEST_S3_SECRET_KEY=$(shell make dev-s3-secret-key) \
	 GIN_MODE=release \
		richgo test -tags=integration $(GOTEST_OPTIONS) ./...   # -run=TestGetFileByID


.PHONY: dev-env-up-short
dev-env-up-short:  ## Starts development environment without hasura-storage
	docker-compose -f ${DOCKER_DEV_ENV_PATH}/docker-compose.yaml up -d postgres graphql-engine minio


.PHONY: dev-env-up-hasura
dev-env-up-hasura: build-docker-image  ## Starts development environment but only hasura-storage
	docker-compose -f ${DOCKER_DEV_ENV_PATH}/docker-compose.yaml up -d storage


.PHONY: _dev-env-up
_dev-env-up:
	docker-compose -f ${DOCKER_DEV_ENV_PATH}/docker-compose.yaml up -d


.PHONY: _dev-env-down
_dev-env-down:
	docker-compose -f ${DOCKER_DEV_ENV_PATH}/docker-compose.yaml down


.PHONY: _dev-env-build
_dev-env-build: build-docker-image
	docker tag $(NAME):$(VERSION) $(NAME):dev
	docker-compose -f ${DOCKER_DEV_ENV_PATH}/docker-compose.yaml build


.PHONY: dev-jwt
dev-jwt:  ## return a jwt valid for development environment
	@sh ./$(DEV_ENV_PATH)/jwt-gen/get-jwt.sh
	@sleep 2


.PHONY: dev-s3-access-key
dev-s3-access-key:  ## return s3 access key for development environment
	@docker exec -i hasura-storage-tests-minio bash -c 'echo $$MINIO_ROOT_USER'


.PHONY: dev-s3-secret-key
dev-s3-secret-key:  ## restun s3 secret key for development environment
	@docker exec -i hasura-storage-tests-minio bash -c 'echo $$MINIO_ROOT_PASSWORD'


.PHONY: migrations-add
migrations-add:  ## add a migration with NAME in the migrations folder
	migrate create -dir ./migrations/postgres -ext sql -seq $(MIGRATION_NAME)
