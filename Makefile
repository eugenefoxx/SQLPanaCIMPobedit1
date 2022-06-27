.PHONY: build
build:
	go build -v -race ./cmd/panasap

.DEFAULT_GOAL := build
