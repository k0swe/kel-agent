VERSION = $(shell < debian/changelog head -1 | egrep -o "[0-9]+\.[0-9]+\.[0-9]+")
GITCOMMIT = $(shell git rev-parse --short HEAD 2> /dev/null || true)

GENERATED = kel-agent kel-agent_*.pkg win/kel-agent_*.msi win/kel-agent.wixobj autorevision.cache ../kel-agent_* ../*.deb

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
	cd .. && tar -cvJf kel-agent_$(VERSION).orig.tar.xz --exclude-vcs kel-agent

.PHONY: deb-orig-tarball
deb-orig-tarball: autorevision.cache
	cd .. && tar -cvJf kel-agent_$(VERSION).orig.tar.xz --exclude-vcs --exclude=debian --exclude=.github --exclude=.idea kel-agent

# TODO: This target can be removed once the package is in Debian stable and Ubuntu stable, 2021-05
wsjtx-go.deb:
	cd .. && \
	wget http://ftp.debian.org/debian/pool/main/g/golang-github-k0swe-wsjtx-go/golang-github-k0swe-wsjtx-go-dev_1.1.0-2_all.deb

# TODO: This target can be removed once the package is in Debian stable and Ubuntu stable, 2021-05
leemcloughlin-jdn.deb:
	cd .. && \
	wget http://ftp.debian.org/debian/pool/main/g/golang-github-leemcloughlin-jdn/golang-github-leemcloughlin-jdn-dev_0.0~git20201102.6f88db6-2_all.deb

.PHONY: deb-package
deb-package: deb-tarball wsjtx-go.deb leemcloughlin-jdn.deb
	# https://wiki.debian.org/sbuild
	sbuild -d stable \
      --extra-package=../golang-github-k0swe-wsjtx-go-dev_1.1.0-2_all.deb \
      --extra-package=../golang-github-leemcloughlin-jdn-dev_0.0~git20201102.6f88db6-2_all.deb

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
