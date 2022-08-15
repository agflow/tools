export # export all `make` variables as environment variables.

CI="agflow/ci:v0.8.15"

GOCI_VERSION = v1.48.0
GOCI = nogit/goci/$(GOCI_VERSION)/golangci-lint

test:
	@go test -v ./...

.PHONY: lint
lint: lint-backend

###############################
# dev tools

.PHONY: format-sql
format-sql:
	@build/run.sh ${CI} build/format/sql.sh

.PHONY: lint-sql
lint-sql:
	@build/run.sh ${CI} build/lint/sql.sh

.PHONY: lint-backend
lint-backend: $(GOCI)
	$(GOCI) run ./...

$(GOCI):
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(dir $(GOCI)) $(GOCI_VERSION)
