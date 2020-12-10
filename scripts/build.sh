#!/usr/bin/env bash
set -eu

VERSION=${VERSION:-$(< debian/changelog head -1 | sed -r 's/.*\(([0-9]+\.[0-9]+\.[0-9]+)-.*\).*/v\1/g')}
GITCOMMIT=${GITCOMMIT:-$(git rev-parse --short HEAD 2> /dev/null || true)}
BUILDTIME=${BUILDTIME:-$(date -u +"%Y-%m-%dT%H:%M:%SZ")}

export LDFLAGS="\
    -w \
    -X \"main.Version=${VERSION}\" \
    -X \"main.GitCommit=${GITCOMMIT}\" \
    -X \"main.BuildTime=${BUILDTIME}\" \
    ${LDFLAGS:-} \
"

go build --ldflags "${LDFLAGS}"
