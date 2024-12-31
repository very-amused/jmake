# Installation vars
ifndef PREFIX
PREFIX=/usr/local
endif
ifndef DATADIR
DATADIR=$(PREFIX)/share
endif

jmake:
	go build -o $@
.PHONY: jmake

install:
	install -d $(DESTDIR)$(PREFIX)/bin
	install jmake $(DESTDIR)$(PREFIX)/bin/jmake
	install -d $(DESTDIR)$(DATADIR)/doc/jmake
	install -m644 README.md $(DESTDIR)$(DATADIR)/doc/jmake/README.md
	@# TODO: license
.PHONY: install

uninstall:
	rm -f $(DESTDIR)$(PREFIX)/bin/jmake
	rm -rf $(DESTDIR)$(DATADIR)/doc/jmake
.PHONY: uninstall
