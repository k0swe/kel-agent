all: kel-agent
.PHONY: all

kel-agent:
	scripts/build.sh

architecture.svg:
	dot -T svg -o architecture.svg < architecture.dot

.PHONY: deb-package
deb-package:
	dpkg-buildpackage -uc -us -b

.PHONY: clean
clean:
	rm -f kel-agent
