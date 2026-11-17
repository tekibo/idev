#!/usr/bin/env bash
# Mirror/control only. Never builds, installs packages system-wide, or uses sudo.
set -euo pipefail
root=$(cd -- "$(dirname -- "$(readlink -f -- "$0")")" && pwd)
idev=$(cd -- "$root/../.." && pwd)
if [[ -n ${IPAD_SCREEN_PYTHON:-} ]]; then
  py=$IPAD_SCREEN_PYTHON
elif [[ -x ${IDEV_HOME:-$idev}/pymobile3-venv/bin/python ]]; then
  py=${IDEV_HOME:-$idev}/pymobile3-venv/bin/python
elif [[ -x $HOME/pymobile3-venv/bin/python ]]; then
  py=$HOME/pymobile3-venv/bin/python
elif [[ -x $HOME/.local/share/pipx/venvs/pymobiledevice3/bin/python ]]; then
  py=$HOME/.local/share/pipx/venvs/pymobiledevice3/bin/python
else
  echo 'ipad-screen: need the idev pymobiledevice3 venv or a pymobiledevice3 pipx environment' >&2
  exit 1
fi
: "${IPAD_SCREEN_CONFIG:=${IDEV_HOME:-$idev}/config/ipad-screen}"
: "${IPAD_SCREEN_STATE:=${IDEV_HOME:-$idev}/state/ipad-screen}"
export IPAD_SCREEN_CONFIG IPAD_SCREEN_STATE
exec "$py" "$root/ipad-screen.py" "$@"
