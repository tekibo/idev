# Design — how Expo builds on Linux

Deep engine reference. The supported entry point is `idev expo run:ios`. The arm is installed at `$IDEV_HOME/expo`. The engine binary is `$IDEV_HOME/bin/expo-build-bin`. Start with `docs/COMPONENTS.md` and `docs/EXPO-CLI.md`.

## Layers

| Layer | Owns | Does not own |
|---|---|---|
| `xtool` | signing, install, SwiftPM projects | `.xcodeproj` |
| Apple tooling | Darwin tools, OAM, shipping helpers | `.xcodeproj` |
| Swift Build | the build engine Xcode uses | `.xcodeproj` parsing |

Swift Build consumes PIF. Xcode's closed front end turns a project into PIF. Swift Build turns PIF into a build.

In Swift Build, the service contract maps `xcodeproj` and `xcworkspace` to `dumpPIF`, and `json` and `pif` to a direct path. On Linux the missing piece is a real `xcodebuild -dumpPIF`. Compilation, linking, module maps, codegen ordering, asset catalogs, storyboards and script phases are already inside Swift Build.

## PIF wire format

A PIF file is a JSON array:

```json
{ "signature": "PROJECT@v11_P1", "type": "project", "contents": {} }
```

`contents` is a plain dictionary. The MessagePack envelope used elsewhere is an optimisation, not a requirement.

The schema vocabulary lives in `SWBProtocol/PIFKeyConstants.swift` (`guid`, `type`, `buildPhases`, `productTypeIdentifier`, `scriptContents`).

Consequence: `payloads/react-native` emits PIF directly from Python. No Swift component, no MessagePack, no duplicated schema in another language.

## Architecture

```text
project.pbxproj + xcconfig + Podfile.lock
        │  payloads/react-native/lib      project model: targets, phases, settings, file refs
        ▼
     app.pif           plain JSON, the document type Xcode hands to Swift Build
        │  expo-build-bin
        ▼
   .app  ──►  xtool install (sole signer)
```

Python owns only reading Xcode and CocoaPods data. Swift Build owns compilation, linking, module maps, resource compilation, script phases and dependency ordering. Nothing reimplements a build system.

## Why not SwiftPM

SwiftPM cannot express a CocoaPods target: one target that compiles Swift and Objective-C with a custom module map and emits `<Pod>-Swift.h` for its Objective-C consumers. That is why the Flutter path requires every plugin to ship a `Package.swift`. Expo cannot take that route.

## Platform registration

The `swbuild` console creates every session with `developerPath: nil`, and never passes `resourceSearchPaths`. With no developer directory, `PlatformRegistry` registers only generic platforms. `iphoneos` is absent, so planning fails with `unable to resolve run destination platform: 'iphoneos'` and then `unable to find sdk 'iphoneos'`. Setting `DEVELOPER_DIR` does not help: the value is passed over the session protocol, not read from the environment.

The library API supports it:

```swift
createSession(name:developerPath:resourceSearchPaths:cachePath:
              inferiorProductsPath:environment:)
```

The engine entry point is a small Swift driver that:

1. creates a session with `developerPath` set to the Darwin SDK bundle `Developer` directory, and `resourceSearchPaths` set to its platform and toolchain paths;
2. sends the generated PIF;
3. runs the build request.

`swbuild` stays useful for diagnostics (`showPlatforms`, `showSDKs`).

## Legacy proof and native proof

The build engine path was proven end to end in the legacy container on Expo SDK 58:

- `payloads/react-native/tools/dump-pif` produces a PIF Swift Build accepts.
- `payloads/react-native/bin/build` prepares a Developer overlay, stages CocoaPods xcframeworks and prebuilt pod products, and runs the engine through `expo-build-bin`.
- Release and Debug builds completed in that environment.
- xtool signed and installed. pymobiledevice3 launched.

That proof is historical. The native host proof is deferred until the user tests after full GUI completion. Do not treat the container result as a current native proof.

Remaining work is not in the core engine:

- native component installer and manifest;
- Expo CLI adapter per SDK and the manager-owned install hook;
- GUI and setup flow.

## DEVELOPER_DIR and SwiftPM

Pointing `DEVELOPER_DIR` at the prepared overlay makes SwiftPM's Swift Build backend work for iOS too:

```bash
export DEVELOPER_DIR=/tmp/expo-developer
swift build -c release --swift-sdk darwin --triple arm64-apple-ios --product <Product>
```

Without it SwiftPM fails with `unable to find platform for 'iphoneos'`, the same failure `swbuild` has. This is the route to build SPM-based pods on Linux without `xcodebuild`.

## Reading Swift Build diagnostics

The engine `output` and `taskOutput` events carry plain base64 text, not a compressed stream. Decode it to expose compiler errors:

```python
base64.b64decode(message["data"]).decode("utf-8", "replace")
```

## Environment facts

- Build output stays unsigned. `CODE_SIGNING_ALLOWED=NO` and related settings are set in the PIF so CocoaPods embed phases never sign. xtool signs at install.
- Script phases run from `SRCROOT` with the caller `PATH` and sandboxing disabled.
- The `[CP-User] Build ExpoModulesJSI xcframework` and `[Expo] Configure project` phases are neutralised: the engine builds the xcframework from its SPM package, and `pod install` writes the generated file.
- `ditto` and `lipo` shims live under `$IDEV_HOME/expo/tools`. `llvm-lipo` backs the latter.
- `HERMES_CLI_PATH` points at the Linux `hermesc` from the `hermes-compiler` package.
- `SKIP_BUNDLING`, `RCT_METRO_PORT`, `RCT_TERMINAL`, `__EXPO_EAGER_BUNDLE_OPTIONS` and `NODE_BINARY` are carried from the caller into the bundle phase.
- Swift compatibility headers are emitted only for pods whose sources import `<Pod/Pod-Swift.h>`.
- The `xcodebuild` shim installs CocoaPods deps itself and prints the two build-output exports the CLI parses.
- `payloads/react-native/bin/build` puts `$HERE/../tools` on `PATH`, so build phases see `ditto`, `lipo`, `ipconfig` and the rest.

## Static frameworks

Static frameworks must not be embedded. CocoaPods can call `install_framework` for a vendored xcframework whose binary is an `ar` archive. xtool's signer parses every framework binary as Mach-O and fails with `Invalid mach-O file`. The archive is already linked into the executable, so `payloads/react-native/tools/prune-static-frameworks` removes frameworks whose main binary is an archive after a successful build.

## Device install

- `usbmuxd` must be the host daemon. A second daemon cannot own the USB interface.
- The manager path calls `idev install-app` to `xtool install` and `pymobiledevice3` launch.
- Developer Mode must be on.
- The development certificate must be trusted on the device.
- Free-provisioning registration is automatic but eventual. `xtool install` registers the device before creating the profile. Apple's profile API can lag. Running `xtool install` again resolves it.
- The build stays unsigned. xtool is the only signer.

## Verification notes

- `taskUpToDate` hides real work when the arena moves. Clear the arena between verification runs.
- Swift Build cannot read serialized diagnostics because it wants a Mach-O `libclang`. Compile errors still surface through task results.
