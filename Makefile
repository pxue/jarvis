.PHONY: help run test coverage docs build dist clean tools dist-tools vendor-list vendor-update

TEST_FLAGS ?=
DATABASE_URL := "postgres://localhost/jarvis"
APP_PORT := 5331

all:
	@echo "****************************"
	@echo "** Jarvis build tool **"
	@echo "****************************"
	@echo "make <cmd>"
	@echo ""
	@echo "commands:"
	@echo "  run                   - run API in dev mode"
	@echo "  build                 - build api into bin/ directory"
	@echo "  tools                 - go get's a bunch of tools for dev"c:w
	@echo ""
	@echo ""

print-%: ; @echo $*=$($*)

##
## Tools
##
tools:
	go get -u github.com/pressly/sup/cmd/sup
	go get -u github.com/pressly/fresh

run:
	@(export DATABASE_URL=${DATABASE_URL}; export APP_PORT=${APP_PORT}; fresh -c runner.conf -p ./cmd/jarvis)

tunnel:
	ngrok http ${APP_PORT}


build:
	@mkdir -p ./bin
	GOGC=off go build -i -o ./bin/jarvis ./cmd/jarvis/main.go

