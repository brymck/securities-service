PROTOS := alpha_vantage_api

SERVICE_NAME := $(notdir $(CURDIR))
GO_FILES := $(shell find . -name '*.go')
PROTO_FILES := $(shell find proto -name '*.proto' 2>/dev/null) $(foreach proto,$(PROTOS),proto/$(proto).proto)
PROTO_PATH := /usr/local/include
GENPROTO_FILES := $(patsubst proto/%.proto,genproto/%.pb.go,$(PROTO_FILES))

all: generate test build

init: .init.stamp

.init.stamp:
	go get -u github.com/golang/protobuf/protoc-gen-go
	go mod download
	touch $@

generate: $(GENPROTO_FILES)

proto genproto:
	mkdir $@

proto/alpha_vantage_api.proto: | proto
	curl --location --output $@ --silent https://raw.githubusercontent.com/brymck/alpha-vantage-service/master/$@

genproto/%.pb.go: proto/%.proto | .init.stamp genproto
	protoc -Iproto -I$(PROTO_PATH) --go_out=plugins=grpc:$(dir $@) $<

test: profile.out

profile.out: $(GO_FILES) $(GENPROTO_FILES) | .init.stamp
	go test -race -coverprofile=profile.out -covermode=atomic ./...

build: service

service: $(GO_FILES) $(GENPROTO_FILES) | .init.stamp
	go build -ldflags='-w -s' -o service cmd/web/*.go

run: service
	./service

docker:
	docker build . --tag gcr.io/$(shell gcloud config get-value project)/$(SERVICE_NAME)

clean:
	rm -rf proto/ genproto/ .init.stamp profile.out client service

.PHONY: all init generate test build run docker clean
