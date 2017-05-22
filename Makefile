.PHONY: build
build:
	go build -ldflags "-X main.appVersion=$$(git describe --tags)" .

.PHONY: deps
deps:
	go get -v github.com/Masterminds/glide
	glide install
