ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

GIT_TAG ?= dirty-tag
GIT_VERSION ?= $(shell git describe --tags --dirty --always --abbrev=4 | sed -e 's/^v//')
GIT_HASH ?= $(shell git rev-parse HEAD)
DATE_FMT = +'%Y-%m-%dT%H:%M:%SZ'
SOURCE_DATE_EPOCH ?= $(shell git log -1 --pretty=%ct)

ifdef SOURCE_DATE_EPOCH
    BUILD_DATE ?= $(shell date -u -d "@$(SOURCE_DATE_EPOCH)" "$(DATE_FMT)" 2>/dev/null || date -u -r "$(SOURCE_DATE_EPOCH)" "$(DATE_FMT)" 2>/dev/null || date -u "$(DATE_FMT)")
else
    BUILD_DATE ?= $(shell date "$(DATE_FMT)")
endif
GIT_TREESTATE = "clean"
DIFF = $(shell git diff --quiet >/dev/null 2>&1; if [ $$? -eq 1 ]; then echo "1"; fi)
ifeq ($(DIFF), 1)
    GIT_TREESTATE = "dirty"
endif

EXE_PATH=btcd-node-handshake
PKG=github.com/A-Roberto-Company/btcd-node-handshake

LDFLAGS="-X $(PKG).GitVersion=$(GIT_VERSION) -X $(PKG).gitCommit=$(GIT_HASH) -X $(PKG).gitTreeState=$(GIT_TREESTATE) -X $(PKG).buildDate=$(BUILD_DATE)"

.PHONY: all build clean

all: clean build

BUILD_FLAGS=-ldflags $(LDFLAGS)

build: main.go log.go
	CGO_ENABLED=0 go build $(BUILD_FLAGS) -o $(EXE_PATH) .

install:
	CGO_ENABLED=0 go install $(BUILD_FLAGS) ./...

.PHONY: test

test:
	go test -race -cover -coverprofile cp.out -count=1 -timeout=30s ./...

clean:
	rm -f $(EXE_PATH)