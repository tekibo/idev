"""Private toolchain paths shared by the Expo arm helpers."""

from __future__ import annotations

import os


def swift_bin() -> str:
    return os.environ.get("IDEV_SWIFT_BIN", "")


def swift_include() -> str:
    value = os.environ.get("IDEV_SWIFT_INTEROP_INCLUDE", "")
    if value:
        return value
    binary = swift_bin()
    if not binary:
        return ""
    return os.path.normpath(os.path.join(binary, "..", "include"))


def apple_tools_root() -> str:
    return os.environ.get("IDEV_APPLE_TOOLS", os.environ.get("I_DEV_APPLE_TOOLS", ""))


def sdk_bundle() -> str:
    return os.environ.get("DARWIN_SDK_BUNDLE", "")
