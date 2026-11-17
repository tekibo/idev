# Agent usage

`idev` is a CLI; agents use the same commands as people.

## Start

```bash
idev next --json
idev doctor --json
```

`next` reports the setup state and the next command. `doctor` reports the full
state, including JSON fields for machine use.

## Install and update

```bash
curl -fsSL https://github.com/tekibo/idev/releases/latest/download/install.sh | sh
idev update check
idev update apply --all
```

`update apply --all` updates outdated components, reconciles the manager binary,
and refreshes the installed skill. Restart `idev` after the manager changes.
`idev update clean` removes leftover downloads and build caches after asking.

## Install the skill

```bash
idev agents install
idev agents status
```

The embedded skill is written to `$IDEV_HOME/skills/idev` and linked into
`~/.claude/skills/idev` and `~/.agents/skills/idev`. It covers setup, updates,
SwiftUI, Expo, Flutter, KMP, devices, shipping, the iPad viewer, and
troubleshooting.

## Rules

- Run `idev next` before changing the host.
- Show the `idev system plan` command and ask before installing packages.
- Never ask for an Apple ID password, 2FA code, API key, or token in chat.
- Long builds and tests run in the background.
- Keep project files and user credentials out of commands and logs.
