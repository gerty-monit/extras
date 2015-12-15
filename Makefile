.PHONY: install build test run clean-ts

build:
	@go build -v ./...

install:
	go get github.com/tools/godep
	godep restore

test:
	@go test ./...

test-cover:
	@echo "mode: set" > acc.coverage-out
	@go test -coverprofile=services.coverage-out ./services
	@cat services.coverage-out | grep -v "mode: set" >> acc.coverage-out
	@go tool cover -html=acc.coverage-out
	@rm *.coverage-out