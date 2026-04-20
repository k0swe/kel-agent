# Building and Packaging

This document now describes the intended build and packaging direction for kel-agent as Hamlib
support becomes a first-class dependency.

The short version is:

- Debian packages should stay distro-native.
- Flatpak should stay self-contained.
- macOS and Windows releases should be self-contained by bundling Hamlib runtime libraries.
- Pure static linking is not the primary goal.

That approach gives us control over desktop release artifacts without fighting Debian packaging
policy or CGO cross-compilation constraints.

## Current state

The current branch still uses the older build flow:

- `make` builds Hamlib from source under `build/Hamlib-4.5.1/`.
- `PKG_CONFIG_PATH` points at that local Hamlib prefix.
- Debian, Flatpak, macOS, and Windows packaging metadata are maintained separately.
- Version updates are still manual in multiple files.

Those details are transitional. The target model is described below.

## Target artifact model

kel-agent should produce three kinds of release artifacts.

### Debian packages

Debian packages should be built against distro-provided Hamlib development packages and should
resolve runtime dependencies through normal Debian shared-library packaging.

This means:

- no repo-local Hamlib prefix in the Debian packaging path
- no requirement to bundle Hamlib into the `.deb`
- `debian/control` and `debian/rules` should describe the Hamlib dependency explicitly

This keeps Debian and Raspberry Pi packaging policy-friendly and easier to maintain.

### Flatpak

Flatpak should remain a hermetic build and runtime environment.

This means:

- Hamlib is built inside the Flatpak build sandbox
- the Flatpak manifest remains responsible for the exact Hamlib version used there
- host system Hamlib should not be required

### macOS and Windows installers

macOS and Windows releases should bundle the Hamlib runtime files needed by kel-agent.

This means:

- build Hamlib natively for each target platform
- package the resulting shared libraries beside the application binary or inside the installer
  payload
- treat full static linking as optional follow-up work, not the initial requirement

This is the most practical way to keep release artifacts self-contained without overcomplicating the
linker configuration.

## Build modes

The build system should move toward three explicit modes.

### 1. Developer build

For local development, kel-agent should build against an available Hamlib installation or a locally
built Hamlib prefix.

Expected properties:

- fast iteration
- easy local testing
- no installer or release artifact creation

### 2. Release build

Release builds should consume a known Hamlib input for the target platform and produce an
application binary plus any runtime files that must ship with it.

Expected properties:

- versioned Hamlib input
- reproducible target-specific output
- no hidden dependence on the developer machine

### 3. Packaging build

Packaging should take already-built release outputs and turn them into platform artifacts.

Expected properties:

- Debian packaging uses distro-native linking
- Flatpak packages a hermetic build
- macOS and Windows packaging bundle the runtime files created during release build

## Hamlib strategy

Hamlib is now a real platform dependency, so it should be managed as its own build input rather than
as an incidental side effect of `make`.

The intended direction is:

- keep a single tracked Hamlib version for the project
- build Hamlib separately per target platform
- publish or cache Hamlib artifacts by version and target triple
- feed those artifacts into kel-agent release builds

Expected artifact layout:

```text
out/
	hamlib/
		<hamlib-version>/
			<platform>-<arch>/
				include/
				lib/
				bin/
```

The exact directory structure may change, but the important part is that Hamlib becomes a clear,
versioned input to the application build.

## CI and release direction

Because Hamlib is accessed through CGO, native builds are preferred over aggressive
cross-compilation.

The intended CI matrix is:

- Linux runners for Debian packaging and Flatpak
- macOS runners for macOS packaging
- Windows runners for MSI packaging

CI should eventually:

- build or restore Hamlib artifacts for each target
- build kel-agent against those artifacts
- package the platform-specific release outputs
- verify runtime dependencies with platform-native inspection tools

Examples of the verification we want:

- Linux: `ldd`
- macOS: `otool -L`
- Windows: dependency inspection as part of the packaging job
- archive and installer content checks in every release workflow

## Optional Hamlib-free builds

Hamlib is runtime-optional in configuration today, but it is still compile-time required because the
current code imports `goHamlib` directly.

The build should eventually support a Hamlib-disabled path so that:

- environments without Hamlib can still build and test the rest of the program
- unsupported targets fail clearly instead of implicitly through CGO setup
- CI can separately validate Hamlib-enabled and Hamlib-disabled builds

This likely requires build tags or a small abstraction boundary around the Hamlib integration.

## Versioning direction

Version metadata is currently duplicated across several packaging files. That does not scale once
Hamlib artifacts are added.

The target state is:

- one canonical kel-agent version source
- one canonical Hamlib version source
- generated or synchronized packaging metadata derived from those sources

Today, version bumps still require checking all of the following manually:

- `debian/changelog`
- `assets/radio.k0swe.Kel_Agent.metainfo.xml`
- `flatpak/radio.k0swe.Kel_Agent.yml`
- `macos/kel-agent.pkgproj`
- `win/kel-agent.wxs`

That manual process should be reduced as the new build system is implemented.

## Transition plan

The implementation work should proceed in roughly this order.

1. Separate developer, release, and packaging workflows in the Makefile and helper scripts.
2. Remove the Debian packaging path's dependence on the repo-local Hamlib build prefix.
3. Define the Hamlib artifact layout and native build flow for Linux, macOS, and Windows.
4. Keep Flatpak bundled, but align its Hamlib version with the rest of the project.
5. Update macOS and Windows packaging to bundle Hamlib runtime files.
6. Add a Hamlib-disabled compile path.
7. Centralize version and dependency metadata.
8. Move release builds into a native CI matrix.

## Current commands

Until the new build flow is implemented, these are still the relevant commands for the existing
branch.

### Current local build

```shell
make
```

### Current Debian packaging flow

```shell
sudo apt install build-essential debhelper dh-golang sbuild autorevision
export ARCH=$(dpkg --print-architecture)
sudo sbuild-createchroot stable /srv/chroot/stable-"$ARCH" http://deb.debian.org/debian
make deb-package
```

### Current Flatpak packaging flow

```shell
sudo apt install flatpak flatpak-builder appstream-util desktop-file-validate
flatpak install flathub runtime/org.freedesktop.Sdk.Extension.golang/x86_64/20.08
make flatpak
```

## Release checklist during the transition

Until versioning and packaging metadata are centralized, a release still requires manual
cross-checking.

1. Update Go dependencies as needed.
2. Refresh `assets/modules.txt` if vendored module metadata changes.
3. Cross-check `go.mod`, the Flatpak manifest, and Debian build dependencies.
4. Update release notes in `debian/changelog` and `assets/radio.k0swe.Kel_Agent.metainfo.xml`.
5. Update versioned packaging metadata in macOS and Windows packaging files.
6. Build and validate Debian packages on supported Linux architectures.
7. Build and validate the Flatpak package.
8. Once implemented, build and validate macOS and Windows self-contained release artifacts.
