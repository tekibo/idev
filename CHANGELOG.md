# Changelog

## 0.1.0

- First public release.
- `idev next` reports the setup state and the exact next command, with `--json`.
- `idev setup` installs the base toolchain and the optional React Native/Expo,
  Flutter, and KMP flows.
- `idev doctor` groups checks by class, prints `ok`/`WARN`/`FAIL`, and ends with
  the next command.
- `idev agents install|status|uninstall` manages the embedded skill for Claude
  Code and Codex.
- `idev update check` and `idev update apply` use GitHub Releases. `apply --all`
  updates every component, reconciles the manager artifact, and refreshes an
  installed agent skill.
- `idev update clean` shows update downloads and build caches, asks before
  removing them, and keeps the SDK, auth, config, state, xips, and payloads.
- `idev xcode` and `idev sdk` manage `Xcode.xip` and the private Darwin SDK.
- `idev devices`, `idev install-app`, `idev launch-app`, and
  `idev uninstall-app` cover the device loop.
- `idev ipad-screen` mirrors and controls a USB iPad over Tailscale. Its docs
  and `doctor` state that installing `tailscale` does not start the daemon;
  they show the daemon, login, and Serve operator commands.
- `idev lsp` starts SourceKit-LSP in the LSP workspace root without the
  duplicate-toolchain `DEVELOPER_DIR`, fixing editor timeouts in monorepos.
- Pins: apple-dev `d0f4980`, xtool `7cc56c3`, pymobiledevice3 `11.58.0`.
