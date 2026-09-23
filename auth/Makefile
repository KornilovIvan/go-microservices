include local.env

LOCAL_BIN := $(CURDIR)/bin
LOCAL_MIGRATION_DIR := $(MIGRATION_DIR)
LOCAL_MIGRATION_DSN := host=localhost port=$(PG_PORT) dbname=$(PG_DATABASE_NAME) user=$(PG_USER) password=$(PG_PASSWORD) sslmode=disable

install-deps:
	mkdir -p "$(LOCAL_BIN)"
	GOBIN="$(LOCAL_BIN)" go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN="$(LOCAL_BIN)" go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN="$(LOCAL_BIN)" go install github.com/pressly/goose/v3/cmd/goose@v3.14.0
	GOBIN="$(LOCAL_BIN)" go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.27.3
	GOBIN="$(LOCAL_BIN)" go install github.com/envoyproxy/protoc-gen-validate@v0.10.1
	GOBIN="$(LOCAL_BIN)" go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.27.3
	GOBIN="$(LOCAL_BIN)" go install github.com/rakyll/statik@v0.1.7

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

generate:
	make vendor-proto
	mkdir -p pkg/swagger
	make generate-auth-api
	make generate-access-api
	"$(LOCAL_BIN)/statik" -src=pkg/swagger/ -include='*.css,*.html,*.js,*.json,*.png'

generate-auth-api:
	mkdir -p pkg/auth_v1 pkg/swagger
	protoc --proto_path api/auth_v1 --proto_path vendor.protogen \
	--go_out=pkg/auth_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go="$(LOCAL_BIN)/protoc-gen-go" \
	--go-grpc_out=pkg/auth_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc="$(LOCAL_BIN)/protoc-gen-go-grpc" \
	--grpc-gateway_out=pkg/auth_v1 --grpc-gateway_opt=paths=source_relative \
	--plugin=protoc-gen-grpc-gateway="$(LOCAL_BIN)/protoc-gen-grpc-gateway" \
	--validate_out=lang=go:pkg/auth_v1 --validate_opt=paths=source_relative \
	--plugin=protoc-gen-validate="$(LOCAL_BIN)/protoc-gen-validate" \
	--openapiv2_out=allow_merge=true,merge_file_name=api:pkg/swagger \
	--plugin=protoc-gen-openapiv2="$(LOCAL_BIN)/protoc-gen-openapiv2" \
	api/auth_v1/auth.proto

generate-access-api:
	mkdir -p pkg/access_v1
	protoc --proto_path api/access_v1 --proto_path vendor.protogen \
	--go_out=pkg/access_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go="$(LOCAL_BIN)/protoc-gen-go" \
	--go-grpc_out=pkg/access_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc="$(LOCAL_BIN)/protoc-gen-go-grpc" \
	api/access_v1/access.proto

vendor-proto:
	@if [ ! -d vendor.protogen/google ]; then \
		git clone --depth 1 https://github.com/googleapis/googleapis vendor.protogen/googleapis && \
		mkdir -p vendor.protogen/google && \
		mv vendor.protogen/googleapis/google/api vendor.protogen/google && \
		rm -rf vendor.protogen/googleapis ; \
	fi
	@if [ ! -d vendor.protogen/validate ]; then \
		git clone --depth 1 --branch v0.10.1 https://github.com/envoyproxy/protoc-gen-validate vendor.protogen/protoc-gen-validate && \
		mkdir -p vendor.protogen/validate && \
		mv vendor.protogen/protoc-gen-validate/validate/*.proto vendor.protogen/validate && \
		rm -rf vendor.protogen/protoc-gen-validate ; \
	fi
	@if [ ! -d vendor.protogen/protoc-gen-openapiv2 ]; then \
		git clone --depth 1 --branch v2.27.3 https://github.com/grpc-ecosystem/grpc-gateway vendor.protogen/openapiv2 && \
		mkdir -p vendor.protogen/protoc-gen-openapiv2/options && \
		mv vendor.protogen/openapiv2/protoc-gen-openapiv2/options/*.proto vendor.protogen/protoc-gen-openapiv2/options && \
		rm -rf vendor.protogen/openapiv2 ; \
	fi

local-up:
	docker compose --env-file local.env up -d --build

local-down:
	docker compose --env-file local.env down

local-migration-status:
	"$(LOCAL_BIN)/goose" -dir "$(LOCAL_MIGRATION_DIR)" postgres "$(LOCAL_MIGRATION_DSN)" status -v

local-migration-up:
	"$(LOCAL_BIN)/goose" -dir "$(LOCAL_MIGRATION_DIR)" postgres "$(LOCAL_MIGRATION_DSN)" up -v

local-migration-down:
	"$(LOCAL_BIN)/goose" -dir "$(LOCAL_MIGRATION_DIR)" postgres "$(LOCAL_MIGRATION_DSN)" down -v

.PHONY: test
test:
	go clean -testcache
	go test ./... -covermode count -coverpkg=github.com/ivankornilov/auth/internal/service/...,github.com/ivankornilov/auth/internal/api/user -count 5

.PHONY: test-coverage
test-coverage:
	go clean -testcache
	go test ./... -coverprofile=coverage.tmp.out -covermode count -coverpkg=github.com/ivankornilov/auth/internal/service/...,github.com/ivankornilov/auth/internal/api/user -count 5
	grep -v 'mocks\|config' coverage.tmp.out > coverage.out
	rm coverage.tmp.out
	go tool cover -html=coverage.out
	go tool cover -func=./coverage.out | grep "total"
	grep -sqFx "/coverage.out" .gitignore || echo "/coverage.out" >> .gitignore
