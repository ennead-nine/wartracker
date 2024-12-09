include .env

.PHONY: db-schema-commit

db-schema-commit:

	@$(SQLITE) $(DBFILE) ".schema --indent" > $(SCHEMAFILE)

.PHONY: cli

cli:

	@go build -o $(BINOUT) $(CLI_DIR)
