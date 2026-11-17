# iPad screen and control

`idev ipad-screen` mirrors and controls a USB-attached iPad. It does not build
apps. It requires Tailscale on the host, a signed-in tailnet, and Serve
operator permission.

## Install

```bash
idev ipad-screen install
```

This writes:

- `$IDEV_HOME/bin/ipad-screen` — the command wrapper
- `$IDEV_HOME/config/ipad-screen/password` — mode 0600
- `$IDEV_HOME/state/ipad-screen/` — lifecycle state and logs
- `~/.local/share/applications/idev-ipad-screen.desktop`
- `~/.config/systemd/user/ipad-screen@.service`

## Start and stop

```bash
idev ipad-screen start --udid <udid>
idev ipad-screen start --udid <udid> --port 8080
idev ipad-screen status
idev ipad-screen stop
```

`start` checks Developer Mode, mounts the DDI with the rootless userspace
tunnel, binds the viewer to loopback, verifies the password, then publishes it
through `tailscale serve`. `status` prints the published HTTPS URL. `stop`
removes only its exact Serve mapping.

`run` is the foreground worker used by the user service. `menu` is the desktop
entry's interactive prompt. Use the service or manual start, not both.

## Prerequisites

1. Install Tailscale: `idev system plan` includes `tailscale`.
2. Start the daemon: `sudo systemctl enable --now tailscaled`.
3. Sign in: `tailscale up`.
4. Grant Serve permission if needed:
   `sudo tailscale set --operator=$USER`.
5. Connect the iPad by USB and enable Developer Mode and trust.
6. Restrict the published port in the tailnet ACL to the intended group.

The package install is automated; enabling and starting `tailscaled`,
`tailscale up`, ACL changes, and the iPad trust prompts are user actions.
`idev doctor` warns when the daemon socket is missing.

## Security

- The password file is private; never print its contents.
- The viewer binds loopback only; do not expose it on a LAN or public interface.
- The viewer is password-protected and Tailscale is the ingress boundary.
- If Serve is denied, stop and show the user the exact `tailscale set` command;
  do not change their Tailscale preferences.

## Files and validation

The config and state are under `$IDEV_HOME`; the password is mode 0600. The
install path and `status` are tested. A live display and Serve publication need
Tailscale and a USB iPad. The viewer is CLI-only.
