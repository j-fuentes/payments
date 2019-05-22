ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

all: build

build:
	cd $(ROOT_DIR) && go build -o out/payments-service ./cmd/payments-service

serve:
	cd $(ROOT_DIR) && go run ./cmd/payments-service/main.go

test: build
	cd $(ROOT_DIR) && go test ./...

generate:
	cd $(ROOT_DIR) && go generate ./...

vet:
	cd $(ROOT_DIR) && go vet ./...

lint: vet
	cd $(ROOT_DIR) && golint

validate: api/swagger.yml
	cd $(ROOT_DIR) && swagger validate ./api/swagger.yml
