#!/usr/bin/make -f

export VERSION := $(shell echo $(shell git describe --always --match "v*") | sed 's/^v//')
export TMVERSION := $(shell go list -m github.com/cometbft/cometbft | sed 's:.* ::')
export COMMIT := $(shell git log -1 --format='%H')

BUILDDIR ?= $(CURDIR)/build
GOBIN ?= $(GOPATH)/bin
NILLION_BINARY ?= nilliond
DOCKER := $(shell which docker)


# Default target executed when no arguments are given to make.
default_target: all

.PHONY: default_target

# process build tags
build_tags = netgo
ifeq ($(LEDGER_ENABLED),true)
  ifeq ($(OS),Windows_NT)
    GCCEXE = $(shell where gcc.exe 2> NUL)
    ifeq ($(GCCEXE),)
      $(error gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false)
    else
      build_tags += ledger
    endif
  else
    UNAME_S = $(shell uname -s)
    ifeq ($(UNAME_S),OpenBSD)
      $(warning OpenBSD detected, disabling ledger support (https://github.com/cosmos/cosmos-sdk/issues/1988))
    else
      GCC = $(shell command -v gcc 2> /dev/null)
      ifeq ($(GCC),)
        $(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
      else
        build_tags += ledger
      endif
    endif
  endif
endif

# Set Linux/arm64 compiler.
UNAME_S = $(shell uname -s)
ifeq ($(UNAME_S),Linux)
  UNAME_M = $(shell uname -m)
  ifneq ($(UNAME_M),aarch64)
    LINUX_ARM64_CC = "aarch64-linux-gnu-gcc"
  else
    LINUX_ARM64_CC = $(CC)
  endif
else
  LINUX_ARM64_CC = $(CC)
endif

ifeq (cleveldb,$(findstring cleveldb,$(COSMOS_BUILD_OPTIONS)))
  build_tags += gcc
endif
build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=nillion \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=nilliond \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)" \
		  -X github.com/cometbft/cometbft/version.TMCoreSemVer=$(TMVERSION)

ifeq (,$(findstring nostrip,$(COSMOS_BUILD_OPTIONS)))
  ldflags += -w -s
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

# Build flags for compiling Go binaries
BUILD_FLAGS = -tags "netgo" -ldflags '$(ldflags)'

###############################################################################
###                                  Build                                  ###
###############################################################################

BUILD_TARGETS := build install

build: BUILD_ARGS=-o $(BUILDDIR)/
$(BUILD_TARGETS): go.sum $(BUILDDIR)/
	CGO_ENABLED="1" go $@ $(BUILD_FLAGS) $(BUILD_ARGS) ./...

build-cross: go.sum $(BUILDDIR)/
build-cross: build-darwin-amd64 build-darwin-arm64
build-cross: build-linux-amd64 build-linux-arm64

build-darwin-amd64 build-darwin-arm64: build-darwin-%:
	mkdir -p $(BUILDDIR)/darwin/$*
	GOOS=darwin GOARCH=$* go build $(BUILD_FLAGS) -o $(BUILDDIR)/darwin/$* ./...

build-linux-amd64:
	mkdir -p $(BUILDDIR)/linux/amd64
	CGO_ENABLED="1" GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BUILDDIR)/linux/amd64 ./...

build-linux-arm64:
	mkdir -p $(BUILDDIR)/linux/arm64
	test -z $(LINUX_ARM64_CC) || command -v $(LINUX_ARM64_CC) >/dev/null
	CC=$(LINUX_ARM64_CC) CGO_ENABLED="1" GOOS=linux GOARCH=arm64 go build $(BUILD_FLAGS) -o $(BUILDDIR)/linux/arm64 ./...

$(BUILDDIR)/:
	mkdir -p $(BUILDDIR)/

clean:
	rm -rf \
    $(BUILDDIR)/ \
    artifacts/ \
    tmp-swagger-gen/

all: build

build-all: tools build lint test vulncheck

.PHONY: distclean clean build-all


###############################################################################
###                                  Run                                  ###
###############################################################################

init:
	cd scripts && sh init.sh

###############################################################################
###                                Protobuf                                 ###
###############################################################################

protoVer=latest
protoImageName=cosmossdk-proto:$(protoVer)
protoImage=$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace --user 0 $(protoImageName)

# ------
# NOTE: If you are experiencing problems running these commands, try deleting
#       the docker images and execute the desired command again.
#
proto-all: proto-format proto-lint proto-gen proto-swagger-gen

proto-gen:
	@echo "Generating Protobuf files"
	$(protoImage) sh ./scripts/protocgen.sh

proto-swagger-gen:
	@echo "Downloading Protobuf dependencies"
	@make proto-download-deps
	@echo "Generating Protobuf Swagger"
	$(protoImage) sh ./scripts/protoc-swagger-gen.sh

proto-format:
	@echo "Formatting Protobuf files"
	$(protoImage) find ./ -name *.proto -exec clang-format -i {} \;

proto-lint:
	@echo "Linting Protobuf files"
	@$(protoImage) buf lint --error-format=json

proto-check-breaking:
	@echo "Checking Protobuf files for breaking changes"
	$(protoImage) buf breaking --against $(HTTPS_GIT)#branch=main

SWAGGER_DIR=./swagger-proto
THIRD_PARTY_DIR=$(SWAGGER_DIR)/third_party

proto-download-deps:
	mkdir -p "$(THIRD_PARTY_DIR)/cosmos_tmp" && \
	cd "$(THIRD_PARTY_DIR)/cosmos_tmp" && \
	git init && \
	git remote add origin "https://github.com/evmos/cosmos-sdk.git" && \
	git config core.sparseCheckout true && \
	printf "proto\nthird_party\n" > .git/info/sparse-checkout && \
	git pull origin "$(DEPS_COSMOS_SDK_VERSION)" && \
	rm -f ./proto/buf.* && \
	mv ./proto/* ..
	rm -rf "$(THIRD_PARTY_DIR)/cosmos_tmp"

	mkdir -p "$(THIRD_PARTY_DIR)/ibc_tmp" && \
	cd "$(THIRD_PARTY_DIR)/ibc_tmp" && \
	git init && \
	git remote add origin "https://github.com/cosmos/ibc-go.git" && \
	git config core.sparseCheckout true && \
	printf "proto\n" > .git/info/sparse-checkout && \
	git pull origin "$(DEPS_IBC_GO_VERSION)" && \
	rm -f ./proto/buf.* && \
	mv ./proto/* ..
	rm -rf "$(THIRD_PARTY_DIR)/ibc_tmp"

	mkdir -p "$(THIRD_PARTY_DIR)/cosmos_proto_tmp" && \
	cd "$(THIRD_PARTY_DIR)/cosmos_proto_tmp" && \
	git init && \
	git remote add origin "https://github.com/cosmos/cosmos-proto.git" && \
	git config core.sparseCheckout true && \
	printf "proto\n" > .git/info/sparse-checkout && \
	git pull origin "$(DEPS_COSMOS_PROTO_VERSION)" && \
	rm -f ./proto/buf.* && \
	mv ./proto/* ..
	rm -rf "$(THIRD_PARTY_DIR)/cosmos_proto_tmp"

	mkdir -p "$(THIRD_PARTY_DIR)/gogoproto" && \
	curl -SSL "https://raw.githubusercontent.com/cosmos/gogoproto/$(DEPS_COSMOS_GOGOPROTO)/gogoproto/gogo.proto" > "$(THIRD_PARTY_DIR)/gogoproto/gogo.proto"

	mkdir -p "$(THIRD_PARTY_DIR)/google/api" && \
	curl -sSL https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto > "$(THIRD_PARTY_DIR)/google/api/annotations.proto"
	curl -sSL https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto > "$(THIRD_PARTY_DIR)/google/api/http.proto"

	mkdir -p "$(THIRD_PARTY_DIR)/cosmos/ics23/v1" && \
	curl -sSL "https://raw.githubusercontent.com/cosmos/ics23/$(DEPS_COSMOS_ICS23)/proto/cosmos/ics23/v1/proofs.proto" > "$(THIRD_PARTY_DIR)/cosmos/ics23/v1/proofs.proto"


.PHONY: proto-all proto-gen proto-format proto-lint proto-check-breaking proto-swagger-gen

