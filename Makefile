check_defined = \
	$(if !$(value $1),, \
		$(error Undefined $1$(if $2, ($2))))

dev_server:
	go run ./cmd/server/main.go

go_tidy:
	go mod tidy

generate_api: 
	cd ./internal/openapi && go generate

build_migrate:
	go build -o ./cmd/bin/migration ./cmd/migration/*.go

migration: $(call check_defined, name) build_migrate
	./cmd/bin/migration create $(name)

migrate: build_migrate
	./cmd/bin/migration migrate

rollback: build_migrate
	./cmd/bin/migration rollback