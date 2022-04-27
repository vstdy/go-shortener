GO_BIN=$(GOPATH)/bin
PATH:=$(PATH):$(GO_BIN)

PROTO_IN_DIR=./api/grpc
PROTO_OUT_DIR=.
PROTO_DEPS_DIR=./api/grpc/deps

.PHONY: test
test:
	go test -v ./... --count=1

.PHONY: deps
deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.0
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.10.0
	go install github.com/envoyproxy/protoc-gen-validate@v0.6.7
	go install github.com/bufbuild/buf/cmd/buf@v1.3.1
	go install github.com/bufbuild/buf/cmd/protoc-gen-buf-breaking@v1.3.1
	go install github.com/bufbuild/buf/cmd/protoc-gen-buf-lint@v1.3.1

.PHONY: build-proto
build-proto:
	protoc \
		-I ${PROTO_IN_DIR} \
		-I ${PROTO_DEPS_DIR} \
		--go_out $(PROTO_OUT_DIR) \
		--go-grpc_out $(PROTO_OUT_DIR) \
		--grpc-gateway_opt logtostderr=true,allow_delete_body=true \
		--grpc-gateway_out $(PROTO_OUT_DIR) \
  	--validate_opt lang=go \
  	--validate_out $(PROTO_OUT_DIR) \
		url_service.proto

.PHONY: generate
generate:
	buf generate
