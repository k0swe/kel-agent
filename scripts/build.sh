#!/usr/bin/env bash
set -eu
echo "$ROOT_DIR"
cd "$ROOT_DIR" || exit 1

VERSION=${VERSION:-v$(< debian/changelog head -1 | egrep -o "[0-9]+\.[0-9]+\.[0-9]+")}
BUILDTIME=${BUILDTIME:-$(date -u +"%Y-%m-%dT%H:%M:%SZ")}

LDFLAGS=$(echo "\
    -w \
    -X \"main.Version=${VERSION}\" \
    -X \"main.GitCommit=${GITCOMMIT}\" \
    -X \"main.BuildTime=${BUILDTIME}\" \
    ${LDFLAGS:-} \
" | xargs echo)

mod=""
if [ -d vendor ]; then
  mod="-mod vendor"
fi

# shellcheck disable=SC2086
set -x
go build $mod --ldflags "$LDFLAGS"
{ set +x; } 2>/dev/null
