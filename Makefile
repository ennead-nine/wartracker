include .env

.PHONY: db-schema-commit

db-schema-commit:

	@$(SQLITE) $(DBFILE) ".schema --indent" > $(SCHEMAFILE)
