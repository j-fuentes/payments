ROOT_PKG=github.com/j-fuentes/payments
ROOT_PKG_DIR=${GOPATH}/src/$(ROOT_PKG)

all: build

build:
	cd $(ROOT_PKG_DIR) && go build .

test: build
	cd $(ROOT_PKG_DIR) && go test ./...

run:
	cd $(ROOT_PKG_DIR) && go run ./cmd/payments-server/main.go

generate:
	cd $(ROOT_PKG_DIR) && go generate ./...

vet:
	cd $(ROOT_PKG_DIR) && go vet ./...

lint: vet
	cd $(ROOT_PKG_DIR) && golint

validate: api/swagger.yml
	cd $(ROOT_PKG_DIR) && swagger validate ./api/swagger.yml
