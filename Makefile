all: kel-agent
.PHONY: all

kel-agent:
	go build

architecture.svg:
	dot -T svg -o architecture.svg < architecture.dot

.PHONY: deb-package
deb-package:
	dpkg-buildpackage -uc -us -b

.PHONY: clean
clean:
	rm kel-agent
