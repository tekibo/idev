package cocoapods

import (
	"os"
	"path/filepath"
)

const wrapper = `#!/bin/sh
set -eu
SELF=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
ROOT=$(dirname -- "$SELF")
GEM_ROOT="$ROOT/gems"
export LANG=C.UTF-8
export LC_ALL=C.UTF-8
export GEM_HOME="$GEM_ROOT"
export GEM_PATH="$GEM_ROOT"
export PATH="$GEM_ROOT/bin:$ROOT/react-native/bin:$ROOT/bin:$PATH"
export COCOAPODS_DISABLE_STATS=1
export CP_HOME_DIR="${CP_HOME_DIR:-$ROOT/cache/cocoapods}"
APP_ROOT="$PWD"
if [ -f "$APP_ROOT/../app.json" ]; then
    APP_ROOT="$(cd "$APP_ROOT/.." && pwd)"
fi
PATCH="$ROOT/react-native/expo/bin/patch-expo-autolinking.py"
if [ -x "$PATCH" ] && [ -d "$APP_ROOT/node_modules" ]; then
    python3 "$PATCH" "$APP_ROOT" || true
fi
exec "$GEM_ROOT/bin/pod" "$@"
`

func writeWrapper(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	if err := os.WriteFile(path, []byte(wrapper), 0755); err != nil {
		return err
	}
	return os.Chmod(path, 0755)
}
