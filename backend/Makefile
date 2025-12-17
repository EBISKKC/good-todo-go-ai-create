.PHONY: run stop init-db dev generate_ent mockgen oapi-gen generate_di migrate_diff migrate_apply migrate_status migrate_down test_unit test_integration test fmt lint vet

# Docker
run:
	docker compose up -d

stop:
	docker compose down

init-db:
	docker compose down -v
	docker compose up -d postgres
	sleep 3
	$(MAKE) migrate_apply

# 開発
dev:
	go run ./cmd/api/main.go

# コード生成
generate_ent:
	go generate ./internal/ent/...

mockgen:
	go generate ./internal/domain/repository/...

oapi-gen:
	oapi-codegen -config openapi/config-public.yaml openapi/openapi-public.yaml

generate_di:
	go run ./scripts/generate_di.go

# マイグレーション
migrate_diff:
	atlas migrate diff --env local

migrate_apply:
	atlas migrate apply --env local

migrate_status:
	atlas migrate status --env local

migrate_down:
	atlas migrate down --env local $(if $(n),--steps $(n),--steps 1)

# テスト
test_unit:
	go test -v -short ./internal/usecase/...

test_integration:
	go test -v ./internal/integration_test/...

test:
	go test -v ./...

# コード品質
fmt:
	go fmt ./...

lint:
	golangci-lint run

vet:
	go vet ./...
