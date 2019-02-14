ghsshauth:
	go build

install: ghsshauth
	mkdir -p /usr/local/sbin
	cp ghsshauth /usr/local/sbin/ghsshauth
.PHONY: install
