all: help

##
## Repository Docker Helper Commands
## Available commands:
##

.PHONY: help vendor up stop pull down logs shell test build

help: Makefile
	@sed -n 's/^##//p' $<

## vendor:      Vendor
vendor:
	go mod tidy && go mod vendor && go mod tidy && ln -sf ../../pre-commit .git/hooks/pre-commit

## unit:        Run unit tests
unit:
	go test ./... | grep -v 'no test files'
