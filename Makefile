GOOS = $(shell go env GOOS)
GOARCH = $(shell go env GOARCH)
BUILD_DIR = dist/${GOOS}_${GOARCH}

ifeq ($(GOOS),windows)
OUTPUT_PATH = ${BUILD_DIR}/baton-jd-edwards.exe
else
OUTPUT_PATH = ${BUILD_DIR}/baton-jd-edwards
endif

# Set the build tag conditionally based on BATON_LAMBDA_SUPPORT
ifdef BATON_LAMBDA_SUPPORT
	BUILD_TAGS=-tags baton_lambda_support
else
	BUILD_TAGS=
endif

.PHONY: build
build:
	go build ${BUILD_TAGS} -o ${OUTPUT_PATH} ./cmd/baton-jd-edwards

.PHONY: update-deps
update-deps:
	go get -d -u ./...
	go mod tidy -v
	go mod vendor

.PHONY: add-dep
add-dep:
	go mod tidy -v
	go mod vendor

.PHONY: lint
lint:
	golangci-lint run
