# Ship and App Store Connect

## SwiftUI

```bash
idev ship --device
idev ship --device TEAMID
idev ship --upload
```

`idev ship --device` builds a device build. `--upload` uploads after signing.
The App Store record must exist for the bundle ID before the first upload.

## Flutter

```bash
idev flutter ship
```

The Flutter ship path uses the private toolchain and the same Apple account
variables as `idev ship`.

## App Store Connect helper

```bash
idev asc stamp <app>
idev asc validate <app>
idev asc upload <app>
idev asc publish <app>
```

`stamp` writes the build-environment keys App Store processing reads.
`validate` checks offline. `upload` sends the build. `publish` performs the
post-upload release steps; it is opt-in and prints what it will do.

`idev asc` and `idev provision-dev` run through the private pymobiledevice3 venv,
so `cryptography` is available. `idev auth`/`idev ship` use the xtool account.

## Development profile and signing

```bash
idev provision-dev --bundle-id com.example.app --udid <udid>
idev sign-dev /path/App.app
```

`provision-dev` creates a development profile. `sign-dev` signs an app with a
development identity. xtool remains the production signer.

## Credentials

- Never ask for a password, 2FA code, app-specific password, or API key in chat.
- The user stores App Store Connect API values in their own environment or
  `$IDEV_HOME` config, not in the repo.
- The signing identity, provisioning profile, and App Store Connect key are the
  user's. Do not copy them into logs, fixtures, or commits.

## Rules

- Do not submit for App Store review or release to users without explicit
  approval.
- TestFlight upload is fine when the user asks for a build.
- `idev ship --upload` is a network and account action; report the exact failure
  without echoing credentials.
- After a failed upload, check `idev asc validate` before retrying.
