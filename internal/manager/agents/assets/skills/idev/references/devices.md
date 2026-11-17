# Devices and signing

## List

```bash
idev devices
idev devices --json
```

The JSON form is:

```json
[{"name":"iPad","connection":"usb","udid":"00008103-..."}]
```

Use the UDID explicitly when more than one device is connected. Do not assume the
first device is the user's target.

## Install, launch, remove

```bash
idev install-app --bundle /path/App.app --bundle-id com.example.app
idev install-app --bundle /path/App.ipa --bundle-id com.example.app --udid <udid>
idev launch-app --bundle-id com.example.app --udid <udid>
idev uninstall-app --bundle-id com.example.app --udid <udid>
```

`install-app` signs with xtool, installs, and launches. `--dev-server-url` sets
the initial dev-server URL for a dev client. `launch-app` resolves the signed
installed identifier; `XTL-<team>.<id>` is that value, not a new App ID.

## Device run

```bash
idev device-run
idev device-run --console 30
idev device-run --lldb
idev device-run --attach
```

`device-run` is the apple-tools helper. It needs the Darwin SDK and the private
pymobiledevice3 venv.

## Account

```bash
idev auth status --json
idev auth login
```

`status` prints the parsed account state; `login` passes through to xtool. The
token never displays. Never ask the user for their Apple ID password or 2FA code
in chat.

## Rules

- xtool is the only signer and provisioner.
- A free Apple ID signs a 7-day development build; a paid team is required for
  distribution, push, and TestFlight.
- Do not create or change a bundle ID without explicit approval.
- If the device does not appear, check USB trust and Developer Mode before
  changing anything.
- `idev doctor` reports device and auth checks; `idev next` shows the setup step
  that is missing.
