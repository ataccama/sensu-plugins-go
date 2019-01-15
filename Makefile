.PHONY: checks clean package meta

VERSION ?= $(shell git rev-parse HEAD)

all: meta

checks:
	$(MAKE) -C checks

clean:
	rm -vrf bin/*

bin:
	mkdir -p bin

package: sensu-atc-assets-$(VERSION).tar.gz
sensu-atc-assets-$(VERSION).tar.gz: bin checks
	tar --xattrs --xattrs-include='*' -czvf $@ bin/
	aws --profile=devops_s3 s3 cp $@ s3://atc-devops/sensu/$@

meta: sensu-atc-assets-$(VERSION).tar.gz
	@echo
	@echo URL: http://atc-devops.s3.eu-central-1.amazonaws.com/sensu/$^
	@echo SUM: $(shell sha512sum $^ | cut -d\  -f1)
