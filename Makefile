#!/usr/bin/make -f

BUILDDIR ?= $(CURDIR)/build
GOBIN ?= $(GOPATH)/bin
NILLION_BINARY ?= nilliond

# Build flags for compiling Go binaries
BUILD_FLAGS = -tags "netgo" -ldflags "-w -s"

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