# Build requirements

Building the basic program should only require a recent version of golang.

```shell
make
```

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
