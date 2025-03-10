default: run

# Install dependencies
install:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest && \
	go install gotest.tools/gotestsum@latest && \
	go install github.com/rubenv/sql-migrate/...@latest && \
	go mod tidy && \
	go mod vendor

# Upgrade packages
upgrade-pkgs:
	go get -u ./... && make install

# Run
run:
	docker compose -f './docker-compose.local.yml' up -d --build 'db' && \
	cd scripts/db/ && ./setup_db.sh && \
	cd ../../ && LOGGER_FORMAT=cli BUNDEBUG=2 go run cmd/main.go

test-e2e:
	gotestsum --format pkgname ./test/...

# Build
build:
	go build cmd/main.go

lint:
	golangci-lint run ./...
