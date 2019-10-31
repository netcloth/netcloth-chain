#!/usr/bin/make -f

PACKAGES_SIMTEST=$(shell go list ./... | grep '/simulation')
VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')

export GO111MODULE = on

# process build tags

build_tags = netgo

update-swagger-docs:
	statik -src=client/lcd/swagger-ui -dest=client/lcd -f -m
.PHONY: update-swagger-docs


ifeq ($(WITH_CLEVELDB),yes)
  build_tags += gcc
endif
build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace :=
whitespace += $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

# process linker flags

ldflags = -X github.com/NetCloth/netcloth-chain/version.Name=nch \
		  -X github.com/NetCloth/netcloth-chain/version.ServerName=nchd \
		  -X github.com/NetCloth/netcloth-chain/version.ClientName=nchcli \
		  -X github.com/NetCloth/netcloth-chain/version.Version=$(VERSION) \
		  -X github.com/NetCloth/netcloth-chain/version.Commit=$(COMMIT) \
		  -X "github.com/NetCloth/netcloth-chain/version.BuildTags=$(build_tags_comma_sep)"

ifeq ($(WITH_CLEVELDB),yes)
  ldflags += -X github.com/NetCloth/netcloth-chain/types.DBBackend=cleveldb
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'

all: install

build: go.sum
ifeq ($(OS),Windows_NT)
	go build $(BUILD_FLAGS) -o build/nchd.exe ./cmd/nchd
	go build $(BUILD_FLAGS) -o build/nchcli.exe ./cmd/nchcli
else
	go build $(BUILD_FLAGS) -o build/nchd ./cmd/nchd
	go build $(BUILD_FLAGS) -o build/nchcli ./cmd/nchcli
endif

build-linux: go.sum
	GOOS=linux GOARCH=amd64 $(MAKE) build

install: go.sum update-swagger-docs
	make update-swagger-docs
	go install $(BUILD_FLAGS) ./cmd/nchd
	go install $(BUILD_FLAGS) ./cmd/nchcli


########################################
### Tools & dependencies

go-mod-cache: go.sum
	@echo "--> Download go modules to local cache"
	@go mod download

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
#	@go mod verify

draw-deps:
	@# requires brew install graphviz or apt-get install graphviz
	go get github.com/RobotsAndPencils/goviz
	@goviz -i ./cmd/nchd -d 2 | dot -Tpng -o dependency-graph.png

clean:
	rm -rf  build/
