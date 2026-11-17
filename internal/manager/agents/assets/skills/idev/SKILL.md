---
name: idev
description: Build, run, sign, ship, and debug iOS apps on Linux with idev and xtool. Use for SwiftUI, Expo, React Native, Flutter, or Kotlin Multiplatform on a Linux host; for idev setup, doctor, next, update, config, devices, signing, TestFlight, App Store Connect, app icons, or the iPad screen viewer. Triggers on "idev", "xtool", "iOS on Linux", "run on my iPhone", "ship to TestFlight", "build the IPA", "Flutter iOS", "KMP framework", "iPad screen", or a project containing xtool.yml, app.json, pubspec.yaml, or an Xcode project.
---

# idev

idev is the native iOS build manager for Linux. It wraps xtool, the Darwin SDK,
Swift, CocoaPods, pymobiledevice3, Flutter, KMP, and its own React Native/Expo
engine. Projects stay ordinary SwiftUI, Expo, React Native, Flutter, or KMP apps.

If `idev` is not installed:

```bash
curl -fsSL https://github.com/tekibo/idev/releases/latest/download/install.sh | sh
```

Start every session with:

```bash
idev next          # setup state and the exact next command
idev doctor        # the full diagnostic report
```

`idev next --json` and `idev doctor --json` are the machine-readable forms.

## Ground rules

- **Run `idev next` before assuming anything about the host.** It is read-only.
- **Ask before changing the host**: packages, groups, sysctl, or anything sudo.
  Show the exact command from `idev system plan`.
- **Never ask for an Apple ID password, 2FA code, app-specific password, or API
  key in chat.** The user signs in through xtool or sets the ASC variables.
- **Xcode.xip is user-supplied.** idev never downloads it.
- **Never add ad-hoc signing.** xtool is the only signer and provisioner.
- **Do not create a new bundle ID without explicit approval.** The project's
  `ios.bundleIdentifier` is the App ID slot.
- **Long builds and tests run in the background.** Report progress; do not sleep
  in a loop.
- **No AUR.** Swift is the official swift.org tarball; system packages use the
  detected manager through `idev system plan`.
- Secrets stay out of commands, logs, docs, tests, and chat.

## Command map

| Task | Start with | Reference |
|---|---|---|
| Install or repair the toolchain | `idev next`, then `idev setup --base` | [setup](references/setup.md) |
| SwiftUI project | `idev swift run`, `idev device-run`, `idev ship` | [swiftui](references/swiftui.md) |
| Expo or bare React Native | `idev expo run:ios`, `idev rn run-ios` | [expo](references/expo.md) |
| Flutter | `idev flutter build`, `idev flutter ship` | [flutter](references/flutter.md) |
| Kotlin Multiplatform | `idev kmp` | [kmp](references/kmp.md) |
| Devices, install, launch, auth | `idev devices`, `idev auth status` | [devices](references/devices.md) |
| TestFlight, App Store Connect | `idev ship`, `idev asc` | [ship](references/ship.md) |
| iPad screen and control | `idev ipad-screen` | [ipad-screen](references/ipad-screen.md) |
| Updates and cleanup | `idev update apply --all`, `idev update clean` | [setup](references/setup.md) |
| Failures | `idev doctor`, `idev next` | [troubleshooting](references/troubleshooting.md) |

Do not guess a flag. Run `idev help <command>` or `idev <command> --help`;
the help tree documents every subcommand.

## Setup in one pass

```bash
idev next
idev system plan          # one exact command; ask before running it
idev system apply         # masked sudo prompt
idev setup --base         # payloads and private components
# put Xcode.xip somewhere, then:
idev sdk install /path/to/Xcode.xip
idev setup --apple-tools
idev setup --flow reactnative    # or flutter, kmp
idev doctor
```

`idev next` prints the one step that is actually missing. Do not run the whole
sequence blindly after a failure.

## Project workflows

- **SwiftUI:** `idev swift run` builds, installs, and launches. `idev ship
  --device` ships; `idev xcodeproj2xtool` converts an Xcode project.
- **Expo:** `idev expo run:ios` builds and launches. `idev expo start` runs Metro.
- **Bare React Native:** `idev rn run-ios`; `idev rn start` runs Metro.
- **Flutter:** `idev flutter build`; `idev flutter ship`.
- **KMP:** `idev kmp` links the framework from Gradle's klibs.

A project is detected by `xtool.yml`, `app.json`, `pubspec.yaml`, or an Xcode
project. `idev project info --json` prints the typed descriptor.

## Device and signing rules

- `idev devices --json` lists connected devices.
- `idev install-app --bundle <app> --bundle-id <id>` signs, installs, and launches.
- `idev launch-app --bundle-id <id>` launches an installed app.
- The installed identifier may be `XTL-<team>.<id>`; that is the signed installed
  value, not a new App ID.
- `idev auth status --json` reads the xtool account. `idev auth login` signs in.

## When something fails

Run `idev next` first, then `idev doctor`. If the report names a missing system
package, show the exact `idev system plan` command and ask before installing.
For source-pin conflicts, run `idev config migrate`. For a stale or unsupported
component, run `idev update check`.

Read [troubleshooting](references/troubleshooting.md) before retrying a command
that already failed twice.
