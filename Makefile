OPENAPI_GENERATOR := java -jar ~/openapi-generator-cli.jar
API_SRC := ./docs/api.yaml
API_BUNDLED := ./docs/api-bundled.yaml
DOCS_OUTPUT_DIR := ./docs/web
DOCS_INTERNAL_DIR := ./docs/web/docs
RESOURCES_DIR := ./pkg/resources

generate-models:
	test -d $(RESOURCES_DIR) || mkdir -p $(RESOURCES_DIR)
	test -d $(dir $(API_SRC)) || mkdir -p $(dir $(API_SRC))
	test -d $(dir $(API_BUNDLED)) || mkdir -p $(dir $(API_BUNDLED))
	test -d $(DOCS_OUTPUT_DIR) || mkdir -p $(DOCS_OUTPUT_DIR)

	rm -rf $(DOCS_INTERNAL_DIR) && mkdir -p $(DOCS_INTERNAL_DIR)
	rm -rf $(RESOURCES_DIR) && mkdir -p $(RESOURCES_DIR)
	swagger-cli bundle $(API_SRC) --outfile $(API_BUNDLED) --type yaml

	$(OPENAPI_GENERATOR) generate \
		-i $(API_BUNDLED) -g go \
		-o $(DOCS_OUTPUT_DIR) \
		--additional-properties=packageName=resources \
		--import-mappings uuid.UUID=github.com/google/uuid --type-mappings string+uuid=uuid.UUID

	mkdir -p $(RESOURCES_DIR)
	find $(DOCS_OUTPUT_DIR) -name '*.go' -exec mv {} $(RESOURCES_DIR)/ \;
	find $(RESOURCES_DIR) -type f -name "*_test.go" -delete

build:
	go build -o ./cmd/profiles-svc/main ./cmd/profiles-svc/main.go

migrate-up:
	go build -o ./cmd/profiles-svc/main ./cmd/profiles-svc/main.go
	set -a && . ./.env && set +a && ./cmd/profiles-svc/main migrate up

migrate-down:
	go build -o ./cmd/profiles-svc/main ./cmd/profiles-svc/main.go
	set -a && . ./.env && set +a && ./cmd/profiles-svc/main migrate down

run-server:
	go build -o ./cmd/profiles-svc/main ./cmd/profiles-svc/main.go
	set -a && . ./.env && set +a && ./cmd/profiles-svc/main run service

docker-uo:
	docker compose up -d

docker-down:
	docker compose down

docker-rebuild:
	docker compose up -d --build --force-recreate