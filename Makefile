VERSION = $(shell < debian/changelog head -1 | egrep -o "[0-9]+\.[0-9]+\.[0-9]+")
GITCOMMIT = $(shell git rev-parse --short HEAD 2> /dev/null || true)

GENERATED = kel-agent kel-agent_*.pkg win/kel-agent_*.msi win/kel-agent.wixobj autorevision.cache ../kel-agent_*

.PHONY: all
all: kel-agent

.PHONY: test
test:
	go test ./...
	go vet ./...

kel-agent: test
	export GITCOMMIT=$(GITCOMMIT) && scripts/build.sh

architecture.svg:
	# apt install graphviz
	dot -T svg -o architecture.svg < architecture.dot

autorevision.cache:
	autorevision -s VCS_SHORT_HASH -o ./autorevision.cache

.PHONY: deb-tarball
deb-tarball: autorevision.cache
	cd .. && tar -cvzf kel-agent_$(VERSION).orig.tar.gz kel-agent

.PHONY: deb-package
deb-package: deb-tarball
	# https://wiki.debian.org/sbuild
	sbuild

.PHONY: mac-package
mac-package: kel-agent
	# http://s.sudre.free.fr/Software/Packages/about.html
	packagesbuild macos/kel-agent.pkgproj
	mv kel-agent.pkg kel-agent_mac.pkg

.PHONY: win-package
win-package: kel-agent
	# https://wixtoolset.org/
	cd win && candle kel-agent.wxs && light kel-agent.wixobj

.PHONY: clean
clean:
	rm -f $(GENERATED)
