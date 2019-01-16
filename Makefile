## ----- VARIABLES -----
## Go module name.
MODULE = $(shell basename "$$(pwd)")
ifeq ($(shell ls -1 go.mod 2> /dev/null),go.mod)
	MODULE = $(shell cat go.mod | grep module | awk '{print $$2}')
endif

## Program version.
VERSION = "none"
ifneq ($(shell git describe --tags 2> /dev/null),)
	VERSION = $(shell git describe --tags | cut -c 2-)
endif

## Custom Go linker flag.
LDFLAGS = -X $(MODULE)/internal/info.Version=$(VERSION)



## ----- TARGETS ------
## Generic:
.PHONY: default version setup install build clean run lint test review release \
        help

default: run
version: ## Show project version (derived from 'git describe').
	@echo $(VERSION)

setup: go-setup ## Set up this project on a new device.
	@echo "Configuring githooks..." && \
	 git config core.hooksPath .githooks && \
	 echo done

install: go-install ## Install project dependencies.
build: go-build ## Build project.
clean: go-clean ## Clean build artifacts.
run: go-run ## Run project (development).
lint: go-lint ## Lint and check code.
test: go-test ## Run tests.
review: go-review ## Lint code and run tests.
release: gr-release ## Release / deploy this project.

## Show usage for the targets in this Makefile.
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	   awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'


## CI:
.PHONY: ci-install ci-test ci-deploy
ci-install: go-install
ci-test: go-test
ci-deploy:
	@echo "No deployment procedure defined."


## Go:
.PHONY: go-deps go-bench go-setup go-install go-build go-clean go-run go-lint \
        go-test go-review

go-deps: ## Verify and tidy project dependencies.
	@echo "Verifying module dependencies..." && \
	 go mod verify && \
	 echo "Tidying module dependencies..." && \
	 go mod tidy && \
	 echo done

go-bench: ## Run benchmarks.
	@echo "Running benchmarks with 'go test -bench=.'..." && \
	 $(__TEST) -run=^$$ -bench=. -benchmem ./...

go-setup: go-install go-deps

go-install:
	@echo "Downloading module dependencies..." && \
	 go mod download && \
	 echo done

BUILDARGS = -ldflags "$(LDFLAGS)" $(BARGS)
BDIR = .
go-build:
	@echo "Building with 'go build'..." && \
	 go build $(BUILDARGS) $(BDIR) && \
	 echo done

go-clean:
	@echo "Cleaning with 'go clean'..." && \
	 go clean $(BDIR) && \
	 echo done

go-run:
	@echo "Running with 'go run'..." && \
	 go run $(BUILDARGS) $(BDIR)

go-lint:
	@if command -v goimports > /dev/null; then \
	   echo "Formatting code with 'goimports'..." && goimports -w .; \
	 else \
	   echo "'goimports' not installed, formatting code with 'go fmt'..." && \
	   go fmt .; \
	 fi && \
	 if command -v golint > /dev/null; then \
	   echo "Linting code with 'golint'..." && golint ./...; \
	 else \
	   echo "'golint' not installed, skipping linting step."; \
	 fi && \
	 echo "Checking code with 'go vet'..." && go vet ./... && \
	 echo done

COVERFILE = coverage.out
TIMEOUT   = 20s
TARGS = -race
__TEST = go test ./... -coverprofile="$(COVERFILE)" -covermode=atomic \
                       -timeout="$(TIMEOUT)" $(TARGS)
go-test:
	@echo "Running tests with 'go test':" && $(__TEST)

go-review: go-lint go-test


## Goreleaser:
.PHONY: gr-release gr-snapshot
__GR = MODULE="$(MODULE)" LDFLAGS="$(LDFLAGS)" goreleaser --rm-dist

gr-release:
	@echo "Releasing with 'goreleaser'..." && $(__GR)

gr-snapshot: ## Make a snapshot release.
	@echo "Making snapshot with 'goreleaser'..." && $(__GR) --snapshot
