# Install and setup

`idev` installs a private iOS toolchain on Linux and wraps xtool, the Darwin
tools, and your project's build system.

## Install the CLI

```bash
curl -fsSL https://github.com/tekibo/idev/releases/latest/download/install.sh | sh
```

The installer verifies the release SHA-256 and writes `~/.local/bin/idev`.
`IDEV_INSTALL_DIR` changes the destination.

```bash
idev next
```

`idev next` is read-only. It prints the setup steps and the one command to run
next. `idev next --json` is the machine-readable form.

## System packages

```bash
idev system plan
idev system plan --all
idev system apply
```

`idev system plan` detects `pacman`, `apt`, or `dnf` and prints one exact
command. Review it, then run `idev system apply`; it executes that command
through the sudo prompt.

`--all` includes the optional SDK-free packages.

## Private toolchain

`$IDEV_HOME` defaults to `~/.local/share/idev`.

```bash
idev setup --base
```

This installs the shared toolchain:

- apple-tools
- Swift
- xtool
- pymobiledevice3
- rcodesign
- ipsw

Repair one component at a time:

```bash
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

`--dry-run` prints the command without running it.

## Xcode and the Darwin SDK

Download `Xcode.xip` from Apple, then import it and install the SDK:

```bash
idev xcode list
idev xcode import /path/to/Xcode.xip
idev sdk install /path/to/Xcode.xip
idev sdk status --json
```

The Darwin SDK and Darwin tools are private. SwiftUI support comes from this
step. `idev setup --apple-tools` builds `actool` and OAM into the SDK.

## Flows

```bash
idev setup --flow reactnative          # Expo + React Native
idev setup --flow reactnative --no-expo # bare React Native
idev setup --flow flutter
idev setup --flow kmp
```

The selection is saved in `$IDEV_HOME/config/flows.json`. SwiftUI uses the base
toolchain and needs no flow.

## iPad viewer

```bash
idev setup --apple-tools
idev ipad-screen install
sudo systemctl enable --now tailscaled
tailscale up
sudo tailscale set --operator=$USER
idev ipad-screen start --udid <udid>
```

Installing the `tailscale` package does not start the daemon. `tailscale up`
signs in, and Tailscale Serve needs operator permission. The iPad needs USB
trust and Developer Mode. `idev doctor` warns when the Tailscale daemon socket
is missing.

## Agent skill

```bash
idev agents install
idev agents status
```

The embedded skill is written under `$IDEV_HOME/skills/idev` and linked into
`~/.claude/skills` and `~/.agents/skills`. It covers install, update, cleanup,
setup, `next`, `doctor`, project flows, devices, shipping, the iPad viewer, and
uninstall.

## Updates

```bash
idev update check
idev update apply --all
idev update apply idev
idev update clean --dry-run
idev update clean
```

The manifest and assets are GitHub Release files:

```text
https://github.com/tekibo/idev/releases/latest/download/manifest.json
```

`apply` verifies the release SHA-256. `apply --all` updates outdated components,
reconciles the manager binary, and refreshes an installed agent skill. Restart
`idev` afterwards.

`update clean` removes update downloads and build caches after showing the total
and asking for confirmation. `--dry-run` shows the plan; `--yes` skips the
prompt. The SDK, auth, config, state, xips, and installed payloads stay.

## Configuration

```bash
idev config migrate --dry-run
idev config migrate
idev config refresh
```

`migrate` merges `sources.json` with the latest defaults and confirms conflicting
editable fields. Each write creates a timestamped backup. `refresh` resets the
user copy to the embedded defaults.

## Verification

```bash
idev doctor
idev doctor --json
idev devices --json
```

## Uninstall

```bash
idev uninstall --dry-run
idev uninstall
```

See [UNINSTALL.md](UNINSTALL.md).
