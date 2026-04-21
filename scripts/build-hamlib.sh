#!/usr/bin/env bash
# Build Hamlib from source into the artifact layout described in BUILDING.md.
#
# Usage:
#   scripts/build-hamlib.sh
#
# The script reads HAMLIB_VERSION from versions.env and installs into:
#   out/hamlib/<version>/<os>-<arch>/
#
# Requirements:
#   - autoconf, automake, libtool, pkg-config, make, gcc (or platform equivalent)
#   - wget or curl
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

# shellcheck source=../versions.env
source "$ROOT_DIR/versions.env"

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"
# Normalize architecture names
case "$ARCH" in
    x86_64)  ARCH="amd64" ;;
    aarch64) ARCH="arm64" ;;
    armv7l)  ARCH="armhf" ;;
esac

HAMLIB_TARBALL="$HAMLIB_VERSION.tar.gz"
HAMLIB_URL="https://github.com/Hamlib/Hamlib/archive/refs/tags/$HAMLIB_TARBALL"
BUILD_DIR="$ROOT_DIR/build/Hamlib-$HAMLIB_VERSION"
PREFIX="$ROOT_DIR/out/hamlib/$HAMLIB_VERSION/$OS-$ARCH"

echo "==> Building Hamlib $HAMLIB_VERSION for $OS-$ARCH"
echo "    build dir: $BUILD_DIR"
echo "    prefix:    $PREFIX"

mkdir -p "$ROOT_DIR/build"
if [ ! -d "$BUILD_DIR" ]; then
    echo "==> Downloading Hamlib $HAMLIB_VERSION"
    if command -v wget >/dev/null 2>&1; then
        wget -q -O "$ROOT_DIR/build/$HAMLIB_TARBALL" "$HAMLIB_URL"
    elif command -v curl >/dev/null 2>&1; then
        curl -fsSL -o "$ROOT_DIR/build/$HAMLIB_TARBALL" "$HAMLIB_URL"
    else
        echo "ERROR: neither wget nor curl is available" >&2
        exit 1
    fi
    tar xzf "$ROOT_DIR/build/$HAMLIB_TARBALL" -C "$ROOT_DIR/build"
    rm -f "$ROOT_DIR/build/$HAMLIB_TARBALL"
fi

cd "$BUILD_DIR"
if [ ! -f configure ]; then
    echo "==> Running bootstrap"
    ./bootstrap
fi

mkdir -p "$PREFIX"
echo "==> Configuring"
./configure --prefix="$PREFIX" --quiet
echo "==> Building"
make -j"$(nproc 2>/dev/null || sysctl -n hw.ncpu 2>/dev/null || echo 1)" --quiet
echo "==> Installing to $PREFIX"
make install --quiet

echo "==> Hamlib $HAMLIB_VERSION installed to $PREFIX"
