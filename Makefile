include .env

GIT = /usr/bin/git
BUILD = $(/usr/bin/git describe)
LDFLAGS=-ldflags "-X build.Build=$(BUILD)"

.PHONY: db-schema-commit

db-schema-commit:
	@$(SQLITE) $(DBFILE) ".schema --indent" > $(SCHEMAFILE)

.PHONY: cli

cli:
	@echo $(BINOUT_AMD_DARWIN)
	@mkdir -p $(BINOUT_AMD_DARWIN) 
	@go build -o $(BINOUT_AMD_DARWIN) $(CLI_SRC)

.PHONY: api api-amd-linux api-arm-linux

api:
	@go build -o $(BINOUT_AMD_DARWIN) $(API_SRC)

api-amd-linux:
	@CGO_ENABLED=1 GOARCH=amd64 GOOS=linux CC="zig cc -target x86_64-linux" go build -o $(BINOUT_AMD_LINUX) $(API_DIR)

api-arm-linux:
	@CGO_ENABLED=1 GOOS=linux GOARCH=arm64 CC="zig cc -target aarch64-linux" go build -o $(BINOUT_ARM_LINUX) $(API_DIR)

.PHONY: ui ui-arm-linux ui-dist

ui:
	@go build -v $(LDFLAGS) -o $(BINOUT) $(UI_DIR)

ui-arm-linux:
	@GOARCH=arm64 GOOS=linux go build -o $(BINOUT_ARM_LINUX) $(UI_DIR) 

ui-dist: ui ui-arm-linux
	@tar cfp