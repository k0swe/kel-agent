#!/usr/bin/env bash
# Build Hamlib from source into the artifact layout described in BUILDING.md.
#
# Usage:
#   scripts/build-hamlib.sh
#
# The script reads HAMLIB_VERSION from versions.env and installs into:
#   out/hamlib/<version>/<os>-<arch>/
#
# Requirements (non-Windows):
#   - autoconf, automake, libtool, pkg-config, make, gcc (or platform equivalent)
#   - wget or curl
# Requirements (Windows/MinGW): curl or wget, unzip, pkg-config (from MSYS2 MinGW64)
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

# On Windows/MinGW, autotools (autoreconf, automake, libtool) are not available.
# Download the official pre-built Windows binaries instead of building from source.
if [[ "$OS" == mingw* || "$OS" == msys* ]]; then
    if [ ! -d "$PREFIX" ]; then
        HAMLIB_WIN_ZIP="hamlib-w64-${HAMLIB_VERSION}.zip"
        HAMLIB_WIN_URL="https://github.com/Hamlib/Hamlib/releases/download/${HAMLIB_VERSION}/${HAMLIB_WIN_ZIP}"
        HAMLIB_WIN_EXTRACT="$ROOT_DIR/build/hamlib-win-extract"

        echo "==> Downloading pre-built Hamlib ${HAMLIB_VERSION} for Windows"
        if command -v curl >/dev/null 2>&1; then
            curl -fsSL -o "$ROOT_DIR/build/${HAMLIB_WIN_ZIP}" "$HAMLIB_WIN_URL"
        else
            wget -q -O "$ROOT_DIR/build/${HAMLIB_WIN_ZIP}" "$HAMLIB_WIN_URL"
        fi

        rm -rf "$HAMLIB_WIN_EXTRACT"
        unzip -q "$ROOT_DIR/build/${HAMLIB_WIN_ZIP}" -d "$HAMLIB_WIN_EXTRACT"
        rm -f "$ROOT_DIR/build/${HAMLIB_WIN_ZIP}"

        HAMLIB_WIN_DIR="$HAMLIB_WIN_EXTRACT/hamlib-w64-${HAMLIB_VERSION}"
        mkdir -p "$PREFIX/bin" "$PREFIX/include" "$PREFIX/lib/pkgconfig"
        cp -r "$HAMLIB_WIN_DIR/bin/." "$PREFIX/bin/"
        cp -r "$HAMLIB_WIN_DIR/include/." "$PREFIX/include/"
        # The GCC import library lives under lib/gcc/ in the Windows release zip
        cp "$HAMLIB_WIN_DIR/lib/gcc/libhamlib.dll.a" "$PREFIX/lib/"

        # Write a pkg-config file so CGo (#cgo pkg-config: hamlib) can locate hamlib
        cat > "$PREFIX/lib/pkgconfig/hamlib.pc" << EOF
prefix=${PREFIX}
exec_prefix=\${prefix}
libdir=\${exec_prefix}/lib
includedir=\${prefix}/include

Name: hamlib
Description: Ham Radio Control Library
Version: ${HAMLIB_VERSION}
Libs: -L\${libdir} -lhamlib
Cflags: -I\${includedir}
EOF
    fi

    echo "==> Hamlib ${HAMLIB_VERSION} installed to ${PREFIX}"
    exit 0
fi

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
