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
	tar czvf $@ bin/
