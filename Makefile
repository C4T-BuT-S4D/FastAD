.PHONY: lint-proto
lint-proto:
	cd proto && buf lint

.PHONY: lint-go
lint-go:
	golangci-lint run -v -c .golangci.yml ./...

.PHONY: lint-front
lint-front:
	cd front && yarn lint

.PHONY: goimports
goimports:
	gofancyimports fix --local github.com/c4t-but-s4d/fastad -w $(shell find . -type f -name '*.go' -not -path "./pkg/proto/*")

.PHONY: proto
proto: lint-proto
	rm -rf pkg/proto
	cd proto && buf generate && buf build --as-file-descriptor-set -o descriptor.pb

.PHONY: tidy
tidy:
	go mod tidy -compat="1.18"

.PHONY: lint
lint: lint-proto lint-front lint-go

.PHONY: test-go
test-go:
	go test -race ./...

.PHONY: test
test: test-go

.PHONY: commit-proto
commit-proto: proto
	git add proto
	git add pkg/proto
	git add front/src/proto
	git commit -m "Regenerate proto"

.PHONY: reset-db
reset-db:
	docker compose exec postgres psql -U fastad fastad -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
	go run ./cmd/migrator/main.go init

.PHONY: migrate-db
migrate-db:
	go run ./cmd/migrator/main.go migrate

.PHONY: db
db: reset-db migrate-db
