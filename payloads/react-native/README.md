# expo — Expo / React Native iOS builds on Linux

This directory is the Expo/React Native arm of `idev`, in the same sense that
The Flutter arm is separate. `xtool` is the base.

The `idev` manager installs this arm at `$IDEV_HOME/react-native` and drives it through
`idev expo run:ios`. The arm is portable native: it does not assume Docker and does not assume a
particular container path. The manager supplies the CLI hook, project paths and
toolchain environment. See `docs/COMPONENTS.md` and `docs/EXPO-CLI.md`.

It does **not** reimplement compilation. It supplies the missing piece between
an Expo app and the engines already on this machine:

- `xtool` builds SwiftPM packages and signs/installs. It does not read
  `.xcodeproj`.
- The Apple tooling payload supplies the Darwin tools, OAM, and shipping helpers
  and tool shims. It does not read `.xcodeproj` either.
- Xcode's own build engine (Swift Build / `SWBCore`) ships inside the Swift
  toolchain but is driven by **PIF**, not by `.xcodeproj`.

So this directory owns exactly one thing: Xcode project semantics.

```
.xcworkspace + app .xcodeproj + Pods.xcodeproj
        │
        ▼
  project model  (targets, build phases, xcconfigs, dependencies)
        │
        ├── build phases     script phases, Hermes bundle, assets, storyboards
        └── PIF              Swift Build target graph
                │
                ▼
          Swift Build  →  .app  →  xtool install (sole signer)
```

## Rules

- Never patch `xtool`, the Apple tooling payload, Expo, React Native sources.
  Everything here is additive.
- Never sign. Build output stays unsigned; `xtool` is the only signer.
- No per-app, no per-SDK special cases. If the project model says it, apply it.
- Upstream fixes belong upstream. This tree only supplies what upstream
  deliberately does not ship (CocoaPods workspace support for Linux).

## Layering around upstream

- `expo/tools/ibtool` loads the Apple tooling ibtool and adds the device ids
  Expo/React Native templates use (`retina6_12`, `retina6_7`, `retina6_9`)
  before calling its `main`. Upstream still reports them as unsupported; the
  fix to send them is to add those entries to `DEVICE_SIZES`.
- Never patch a file under the private Apple tooling checkout, the xtool checkout. If a
  tool needs different behaviour, wrap it here.
