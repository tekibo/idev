# Troubleshooting

## Order

1. `idev next` — the setup state and the exact next command.
2. `idev doctor` — the full report.
3. Only then retry the failed command.

Do not fix a diagnostic by guessing. The code at hand may be unnecessary; walk
back to the project state and the setup step first.

## Common failures

| Symptom | First check | Fix |
|---|---|---|
| `idev setup` stops on a source pin | `idev config migrate --dry-run` | Confirm the required value, then `idev config migrate` |
| System package missing | `idev system plan` | Show the exact command; ask before `idev system apply` |
| `xtool` missing | `idev doctor --json` | `idev setup --xtool` |
| Darwin SDK missing | `idev sdk status --json` | `idev sdk install /path/Xcode.xip` |
| `actool` missing | `idev doctor` (`private.apple_tools`) | `idev setup --apple-tools` |
| Device not listed | `idev devices --json` | USB trust and Developer Mode |
| App installs but does not launch | `idev auth status --json` | Sign in through xtool; check the team |
| Expo/RN build fails after a native change | `idev expo run:ios` log | Rebuild; Fast Refresh only covers JS |
| Flutter flow missing | `idev next` | `idev setup --flow flutter` |
| KMP link fails on a missing host archive | `idev kmp` log | Let the pinned download finish; check `KMP_X64_SUMS` if pinning another version |
| iPad viewer cannot publish | `idev doctor` (`private.tailscale`) | `sudo systemctl enable --now tailscaled`, `tailscale up`, and Serve operator permission |
| Update shows `unsupported` | `idev update check` | Follow the pin warning; `idev config migrate` |
| Agent skill is stale | `idev agents status --json` (`stale`) | `idev agents install` |
| Disk usage from updates and builds | `idev update clean --dry-run` | `idev update clean` |

## System plan

`idev system plan` prints the detected manager and one exact command. The adapter
refuses upgrades and replacements and never uses an AUR helper or repository.

## Updates and pins

`idev update check` compares the manager and every component with the release
manifest. `idev update apply <component>` installs one. `idev update apply
--all` also reconciles the manager artifact and refreshes an installed agent
skill. An unsupported required pin is shown as `unsupported`; `idev config
migrate` resolves it.

`idev update clean` removes only update downloads, build caches, removed-feature
logs and config, and asks before deleting. It never touches the SDK, auth,
config state, xips, or installed payloads.

## When a command is unknown

Run `idev help`, then `idev help <command>`, then
`idev <command> --help`. Do not invent flags; the help tree is the contract.

## Secrets

Never ask for a password, 2FA code, API key, or token in chat. The user sets
secrets in their own environment or config. If a secret appears in output, stop
and redact it before continuing.
