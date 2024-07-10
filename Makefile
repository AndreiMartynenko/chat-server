ifneq ($(ENV),)
	include $(ENV).env
endif

LOCAL_BIN:=$(CURDIR)/bin
BINARY_NAME=main
CLI_BINARY_NAME=chat-client
CONFIG=$(ENV).env
LOCAL_MIGRATION_DIR=$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN="host=localhost port=$(POSTGRES_PORT_LOCAL) dbname=$(POSTGRES_DB) user=$(POSTGRES_USER) password=$(POSTGRES_PASSWORD) sslmode=disable"
TLS_PATH=tls
TESTS_PATH=github.com/AndreiMartynenko/chat-server/internal/service/...,github.com/AndreiMartynenko/chat-server/internal/api/...
TESTS_ATTEMPTS=5
TESTS_COVERAGE_FILE=coverage.out

# #################### #
# DEPENDENCIES & TOOLS #
# #################### #

install-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.56.1
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.32.0
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.18.0
	GOBIN=$(LOCAL_BIN) go install github.com/gojuno/minimock/v3/cmd/minimock@v3.3.6

get-protoc-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

lint:
	GOBIN=$(LOCAL_BIN) bin/golangci-lint run ./... --config .golangci.pipeline.yaml

generate-api:
	make generate-api-v1

generate-api-v1:
	mkdir -p pkg/chat_v1
	protoc --proto_path api/chat_v1 \
	--go_out=pkg/chat_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/chat_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	api/chat_v1/chat.proto

generate-cert: $(TLS_PATH)/ca.key $(TLS_PATH)/ca.pem
	openssl genrsa -out $(TLS_PATH)/chat.key 4096
	openssl req -new -key $(TLS_PATH)/chat.key -config openssl.cnf -out $(TLS_PATH)/chat.csr
	openssl x509 -req -in $(TLS_PATH)/chat.csr -CA $(TLS_PATH)/ca.pem -CAkey $(TLS_PATH)/ca.key -extfile openssl.cnf -extensions req_ext -out $(TLS_PATH)/chat.pem -days 365 -sha256
	rm -rf $(TLS_PATH)/chat.csr

generate-mocks:
	go generate ./internal/repository
	go generate ./internal/service

check-env:
ifeq ($(ENV),)
	$(error No environment specified)
endif

# ##### #
# TESTS #
# ##### #

test:
	go clean -testcache
	-go test ./... -v -covermode count -coverpkg=$(TESTS_PATH) -count $(TESTS_ATTEMPTS)

test-coverage:
	go clean -testcache
	-go test ./... -v -coverprofile=$(TESTS_COVERAGE_FILE).tmp -covermode count -coverpkg=$(TESTS_PATH) -count $(TESTS_ATTEMPTS)
	grep -v "mocks/" $(TESTS_COVERAGE_FILE).tmp > $(TESTS_COVERAGE_FILE)
	rm $(TESTS_COVERAGE_FILE).tmp
	go tool cover -html=$(TESTS_COVERAGE_FILE) -o coverage.html
	go tool cover -func=$(TESTS_COVERAGE_FILE) | grep "total"

# ##### #
# BUILD #
# ##### #

build-chat-client:
	GOOS=linux GOARCH=amd64 go build -o $(LOCAL_BIN)/${CLI_BINARY_NAME} cli/cmd/main.go

build-app:
	GOOS=linux GOARCH=amd64 go build -o $(LOCAL_BIN)/${BINARY_NAME} cmd/chat/main.go

docker-net:
	docker network create -d bridge service-net

docker-build: docker-build-app docker-build-migrator

docker-build-app: check-env
	docker buildx build --no-cache --platform linux/amd64 -t chat-server:${APP_IMAGE_TAG} --build-arg="ENV=${ENV}" --build-arg="CONFIG=${CONFIG}" .

docker-build-migrator: check-env
	docker buildx build --no-cache --platform linux/amd64 -t migrator-chat:${MIGRATOR_IMAGE_TAG} -f migrator.Dockerfile --build-arg="ENV=${ENV}" .

# ###### #
# DEPLOY #
# ###### #

docker-deploy: check-env docker-build
	docker compose --env-file=$(ENV).env up -d

local-migration-status: check-env
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} status -v

local-migration-up: check-env
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up -v

local-migration-down: check-env
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v

# #### #
# STOP #
# #### #

docker-stop: check-env
	docker compose --env-file=$(ENV).env down
