.PHONY: build
build:
	go build -ldflags "-X main.appVersion=$$(git describe --tags)" .
