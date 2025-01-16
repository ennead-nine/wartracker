include .env

.PHONY: db-schema-commit

db-schema-commit:
	@$(SQLITE) $(DBFILE) ".schema --indent" > $(SCHEMAFILE)

.PHONY: cli

cli:
	@CGO_CFLAGS_ALLOW="-Xpreprocessor" go build -o $(BINOUT) $(CLI_DIR)

.PHONY: api api-amd-linux api-arm-linux

api:
	@go build -o $(BINOUT) $(API_DIR)

api-amd-linux:
	@CGO_ENABLED=1 GOARCH=amd64 GOOS=linux CC="zig cc -target x86_64-linux" go build -o $(BINOUT_AMD_LINUX) $(API_DIR)

api-arm-linux:
	@CGO_ENABLED=1 GOOS=linux GOARCH=arm64 CC="zig cc -target aarch64-linux" go build -o $(BINOUT_ARM_LINUX) $(API_DIR)

.PHONY: ui ui-arm-linux

ui:
	@go build -o $(BINOUT) $(UI_DIR)

ui-arm-linux:
	@GOARCH=arm64 GOOS=linux go build -o $(BINOUT_ARM_LINUX) $(UI_DIR) 