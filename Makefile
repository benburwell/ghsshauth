SOURCES     := $(shell find . -name '*.go')
BINARY      := ghsshauth
CLI_PACKAGE := ghsshauthcli

$(BINARY): $(SOURCES)
	go build -o "$@" "./$(CLI_PACKAGE)"

install: $(BINARY)
	mkdir -p /usr/local/sbin
	cp "$(BINARY)" "/usr/local/sbin/$(BINARY)"
.PHONY: install
