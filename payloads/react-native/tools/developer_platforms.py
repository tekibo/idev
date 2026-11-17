#!/usr/bin/env python3
"""Mirror the Darwin Platforms tree with the actool shim in place."""

from __future__ import annotations

import os
import shlex

PLATFORM = "iPhoneOS.platform"
SDK_SHIMS = ("actool", "ibtool")


def link(source: str, target: str) -> None:
    if not os.path.lexists(target):
        os.symlink(source, target)


def wrapper(path: str, shim: str) -> None:
    with open(path, "w", encoding="utf-8") as handle:
        handle.write("#!/bin/sh\nexec " + shlex.quote(shim) + ' "$@"\n')
    os.chmod(path, 0o755)


def mirror(source_dir: str, target_dir: str, tools_dir: str) -> None:
    os.makedirs(target_dir, exist_ok=True)
    for entry in os.listdir(source_dir):
        target = os.path.join(target_dir, entry)
        shim = os.path.join(tools_dir, entry)
        if entry in SDK_SHIMS and os.path.isfile(shim):
            wrapper(target, shim)
        else:
            link(os.path.join(source_dir, entry), target)


def overlay_platforms(real: str, output: str, tools_dir: str) -> None:
    source = os.path.join(real, "Platforms")
    target = os.path.join(output, "Platforms")
    os.makedirs(target, exist_ok=True)
    for entry in os.listdir(source):
        if entry != PLATFORM:
            link(os.path.join(source, entry), os.path.join(target, entry))
            continue
        platform = os.path.join(source, entry)
        platform_target = os.path.join(target, entry)
        os.makedirs(platform_target, exist_ok=True)
        for child in os.listdir(platform):
            if child != "Developer":
                link(os.path.join(platform, child), os.path.join(platform_target, child))
                continue
            developer = os.path.join(platform, child)
            developer_target = os.path.join(platform_target, child)
            os.makedirs(developer_target, exist_ok=True)
            for item in os.listdir(developer):
                if item != "usr":
                    link(os.path.join(developer, item), os.path.join(developer_target, item))
                    continue
                usr = os.path.join(developer, item)
                usr_target = os.path.join(developer_target, item)
                os.makedirs(usr_target, exist_ok=True)
                for part in os.listdir(usr):
                    if part != "bin":
                        link(os.path.join(usr, part), os.path.join(usr_target, part))
                    else:
                        mirror(os.path.join(usr, part), os.path.join(usr_target, part), tools_dir)
