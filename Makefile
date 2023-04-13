TEST_DIR := test
TOOLS_DIR := tools


# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

all: build

fmt: ## Run go fmt against code.
	go fmt ./...

vet: ## Run go vet against code.
	go vet ./...

build: fmt vet ## Build manager binary.
	go build -o bin/astermule cmd/main.go

run: fmt vet ## Run code from your host.
	go run cmd/main.go

test:
	go test ./... -coverprofile cover.out

STAGING_REGISTRY ?= pluviophile225
IMAGE_NAME ?= astermule_param
TAG ?= latest

IMG ?= ${STAGING_REGISTRY}/${IMAGE_NAME}:${TAG}
docker-build:
	docker buildx build -t ${IMG} . --load

docker-push:
	docker push ${IMG} 

# go-get-tool will 'go get' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Downloading $(2)" ;\
GOBIN=$(PROJECT_DIR)/bin go install $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef

install-golint: ## check golint if not exist install golint tools
ifeq (, $(shell which golangci-lint))
	@{ \
	set -e ;\
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.1 ;\
	}
GOLINT_BIN=$(shell go env GOPATH)/bin/golangci-lint
else
GOLINT_BIN=$(shell which golangci-lint)
endif

lint: install-golint ## Run go lint against code.
	$(GOLINT_BIN) run -v

# for debug
debug:
	docker pull kasterism/test_a
	docker pull kasterism/test_b
	docker pull kasterism/test_c
	docker pull kasterism/test_d
	docker run --name test_a -p 8000:8000 -itd kasterism/test_a
	docker run --name test_b -p 8001:8001 -itd kasterism/test_b
	docker run --name test_c -p 8002:8002 -itd kasterism/test_c
	docker run --name test_d -p 8003:8003 -itd kasterism/test_d
	go run cmd/main.go --dag '{"nodes":[{"name":"A","action":"GET","url":"http://localhost:8000/test"},{"name":"B","action":"POST","url":"http://localhost:8001/test","dependencies":["A"]},{"name":"C","action":"POST","url":"http://localhost:8002/test","dependencies":["A"]},{"name":"D","action":"POST","url":"http://localhost:8003/test","dependencies":["B","C"]}]}'

clean:
	docker rm -f test_a
	docker rm -f test_b
	docker rm -f test_c
	docker rm -f test_d

testbed-build:
	$(MAKE) -C $(TOOLS_DIR)/testbed docker-build

testbed-push:
	$(MAKE) -C $(TOOLS_DIR)/testbed docker-push