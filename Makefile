.PHONY: checks clean package

VERSION ?= $(shell git rev-parse HEAD)

all: package
checks:
	$(MAKE) -C checks

clean:
	rm -vrf bin/*

bin:
	mkdir -p bin

package: sensu-atc-assets-$(VERSION).tar.gz
sensu-atc-assets-$(VERSION).tar.gz: bin checks
	sudo setcap 'cap_setuid,cap_setgid=eip' bin/check-systemd-unit
	tar --xattrs --xattrs-include='*' -czvf $@ bin/
