#!/usr/bin/env bash
set -eu
echo "$ROOT_DIR"
cd "$ROOT_DIR" || exit 1

# VERSION can be set by the caller (e.g. Makefile); fall back to versions.env.
if [ -z "${VERSION:-}" ]; then
  # shellcheck source=../versions.env
  source "$ROOT_DIR/versions.env"
  VERSION="v${KEL_AGENT_VERSION}"
fi

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

TAGS="${TAGS:-hamlib}"

# shellcheck disable=SC2086
set -x
go build $mod -tags "$TAGS" --ldflags "$LDFLAGS"
{ set +x; } 2>/dev/null
