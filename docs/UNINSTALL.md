# Uninstall

`idev uninstall` removes the paths idev manages. System packages, Apple
credentials, device pairing records, and project files stay in place.

## Preview

```bash
idev uninstall --dry-run
idev uninstall --dry-run --json
idev uninstall --keep swift --keep xtool
```

The report lists every managed component, its exact paths, the system packages
idev installed, and the user configs it will not remove.

## Apply

```bash
idev uninstall
idev uninstall --keep pymobiledevice3
```

The CLI asks for the exact confirmation phrase `remove idev`. The latest report
is saved under `$IDEV_HOME/state`.

## Keepable tools

| ID | Tool |
|---|---|
| `swift` | Swift toolchain and SourceKit-LSP |
| `xtool` | xtool signer |
| `pymobiledevice3` | pymobiledevice3 and its private venv |
| `cocoapods` | CocoaPods gem tree |
| `rcodesign` | rcodesign |
| `ipsw` | ipsw |

Keeping xtool keeps Swift. A partial uninstall leaves the kept tools in place;
idev will no longer update them.

## Removed

- `$IDEV_HOME/apple-tools`, `flutter`, `kmp`, `react-native`
- `$IDEV_HOME/skills/idev` and the links in `~/.claude/skills` and `~/.agents/skills`
- `$IDEV_HOME/bin/ipad-screen`, its desktop entry, and its user unit
- the private Darwin SDK
- Swift Build and the React Native engine
- `$IDEV_HOME/config`, `state`, and `cache`
- `~/.local/bin/idev` when it is the installed binary
- `~/.local/bin/xtool` when it is the idev shim
- `~/.local/bin/pymobiledevice3` when it is the idev link

## Kept

- system packages
- `$IDEV_HOME/xcode` and `Xcode.xip`
- `~/.local/share/xtool` and `~/.config/xtool`
- `~/.pymobiledevice3` and `/var/lib/lockdown`
- `~/.appstoreconnect`
- the `<!-- idev:begin -->` section in agent instructions; run
  `idev agents uninstall` first to remove it
- project files, `node_modules`, `ios/Pods`, and build output
