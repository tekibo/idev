# `idev doctor`

`idev doctor` inspects the host, the private toolchain, and the current project.
It is read-only.

## Output

Text output groups checks by class and prints `ok`, `WARN`, or `FAIL`, followed
by the next command. `idev doctor --json` is the stable machine-readable form.
Each check includes `id`, `class`, `status`, `fixable_by`, `detail`, `hint`,
`package`, and `critical`.

## What it checks

| Area | Examples |
|---|---|
| System | USB support, Tailscale, compiler and runtime packages |
| Private toolchain | Swift, xtool, pymobiledevice3, CocoaPods, rcodesign, ipsw |
| Apple tooling | Darwin SDK, `actool`, OAM |
| Project flows | Flutter, KMP, React Native and Expo |
| iPad viewer | Tailscale daemon, viewer script, device transport |
| Project | project kind, bundle ID, and build files |

## System plan

`idev system plan` prints one package command for the detected manager:

| Manager | Command |
|---|---|
| pacman | `sudo pacman -S --needed --noconfirm <packages>` |
| apt | `sudo apt install -y <packages>` |
| dnf | `sudo dnf install -y <packages>` |

`idev system apply` runs the command through the sudo prompt and re-checks the
host.

## Exit codes

`0` when no critical check is missing; `1` when a critical check fails. `WARN`
lines are advice.

## Next

`idev next` is the shorter setup view: it prints only the steps and the one
command to run next. Use it before retrying a failed command.
