ARGS=

.PHONY: build
build:
	go build -ldflags "-X main.appVersion=$$(git describe --tags)" .

.PHONY: debug
debug: build
	./kokoro-go ${ARGS}

.PHONY: deps
deps:
	go get -v github.com/Masterminds/glide
	go get -v github.com/jteeuwen/go-bindata/...
	glide install
