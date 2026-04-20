include versions.env
export ROOT_DIR     = $(shell git rev-parse --show-toplevel)
export GITCOMMIT    = $(shell git rev-parse --short HEAD 2>/dev/null || true)

OS            := $(shell uname -s | tr '[:upper:]' '[:lower:]')
RAW_ARCH      := $(shell uname -m)
ARCH          := $(if $(filter x86_64,$(RAW_ARCH)),amd64,$(if $(filter aarch64,$(RAW_ARCH)),arm64,$(if $(filter armv7l,$(RAW_ARCH)),armhf,$(RAW_ARCH))))
HAMLIB_PREFIX := $(ROOT_DIR)/out/hamlib/$(HAMLIB_VERSION)/$(OS)-$(ARCH)

# Use the local Hamlib prefix if it exists; otherwise fall through to system paths.
ifneq ($(wildcard $(HAMLIB_PREFIX)/lib/pkgconfig),)
export PKG_CONFIG_PATH = $(HAMLIB_PREFIX)/lib/pkgconfig
export PATH := $(HAMLIB_PREFIX)/bin:$(PATH)
endif

VERSION       := $(KEL_AGENT_VERSION)

GENERATED = kel-agent kel-agent.exe kel-agent_*.pkg win/kel-agent_*.msi win/kel-agent.wixobj \
  autorevision.cache ../kel-agent_* ../*.deb \
  flatpak/repo/ flatpak/.flatpak-builder/ flatpak/kel_agent.flatpak flatpak/flatpak_app/ flatpak/build-out/

# ---------------------------------------------------------------------------
# 1. Developer build
# ---------------------------------------------------------------------------

.PHONY: all
all: kel-agent

.PHONY: test
test:
	go test -tags hamlib ./...
	go vet  -tags hamlib ./...
	if command -v appstream-util >/dev/null; then appstream-util validate --nonet assets/radio.k0swe.Kel_Agent.metainfo.xml; fi
	if command -v desktop-file-validate >/dev/null; then desktop-file-validate assets/radio.k0swe.Kel_Agent.desktop; fi

.PHONY: test-nohamlib
test-nohamlib:
	go test ./...
	go vet  ./...

kel-agent: test
	export GITCOMMIT=$(GITCOMMIT) VERSION=v$(VERSION) && scripts/build.sh

# ---------------------------------------------------------------------------
# 2. Release build
# ---------------------------------------------------------------------------

.PHONY: hamlib
hamlib:
	scripts/build-hamlib.sh

.PHONY: release
release: hamlib test
	export GITCOMMIT=$(GITCOMMIT) VERSION=v$(VERSION) && scripts/build.sh

# ---------------------------------------------------------------------------
# 3. Packaging
# ---------------------------------------------------------------------------

assets/modules.txt:
	go mod vendor
	mv vendor/modules.txt assets/
	rm -rf vendor

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

# TODO: This target can be removed once the package is in Debian stable and Ubuntu stable
../golang-github-k0swe-wsjtx-go-dev_4.0.6-1_all.deb:
	wget https://github.com/k0swe/wsjtx-go/releases/download/v4.0.1/golang-github-k0swe-wsjtx-go-dev_4.0.6-1_all.deb \
	-O ../golang-github-k0swe-wsjtx-go-dev_4.0.6-1_all.deb

# TODO: This target can be removed once the package is in Debian stable and Ubuntu stable
../golang-github-mazznoer-csscolorparser-dev_0.1.3-1_all.deb:
	wget https://github.com/k0swe/wsjtx-go/releases/download/v4.0.1/golang-github-mazznoer-csscolorparser-dev_0.1.3-1_all.deb \
	-O ../golang-github-mazznoer-csscolorparser-dev_0.1.3-1_all.deb

# TODO: This target can be removed once the package is in Debian stable and Ubuntu stable
../golang-github-adrg-xdg-dev_0.4.0-1_all.deb:
	wget http://ftp.us.debian.org/debian/pool/main/g/golang-github-adrg-xdg/golang-github-adrg-xdg-dev_0.4.0-1_all.deb \
	-O ../golang-github-adrg-xdg-dev_0.4.0-1_all.deb

.PHONY: deb-package
deb-package: deb-tarball ../golang-github-k0swe-wsjtx-go-dev_4.0.6-1_all.deb ../golang-github-adrg-xdg-dev_0.4.0-1_all.deb ../golang-github-mazznoer-csscolorparser-dev_0.1.3-1_all.deb
	# https://wiki.debian.org/sbuild
	sbuild -d stable \
      --extra-package=../golang-github-k0swe-wsjtx-go-dev_4.0.6-1_all.deb \
      --extra-package=../golang-github-adrg-xdg-dev_0.4.0-1_all.deb \
      --extra-package=../golang-github-mazznoer-csscolorparser-dev_0.1.3-1_all.deb

.PHONY: deb-package-ci
deb-package-ci: deb-tarball ../golang-github-k0swe-wsjtx-go-dev_4.0.6-1_all.deb ../golang-github-adrg-xdg-dev_0.4.0-1_all.deb ../golang-github-mazznoer-csscolorparser-dev_0.1.3-1_all.deb
	@test "$$(id -u)" -eq 0 || (echo "deb-package-ci requires root access (designed for CI container environments)" && exit 1)
	apt-get update
	apt-get install -y --no-install-recommends \
      ../golang-github-k0swe-wsjtx-go-dev_4.0.6-1_all.deb \
      ../golang-github-adrg-xdg-dev_0.4.0-1_all.deb \
      ../golang-github-mazznoer-csscolorparser-dev_0.1.3-1_all.deb
	dpkg-buildpackage -b -us -uc

.PHONY: flatpak
flatpak: kel-agent
	cd flatpak && \
      flatpak-builder --force-clean build-out radio.k0swe.Kel_Agent.yml --repo=repo && \
      flatpak build-bundle repo kel_agent.flatpak radio.k0swe.Kel_Agent main

.PHONY: stage-hamlib
stage-hamlib: hamlib
	@echo "==> Staging Hamlib runtime files for packaging"
	mkdir -p out/hamlib/lib out/hamlib/bin
	cp -a $(HAMLIB_PREFIX)/lib/libhamlib* out/hamlib/lib/ 2>/dev/null || true
	cp -a $(HAMLIB_PREFIX)/bin/* out/hamlib/bin/ 2>/dev/null || true

.PHONY: mac-package
mac-package: release stage-hamlib
	# http://s.sudre.free.fr/Software/Packages/about.html
	packagesbuild --package-version $(VERSION) macos/kel-agent.pkgproj
	productsign --keychain `security list-keychains | grep k0swe | tr -d \"` \
      --sign "Developer ID Installer: Chris Keller (2UK8VD3UP4)" \
      kel-agent.pkg kel-agent-signed.pkg
	mv kel-agent-signed.pkg kel-agent_mac.pkg

.PHONY: win-package
win-package: release stage-hamlib
	# https://wixtoolset.org/docs/tools/wixexe/
	cd win && wix build kel-agent.wxs -arch x64 -o kel-agent.msi

# ---------------------------------------------------------------------------
# Verification helpers
# ---------------------------------------------------------------------------

.PHONY: verify-deps
verify-deps: kel-agent
ifeq ($(OS),linux)
	@echo "==> Linux shared library dependencies:"
	ldd kel-agent
endif
ifeq ($(OS),darwin)
	@echo "==> macOS shared library dependencies:"
	otool -L kel-agent
endif

# ---------------------------------------------------------------------------
# Housekeeping
# ---------------------------------------------------------------------------

.PHONY: clean
clean:
	rm -rf $(GENERATED) out/
