#!/bin/bash
# Legacy Hamlib build script used by the Flatpak manifest.
# For non-Flatpak builds, prefer scripts/build-hamlib.sh which installs into
# the standard artifact layout under out/.

cd "$ROOT_DIR" || exit 1

# shellcheck source=../versions.env
source "$ROOT_DIR/versions.env"

mkdir -p build && cd build || exit 1
if [ ! -d "Hamlib-$HAMLIB_VERSION" ]; then
    wget "https://github.com/Hamlib/Hamlib/archive/refs/tags/$HAMLIB_VERSION.tar.gz"
    tar xvzf "$HAMLIB_VERSION.tar.gz"
    rm "$HAMLIB_VERSION.tar.gz"
fi
cd "Hamlib-$HAMLIB_VERSION" || exit 1
hamlib_repo=$(pwd)
hamlib_prefix="$hamlib_repo"/prefix/usr/local
mkdir -p "$hamlib_prefix"
./bootstrap
./configure --prefix="$hamlib_prefix"
make
make install
