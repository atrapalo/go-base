all: help

##
## Repository Docker Helper Commands
## Available commands:
##

.PHONY: help setup up stop pull down logs shell test build

help: Makefile
	@sed -n 's/^##//p' $<

## setup:       Setup
setup:
	go mod vendor && ln -sf ../../pre-commit .git/hooks/pre-commit

## unit:        Run unit tests
unit:
	go test ./... | grep -v 'no test files'
