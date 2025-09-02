VERSION ?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || cat $(CURDIR)/.version 2> /dev/null || echo v0)
BLDVER = module:$(MODULE),version:$(VERSION),build:$(CIRCLE_BUILD_NUM)
BASE = $(CURDIR)
MODULE = tucp

.PHONY: all $(MODULE) install
all: $(MODULE)

$(MODULE):| $(BASE)
	@GO111MODULE=on GOFLAGS=-mod=vendor go build -v -trimpath -o $(BASE)/bin/$@

$(BASE):
	@mkdir -p $(dir $@)

install:
	@GO111MODULE=on GOFLAGS=-mod=vendor go install -v
