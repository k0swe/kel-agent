#!/bin/bash

cd "$ROOT_DIR" || exit 1
mkdir -p build && cd build || exit 1
if [ ! -d "Hamlib-4.5.1" ]; then
    wget https://github.com/Hamlib/Hamlib/archive/refs/tags/4.5.1.tar.gz
    tar xvzf 4.5.1.tar.gz
    rm 4.5.1.tar.gz
fi
cd Hamlib-4.5.1 || exit 1
hamlib_repo=$(pwd)
hamlib_prefix="$hamlib_repo"/prefix/usr/local
mkdir -p "$hamlib_prefix"
./bootstrap
./configure --prefix="$hamlib_prefix"
make
make install
