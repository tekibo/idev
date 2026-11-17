#!/bin/sh
# Install or update the idev CLI from a GitHub release.
set -eu

REPO="${IDEV_REPO:-tekibo/idev}"
BASE="${IDEV_RELEASE_BASE:-https://github.com/$REPO/releases/latest/download}"
INSTALL_DIR="${IDEV_INSTALL_DIR:-$HOME/.local/bin}"

case "$(uname -m)" in
	x86_64 | amd64) arch=amd64 ;;
	aarch64 | arm64) arch=arm64 ;;
	*) echo "install: unsupported architecture $(uname -m)" >&2; exit 1 ;;
esac

asset="idev-linux-$arch.tar.gz"
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT HUP INT TERM
cd "$tmp"

echo "install: $BASE/$asset"
curl -fsSL "$BASE/SHA256SUMS" -o SHA256SUMS
curl -fsSL "$BASE/$asset" -o "$asset"
sum="$(grep "  $asset\$" SHA256SUMS || true)"
test -n "$sum" || { echo "install: $asset is not in SHA256SUMS" >&2; exit 1; }
printf '%s\n' "$sum" | sha256sum -c -
tar -xzf "$asset"

test -x idev || { echo "install: archive has no idev binary" >&2; exit 1; }
mkdir -p "$INSTALL_DIR"
target="$INSTALL_DIR/idev"
cp idev "$target.tmp"
chmod 755 "$target.tmp"
mv "$target.tmp" "$target"

echo "installed: $target"
"$target" --version
echo "next: idev next"
echo "agent skill: idev agents install"
