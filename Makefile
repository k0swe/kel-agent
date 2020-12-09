all: kel-agent
.PHONY: all

kel-agent:
	export VERSION=$$(< debian/changelog head -1 | sed -r 's/.*\(([0-9]+\.[0-9]+\.[0-9]+)-.*\).*/v\1/g')
	export GIT_REV=$$(git rev-parse --short HEAD)
	go build -ldflags "-X main.GitRev=$$GIT_REV -X main.Version=$$VERSION"

architecture.svg:
	dot -T svg -o architecture.svg < architecture.dot

.PHONY: deb-package
deb-package:
	dpkg-buildpackage -uc -us -b

.PHONY: clean
clean:
	rm kel-agent
