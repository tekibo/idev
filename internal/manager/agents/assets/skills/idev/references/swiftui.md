# SwiftUI and xtool

## Create

```bash
idev new MyApp --bundle-id com.example.myapp
```

`--skip-setup` passes xtool's non-interactive setup flag. Without a bundle ID,
the project is created unset; set it before the first signed build.

## Inspect and change

```bash
idev project type
idev project info --json
idev project set-bundle-id --bundle-id com.example.myapp
idev project clean
```

The bundle ID lives in `xtool.yml` as the App ID slot. Do not create a new App ID
without the user's approval. The installed identifier may be
`XTL-<team>.<id>`; that is the signed installed value, not a new slot.

## Build and run

```bash
idev swift build
idev swift run
idev swift new MyApp
idev swift generate-xcode-project
idev swift build-server
```

`idev swift run` builds, installs, and launches through xtool and the private
device tooling. Run long builds in the background.

## Device run

```bash
idev device-run
idev device-run --console 30
idev device-run --lldb
idev device-run --attach
```

The helper uses the apple-tools payload and the Darwin SDK. It is the same
device-run path used by the SwiftUI project surface.

## Ship

```bash
idev ship --device
idev ship --device TEAMID
idev ship --upload
```

`--device` builds for a device. `--upload` uploads after signing. The App Store
record must already exist for the bundle ID. Use `idev asc` for the post-upload
release steps and `idev ship --upload` only after the user confirms the build and
signing are correct.

## App icon

```bash
idev compile-app-icon /path/to/App.app
```

The compiler resolves the SDK from `DARWIN_SDK_BUNDLE` and uses `actool`. The
Icon Composer `automatic-gradient` background fill is supported by the pinned
AssetKit revision.

## Xcode project conversion

```bash
idev xcodeproj2xtool /path/to/App.xcodeproj --out omarchy-xtool
idev xcodeproj2xtool --self-test
```

The converter preserves Swift package pins, exact versions, revisions, and
branch requirements. A branch pin warns when `Package.resolved` holds a commit
the build does not follow.

## Rules

- xtool is the only signer and provisioner. Never add ad-hoc signing.
- Keep the project an ordinary xtool/SwiftUI app: no per-project dev dependency,
  no `postinstall`, no patched `node_modules`, no project shim.
- Use `idev next` after a failure before retrying.
