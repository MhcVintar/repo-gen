.PHONY: default
default:
	$(MAKE) test
	$(MAKE) build

.PHONY: test
test:
	go test -v ./...

.PHONY: build
build:
	go build -o bin/repogen cmd/repogen.go