#!/usr/bin/env bash
set -eu

VERSION=${VERSION:-v$(< debian/changelog head -1 | egrep -o "[0-9]+\.[0-9]+\.[0-9]+")}
BUILDTIME=${BUILDTIME:-$(date -u +"%Y-%m-%dT%H:%M:%SZ")}

export LDFLAGS="\
    -w \
    -X \"main.Version=${VERSION}\" \
    -X \"main.GitCommit=${GITCOMMIT}\" \
    -X \"main.BuildTime=${BUILDTIME}\" \
    ${LDFLAGS:-} \
"

mod=""
if [ -d vendor ]; then
  mod="-mod vendor"
fi

# shellcheck disable=SC2086
go build $mod --ldflags "$LDFLAGS"
