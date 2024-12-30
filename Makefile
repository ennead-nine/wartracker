include .env

.PHONY: db-schema-commit

db-schema-commit:
	@$(SQLITE) $(DBFILE) ".schema --indent" > $(SCHEMAFILE)

.PHONY: cli

cli:
	@CGO_CFLAGS_ALLOW="-Xpreprocessor" go build -o $(BINOUT) $(CLI_DIR)

.PHONY: api

api:
	@go build -o $(BINOUT) $(API_DIR)
