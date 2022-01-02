VERSION = $(shell < debian/changelog head -1 | egrep -o "[0-9]+\.[0-9]+\.[0-9]+")
GITCOMMIT = $(shell git rev-parse --short HEAD 2> /dev/null || true)

GENERATED = kel-agent kel-agent_*.pkg win/kel-agent_*.msi win/kel-agent.wixobj autorevision.cache \
  ../kel-agent_* ../*.deb flatpak/repo/ flatpak/.flatpak-builder/ flatpak/kel_agent.flatpak \
  flatpak/flatpak_app/ flatpak/build-out/

.PHONY: all
all: kel-agent

internal/webview/kel-agent-gui/dist:
	cd internal/webview/kel-agent-gui && npm install && npm run build

.PHONY: test
test: internal/webview/kel-agent-gui/dist
	go test ./...
	go vet ./...
	if command -v appstream-util; then appstream-util validate-relax --nonet assets/radio.k0swe.Kel_Agent.metainfo.xml; fi
	if command -v desktop-file-validate; then desktop-file-validate assets/radio.k0swe.Kel_Agent.desktop; fi

assets/modules.txt:
	go mod vendor
	mv vendor/modules.txt assets/
	rm -rf vendor

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

# TODO: This target can be removed once the package is in Debian stable and Ubuntu stable
wsjtx-go.deb:
	cd .. && \
	wget https://github.com/k0swe/wsjtx-go/releases/download/v3.1.0/golang-github-k0swe-wsjtx-go-dev_3.1.0-1_all.deb

# TODO: This target can be removed once the package is in Debian stable and Ubuntu stable
adrg-xdg.deb:
	cd .. && \
	wget http://ftp.us.debian.org/debian/pool/main/g/golang-github-adrg-xdg/golang-github-adrg-xdg-dev_0.3.3-2_all.deb

.PHONY: deb-package
deb-package: deb-tarball wsjtx-go.deb adrg-xdg.deb
	# https://wiki.debian.org/sbuild
	sbuild -d stable \
      --extra-package=../golang-github-k0swe-wsjtx-go-dev_3.1.0-1_all.deb \
      --extra-package=../golang-github-adrg-xdg-dev_0.3.3-2_all.deb

.PHONY: flatpak
flatpak: kel-agent
	cd flatpak && \
      flatpak-builder --force-clean build-out radio.k0swe.Kel_Agent.yml --repo=repo && \
      flatpak build-bundle repo kel_agent.flatpak radio.k0swe.Kel_Agent main

.PHONY: mac-package
mac-package: kel-agent
	# http://s.sudre.free.fr/Software/Packages/about.html
	packagesbuild --package-version $(VERSION) macos/kel-agent.pkgproj
	productsign --keychain `security list-keychains | grep k0swe | tr -d \"` \
      --sign "Developer ID Installer: Chris Keller (2UK8VD3UP4)" \
      kel-agent.pkg kel-agent-signed.pkg
	mv kel-agent-signed.pkg kel-agent_mac.pkg

.PHONY: win-package
win-package: kel-agent
	# https://wixtoolset.org/
	cd win && candle kel-agent.wxs && light kel-agent.wixobj

.PHONY: clean
clean:
	rm -rf $(GENERATED)
