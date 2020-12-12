all: kel-agent
.PHONY: all

.PHONY: test
test:
	go test ./...
	go vet ./...

kel-agent: test
	scripts/build.sh

architecture.svg:
	# apt install graphviz
	dot -T svg -o architecture.svg < architecture.dot

.PHONY: deb-package
deb-package:
	# apt install build-essential devscripts
	dpkg-buildpackage -uc -us -b

.PHONY: mac-package
mac-package: kel-agent
	# http://s.sudre.free.fr/Software/Packages/about.html
	packagesbuild macos/kel-agent.pkgproj
	mv kel-agent.pkg kel-agent_mac.pkg

.PHONY: clean
clean:
	rm -f kel-agent kel-agent_mac.pkg
