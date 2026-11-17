#!/usr/bin/env python3
"""Point Expo's precompiled-modules extraction at bsdtar."""

from __future__ import annotations

import argparse
import glob
import os
import sys


def find_targets(app_root: str) -> list[str]:
    patterns = [
        os.path.join(
            app_root, "node_modules", ".pnpm", "expo-modules-autolinking@*",
            "node_modules", "expo-modules-autolinking", "scripts", "ios",
            "precompiled_modules.rb",
        ),
        os.path.join(
            app_root, "node_modules", "expo-modules-autolinking", "scripts",
            "ios", "precompiled_modules.rb",
        ),
    ]
    targets = []
    for pattern in patterns:
        targets.extend(glob.glob(pattern))
    return sorted(set(targets))


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("app_root", help="Expo app root containing node_modules")
    parser.add_argument("--check", action="store_true", help="report without writing")
    args = parser.parse_args()

    targets = find_targets(os.path.abspath(args.app_root))
    if not targets:
        print("No expo-modules-autolinking precompiled_modules.rb found.", file=sys.stderr)
        return 1

    changed = 0
    for path in targets:
        with open(path, "r", encoding="utf-8") as handle:
            text = handle.read()
        if "'/usr/bin/bsdtar'" in text:
            print(f"already patched: {path}")
            continue
        if "'/usr/bin/tar'" not in text:
            print(f"no /usr/bin/tar call found: {path}", file=sys.stderr)
            continue
        if args.check:
            print(f"needs patch: {path}")
            changed += 1
            continue
        text = text.replace("'/usr/bin/tar'", "'/usr/bin/bsdtar'")
        with open(path, "w", encoding="utf-8") as handle:
            handle.write(text)
        print(f"patched: {path}")
        changed += 1

    if args.check and changed:
        return 2
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
