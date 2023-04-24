GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --abbrev=0)
INTERNAL_PROTO_FILES=$(shell find internal -name *.proto)
API_PROTO_FILES=$(shell find api -name *.proto)


.PHONY: init
# init env
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest

.PHONY: api
# generate api proto
api:
	protoc --proto_path=./api \
	       --proto_path=./third_party \
 	       --go_out=paths=source_relative:./api \
 	       --go-http_out=paths=source_relative:./api \
 	       --go-grpc_out=paths=source_relative:./api \
 	       --go-errors_out=paths=source_relative:./api \
 	       --validate_out=paths=source_relative,lang=go:./api \
	       $(API_PROTO_FILES)

.PHONY: swagger
# generate api proto
swagger:
	protoc --proto_path=. \
	       --proto_path=./third_party \
 	       --openapiv2_out . \
           --openapiv2_opt logtostderr=true \
           --openapiv2_opt json_names_for_fields=false \
           $(API_PROTO_FILES)
.PHONY: build
# build
build:
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...

.PHONY: image
# image
image:
	make all;
	for file in `ls bin`;\
	do \
	docker build --build-arg filename=$$file -t $$file:$(VERSION) .;\
	done

.PHONY: tidy
# generate
tidy:
	go mod tidy -compat=1.17


.PHONY: generate
# generate
generate:
	go get github.com/google/wire/cmd/wire@latest
	go generate ./...

.PHONY: wire
# wire
wire:
	go generate ./...

.PHONY: easyjson
# easyjson
easyjson:
	go get github.com/mailru/easyjson/easyjson

.PHONY: all
# generate all
all:
	make init;
	make api;
	make generate;
	make build;

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
