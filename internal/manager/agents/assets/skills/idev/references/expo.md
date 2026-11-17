# Expo and React Native

## Expo

```bash
idev expo run:ios
idev expo run:ios --configuration Release
idev expo run:ios --device generic --output ./build
idev expo start
idev expo prebuild
idev expo export
```

`idev expo run:ios` builds, installs, and launches. `--device generic` builds
without installing. `--no-bundler` skips Metro; `--binary` installs an existing
`.app` or `.ipa`.

`idev expo start` runs Metro with the idev adapter. The Expo CLI is the project's
own local CLI, not a global install. The adapter sets the private Swift, Darwin
SDK, CocoaPods, and React Native environment.

## Bare React Native

```bash
idev rn run-ios
idev rn build-ios
idev rn start
idev rn log-ios
```

`idev rn run-ios` accepts `--mode Debug|Release`, `--scheme`, `--device`,
`--udid`, `--binary-path`, `--no-packager`, `--port`, `--buildFolder`,
`--only-pods`, and `--force-pods`.

## Direct Expo/RN commands

`idev expo ...` and `idev rn ...` are passthroughs. Use `idev help expo` and
`idev help rn` for the current flags. The manager sets:

- `IDEV_HOME`
- `IDEV_APPLE_TOOLS`
- `IDEV_SWIFT_BIN`
- `DARWIN_SDK_BUNDLE`
- `PATH` with the private tools and project CLIs

## CocoaPods

```bash
idev pod install
idev pod update
```

The project's CocoaPods runs with the private gem tree and the Expo autolinking
patch. No system gem install.

## Rules

- Keep the project an ordinary Expo/React Native app. No per-project dev
  dependency, no `postinstall`, no patched `node_modules`.
- `idev setup --flow reactnative` installs the arm. `--no-expo` keeps bare RN.
- After a native dependency or config-plugin change, rebuild; Fast Refresh only
  covers JavaScript.
- `idev devices` before a device run; `idev launch-app --bundle-id <id>` for an
  app that is already installed.
