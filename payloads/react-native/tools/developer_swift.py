#!/usr/bin/env python3
"""Install the Swift frontend wrapper in the build overlay."""

from __future__ import annotations

import os
import shlex

WRAPPER = """#!/bin/sh
header=
previous=
for argument in "$@"; do
    if [ "$previous" = "-emit-objc-header-path" ]; then
        header=$argument
        break
    fi
    previous=$argument
done
{real} "$@"
status=$?
if [ "$status" -eq 0 ] && [ -n "$header" ] && [ -f "$header" ]; then
    python3 {strip} "$header" >/dev/null 2>&1 || true
fi
exit "$status"
"""


def overlay_swift_frontend(real_bin: str, overlay_bin: str, strip_script: str) -> None:
    real = os.path.join(real_bin, "swift-frontend")
    target = os.path.join(overlay_bin, "swift-frontend")
    if os.path.lexists(target):
        os.remove(target)
    with open(target, "w", encoding="utf-8") as handle:
        handle.write(
            WRAPPER.format(
                real=shlex.quote(real),
                strip=shlex.quote(strip_script),
            )
        )
    os.chmod(target, 0o755)
