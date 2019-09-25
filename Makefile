NAME = mgr
VERSION ?= v0.0.1
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

bin/$(NAME): build

build:
	CGO_ENABLED=0 go build -i -ldflags="-s -w -X main.version=$(VERSION)" -o bin/$(NAME) cmd/$(NAME)/*.go

dist: build
	mkdir -p dist
	tar cfz dist/$(NAME)-$(VERSION)_$(GOOS)_$(GOARCH).tar.gz -C bin/ $(NAME)

clean:
	rm -rf bin dist

redis-up:
	docker run -d \
		-p 6379:6379 \
		--name=$(NAME)-redis redis:4

redis-down:
	docker rm -f $(NAME)-redis
