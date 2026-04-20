# Building and Packaging

This document describes the build and packaging strategy for kel-agent now that Hamlib support is
a first-class dependency.

The short version is:

- Debian packages stay distro-native.
- Flatpak stays self-contained.
- macOS and Windows releases are self-contained by bundling Hamlib runtime libraries.
- Pure static linking is not the primary goal.

That approach gives us control over desktop release artifacts without fighting Debian packaging
policy or CGO cross-compilation constraints.

## Version sources

All version metadata lives in `versions.env` at the repository root:

```shell
KEL_AGENT_VERSION=0.4.6
HAMLIB_VERSION=4.5.1
```

The Makefile, build scripts, and CI workflows read from this file. Packaging metadata in
`debian/changelog`, `assets/radio.k0swe.Kel_Agent.metainfo.xml`, `flatpak/radio.k0swe.Kel_Agent.yml`,
`macos/kel-agent.pkgproj`, and `win/kel-agent.wxs` should be kept in sync manually during releases
until further automation is added.

## Build tags

Hamlib is compile-time optional via the Go build tag `hamlib`.

- **With Hamlib** (`-tags hamlib`): full rig-control integration via goHamlib and libhamlib.
  Requires `libhamlib-dev` (or equivalent) to be installed.
- **Without Hamlib** (default): a stub logs a warning when Hamlib is configured; the rest of the
  program builds and runs normally.

This allows environments without Hamlib to build and test, and lets CI validate both paths.

## Build modes

### 1. Developer build

```shell
make                  # Hamlib-enabled (needs libhamlib-dev or a local Hamlib prefix)
make test-nohamlib    # Hamlib-disabled quick iteration
```

The developer build compiles against whichever Hamlib is available to `pkg-config`. The Makefile
derives `PKG_CONFIG_PATH` from the Hamlib artifact layout under `out/` when a local build exists.

### 2. Release build

```shell
make hamlib    # build Hamlib from source into out/hamlib/<version>/<os>-<arch>/
make release   # build kel-agent against the Hamlib artifacts
```

`scripts/build-hamlib.sh` downloads and builds Hamlib into a versioned, platform-specific prefix:

```text
out/
  hamlib/
    <hamlib-version>/
      <os>-<arch>/
        include/
        lib/
        bin/
```

Release builds consume these artifacts through `PKG_CONFIG_PATH` so that the binary and any bundled
runtime files are reproducible for the target platform.

### 3. Packaging build

Packaging targets take already-built outputs and produce platform artifacts.

```shell
make deb-package   # Debian .deb via sbuild (uses distro libhamlib-dev, not repo-local prefix)
make flatpak       # hermetic Flatpak bundle (builds Hamlib inside the sandbox)
make mac-package   # macOS .pkg (bundles Hamlib dylib from out/)
make win-package   # Windows .msi (bundles Hamlib DLL from out/)
```

## Artifact model

### Debian packages

Debian packages are built against distro-provided Hamlib development packages and resolve runtime
dependencies through normal Debian shared-library packaging.

- `debian/control` declares `libhamlib-dev` as a Build-Depends and `libhamlib4` as a runtime Depends
- `debian/rules` passes `-tags hamlib` to `dh_auto_build`
- No repo-local Hamlib prefix is used in the Debian packaging path

### Flatpak

Flatpak remains a hermetic build and runtime environment. The Flatpak manifest builds Hamlib from
source inside the sandbox. The Hamlib version in the manifest should be kept aligned with
`HAMLIB_VERSION` in `versions.env`.

### macOS and Windows installers

macOS and Windows releases bundle the Hamlib runtime files needed by kel-agent. The `mac-package`
and `win-package` Makefile targets expect Hamlib artifacts to have been built via `make hamlib`
before packaging. The installers include the Hamlib shared libraries alongside the application
binary.

## CI and release

Because Hamlib is accessed through CGO, native builds are used instead of cross-compilation.

### Test workflow

The test workflow runs two jobs:

- **test-nohamlib**: builds and tests without the `hamlib` tag (no native dependencies)
- **test-hamlib**: installs `libhamlib-dev` and runs with the `hamlib` tag

### Release workflow

The release workflow uses a native CI matrix:

- Linux runners for Debian packaging and Flatpak
- macOS runners build Hamlib, then build and package kel-agent
- Windows runners build Hamlib, then build and package kel-agent

Each platform job includes verification steps to inspect runtime dependencies (`ldd`, `otool -L`,
package content checks).

## Verification helpers

```shell
make verify-deps   # show shared library dependencies of the built binary (ldd or otool)
```

## Release checklist

1. Update `KEL_AGENT_VERSION` and/or `HAMLIB_VERSION` in `versions.env`.
2. Update Go dependencies as needed.
3. Refresh `assets/modules.txt` if vendored module metadata changes.
4. Cross-check `go.mod`, the Flatpak manifest, and Debian build dependencies.
5. Update release notes in `debian/changelog` and `assets/radio.k0swe.Kel_Agent.metainfo.xml`.
6. Update versioned packaging metadata in `win/kel-agent.wxs` and `macos/kel-agent.pkgproj`.
7. Build and validate Debian packages on supported Linux architectures.
8. Build and validate the Flatpak package.
9. Build and validate macOS and Windows self-contained release artifacts.
