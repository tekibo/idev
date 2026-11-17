# Flutter

```bash
idev setup --flow flutter
idev flutter setup
idev flutter build
idev flutter ship
```

`idev flutter setup` installs the pinned Flutter toolchain into the private
prefix. `build` runs the Flutter build with the Darwin SDK and the private
CocoaPods/Swift environment. `ship` builds and ships the app.

Useful environment:

- `DARWIN_SDK_BUNDLE` is set by the manager.
- `IDEV_HOME` is set by the manager.
- The Dart work cache stays under `$IDEV_HOME/cache/flutter`.

For the no-Xcode path:

```bash
idev sdk-free setup
idev sdk-free flutter-build --debug /path/to/app
idev sdk-free ship /path/to/app
```

`sdk-free` is experimental. Its sysroot lives under `$IDEV_HOME/sdk-free`.
`idev sdk-free flutter-build` uses the Flutter parser at
`apple-tools/../flutter/tools/gen-info-plist.py` and passes `IOS_DIR` plus the
build configuration.

## Rules

- The Flutter SDK is optional; install it only when the project needs it.
- A real Flutter build needs the Flutter SDK and a project. `idev next` shows
  whether the flow is installed.
- `idev doctor` reports the Flutter payload as `private.flutter_arm`.
