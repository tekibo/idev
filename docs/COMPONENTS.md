# Components

`idev` keeps its toolchain under `$IDEV_HOME` and uses the host package manager
for system packages.

## System packages

`idev system plan` prints one exact command for your package manager. The plan
covers USB device support, Tailscale, compilers, Python, Ruby, Node, and shared
libraries used by the toolchain. `idev system plan --all` adds the optional
SDK-free tools.

Review the command, then run `idev system apply` to execute it through the sudo
prompt.

## Private components

| Component | Purpose |
|---|---|
| `apple-tools` | SwiftUI, shipping, device-run, Darwin tools, SDK-free, simview, machoload |
| `Swift` | Official swift.org toolchain and SourceKit-LSP |
| `xtool` | Signing, provisioning, and device install |
| `pymobiledevice3` | Device transport and the iPad viewer |
| `rcodesign` | Code signing helper |
| `ipsw` | IPSW downloads and inspection |
| `Swift Build` | Build engine for SwiftUI and React Native |
| `CocoaPods` | Private gem tree and `pod` wrapper |
| `Flutter` | Flutter setup, build, and ship scripts |
| `KMP` | Kotlin Multiplatform framework link |
| `React Native` | idev-owned React Native and Expo engine |
| `Agent skill` | Embedded `idev` skill in `$IDEV_HOME/skills/idev` |

## Paths

- Swift: `$IDEV_HOME/toolchains/swift-6.4.2`
- xtool: `$IDEV_HOME/lib/xtool/<sha>`
- apple-tools: `$IDEV_HOME/apple-tools`
- Flutter, KMP, React Native: `$IDEV_HOME/flutter`, `kmp`, `react-native`
- pymobiledevice3: `$IDEV_HOME/pymobile3-venv`
- CocoaPods: `$IDEV_HOME/gems`
- Agent skill: `$IDEV_HOME/skills/idev`
- Xcode.xip: `$IDEV_HOME/xcode`

## Updates

```bash
idev update check
idev update apply --all
idev update apply idev
idev update clean
```

`idev update check` lists every component. `idev update apply <component>`
installs one; `--all` updates outdated components and reconciles the manager.
`idev update clean` removes leftover downloads and build caches.

## User-supplied

- `Xcode.xip`: downloaded from Apple.
- Apple ID and team: you sign in through xtool.
- Device trust and Developer Mode.
- Tailscale daemon started (`sudo systemctl enable --now tailscaled`), login
  (`tailscale up`), and Serve operator permission.
