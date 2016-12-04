.PHONY: help run test coverage docs build dist clean tools dist-tools vendor-list vendor-update

TEST_FLAGS ?=

all:
	@echo "****************************"
	@echo "** Jarvis build tool **"
	@echo "****************************"
	@echo "make <cmd>"
	@echo ""
	@echo "commands:"
	@echo "  run                   - run API in dev mode"
	@echo "  build                 - build api into bin/ directory"
	@echo ""
	@echo ""

print-%: ; @echo $*=$($*)

run:
	@go run cmd/jarvis/main.go

build:
	@mkdir -p ./bin
	GOGC=off go build -i -o ./bin/jarvis ./cmd/jarvis/main.go

