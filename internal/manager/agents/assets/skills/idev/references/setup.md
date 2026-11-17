# Setup and updates

## Install the CLI

```bash
curl -fsSL https://github.com/tekibo/idev/releases/latest/download/install.sh | sh
```

The installer verifies the release SHA-256 and writes `~/.local/bin/idev`.
`IDEV_INSTALL_DIR` changes the destination.

## Read the state first

```bash
idev next
idev next --json
```

`idev next` checks source pins, system packages, the private toolchain, the
Darwin SDK, Darwin tools/OAM, flows, and the project. It never changes anything.

`idev doctor --json` is the full report. Each check has `id`, `class`,
`status`, `detail`, `hint`, `package`, and `critical`.

## System packages

```bash
idev system plan
idev system plan --all
```

The plan detects `pacman`, `apt`, or `dnf` and prints one exact command. Show it
to the user and ask before running it.

```bash
idev system apply
```

`apply` runs the exact previewed transaction through the masked sudo prompt. It
never upgrades, never replaces, and never uses an AUR helper or repository.

`--all` includes optional packages such as `cmake`, `ninja`, `patchelf`, and
`fuse3` for SDK-free mode.

## Private toolchain

`$IDEV_HOME` defaults to `~/.local/share/idev`.

```bash
idev setup --base
```

Installs the shared payloads and components: apple-tools, Swift, xtool,
pymobiledevice3, rcodesign, and ipsw. Use the narrower steps when repairing:

```bash
idev setup --swift --dry-run
idev setup --swift
idev setup --xtool
idev setup --apple-tools
idev setup --pymobile
idev setup --cocoapods
idev setup --payloads
idev setup --rcodesign
idev setup --ipsw
idev setup --swift-build
```

Each step is independent. `--dry-run` prints the command without running it.

## Darwin SDK and Xcode.xip

The user downloads `Xcode.xip` from Apple. idev never does.

```bash
idev xcode list
idev xcode import /path/to/Xcode.xip
idev sdk install /path/to/Xcode.xip
idev sdk status --json
```

The Darwin SDK and Darwin tools are private. SwiftUI support comes from this
step. `idev setup --apple-tools` builds `actool` and OAM into the SDK.

## Flows

A flow selects the optional framework arm:

```bash
idev setup --flow reactnative          # Expo + React Native
idev setup --flow reactnative --no-expo # bare React Native
idev setup --flow flutter
idev setup --flow kmp
```

`expo` and `reactnative` share the React Native/Expo arm. `--no-expo` keeps bare
React Native and disables Expo. The selection is saved in
`$IDEV_HOME/config/flows.json`.

## Source pins and migration

```bash
idev config migrate --dry-run
idev config migrate
idev config migrate --yes
idev config refresh
```

`idev config migrate` merges `sources.json` against the last applied defaults.
Conflicting editable fields are confirmed; a changed required pin is reported
with its current and required values. Non-interactive use without `--yes` stops
instead of guessing. Every write creates a timestamped backup.

`idev config refresh` resets the user copy to the embedded defaults.

## Updates

```bash
idev update check
idev update check --json
idev update apply apple-tools
idev update apply pymobiledevice3
idev update apply --all
idev update clean --dry-run
idev update clean
```

Updates cover the manager plus every base, payload, and framework component.
`idev update check` is the source of truth. `idev update apply --all` updates
every outdated component and reconciles the manager from the release artifact,
then refreshes an installed agent skill. `idev update apply idev` updates only
the CLI and prints a restart reminder.

`idev update clean` shows the update downloads and build caches, asks before
removing them, and keeps the SDK, auth, current config and state, xips, and
installed payloads. Use `--dry-run` to inspect and `--yes` for non-interactive runs. The
manifest and assets are GitHub Release files:

```text
https://github.com/tekibo/idev/releases/latest/download/manifest.json
```

The manager update downloads the archive, verifies its SHA-256, and atomically
replaces the running binary. An unsupported required pin is shown as `unsupported`.

## Uninstall

```bash
idev uninstall --dry-run
idev uninstall
```

The report lists managed paths, system packages, and user configs. The CLI asks
for the exact phrase `remove idev`. Run `idev agents uninstall` first to remove
the agent instruction section.

`idev uninstall` is the kill switch; it removes only idev-managed paths and
lists the user configs it will not remove. Preview it with
`idev uninstall --dry-run`.
