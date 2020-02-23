#!/usr/bin/make -f

PACKAGES_SIMTEST=$(shell go list ./... | grep '/simulation')
VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')

export GO111MODULE = on

# process build tags

build_tags = netgo

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

ldflags = -X github.com/netcloth/netcloth-chain/version.Name=nch \
		  -X github.com/netcloth/netcloth-chain/version.ServerName=nchd \
		  -X github.com/netcloth/netcloth-chain/version.ClientName=nchcli \
		  -X github.com/netcloth/netcloth-chain/version.Version=$(VERSION) \
		  -X github.com/netcloth/netcloth-chain/version.Commit=$(COMMIT) \
		  -X "github.com/netcloth/netcloth-chain/version.BuildTags=$(build_tags_comma_sep)"

ifeq ($(WITH_CLEVELDB),yes)
  ldflags += -X github.com/netcloth/netcloth-chain/types.DBBackend=cleveldb
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'

all: get_tools install

get_tools:
	cd scripts && $(MAKE) get_tools

check_dev_tools:
	cd scripts && $(MAKE) check_dev_tools

get_dev_tools:
	cd scripts && $(MAKE) get_dev_tools

### Generate swagger docs for nchd
# update_nchd_swagger_docs:
#     @statik -src=lite/swagger-ui -dest=lite -f

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

install: go.sum
	go install $(BUILD_FLAGS) ./cmd/nchd
	go install $(BUILD_FLAGS) ./cmd/nchcli


########################################
### Tools & dependencies

go-mod-cache: go.sum
	@echo "--> Download go modules to local cache"
	@go mod download

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify

draw-deps:
	@# requires brew install graphviz or apt-get install graphviz
	go get github.com/RobotsAndPencils/goviz
	@goviz -i ./cmd/nchd -d 2 | dot -Tpng -o dependency-graph.png

clean:
	rm -rf  build/
