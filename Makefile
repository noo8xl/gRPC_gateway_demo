run-build:
	export GO_ENV=production && \
	rm -f ./bin/gateway && \
	GOFLAGS=-mod=mod GOOS=linux GOARCH=amd64 && \
	go build -o bin/gateway cmd/main.go

run-dev:
	export GO_ENV=development && \
	go run cmd/main.go

run-test:
	export GO_ENV=test && \
	go test ./... -v

gen-doc-u:
	@echo "Generating comprehensive documentation for gateway_demo..."
	@mkdir -p docs
	@echo "# Anvil Gateway API Documentation" > docs/user-api.md
	@echo "\n## Generated on: `date`\n" >> docs/user-api.md
	@go doc -all ./routes/user/ >> docs/user-api.md
	@echo "Documentation generated at docs/user-api.md"

gen-doc-a:
	@echo "Generating comprehensive documentation for gateway_demo..."
	@mkdir -p docs
	@echo "# Anvil Gateway API Documentation" > docs/admin-api.md
	@echo "\n## Generated on: `date`\n" >> docs/admin-api.md
	@go doc -all ./routes/admin/ >> docs/admin-api.md
