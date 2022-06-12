# Build requirements

Building the basic program should only require a recent version of golang.

```shell
make
```

## Bumping the version

1. `go get -u && go mod tidy`
2. `rm assets/modules.txt && make assets/modules.txt`
3. Cross-reference `assets/modules.txt` with `flatpak/radio.k0swe.Kel_Agent.yml`
4. Cross-reference `go.mod` with `Makefile`
5. Cross-reference `go.mod` with `debian/control` `Build-Depends`
6. Run `make deb-package` on Linux amd64 and Linux arm to make sure `chroot`s are set up
7. Run `make flatpak` on Linux amd64 to make sure that's building
8. Add changelog entries in `debian/changelog` and `assets/radio.k0swe.Kel_Agent.metainfo.xml`
9. Bump versions in `macos/kel-agent.pkgproj` and `win/kel-agent.wxs`

## Packaging for Debian Linux (incl. Raspberry Pi)

```shell
sudo apt install build-essential debhelper dh-golang sbuild autorevision
export ARCH=$(dpkg --print-architecture)
sudo sbuild-createchroot stable /srv/chroot/stable-"$ARCH" http://deb.debian.org/debian
make deb-package
```

## Packaging for Flatpak

```shell
sudo apt install flatpak flatpak-builder appstream-util desktop-file-validate
flatpak install flathub runtime/org.freedesktop.Sdk.Extension.golang/x86_64/20.08
make flatpak
```
