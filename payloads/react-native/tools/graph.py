"""Dump the Xcode project graph the way `expo/` sees it.

    python3 expo/tools/graph.py <ios-dir> [--configuration Release] [--json]

This is the debugging surface for the backend: if a target, build phase or
setting looks wrong here, it is wrong before any compilation starts.
"""

from __future__ import annotations

import argparse
import json
import os
import sys

sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from lib import xcodeproj  # noqa: E402


def describe(project: xcodeproj.Project, configuration: str, verbose: bool) -> dict:
    targets = []
    for target in project.targets:
        settings = target.settings(configuration)
        entry = {
            "name": target.name,
            "type": target.product_type.rsplit(".", 1)[-1] if target.product_type else "unknown",
            "sources": len(target.source_files()),
            "resources": len(target.resource_files()),
            "scriptPhases": [
                {
                    "name": phase.name,
                    "shell": phase.shell_script.strip().splitlines()[:3],
                }
                for phase in target.script_phases()
            ],
            "dependencies": [dependency.name for dependency in target.dependencies],
            "productName": settings.raw("PRODUCT_NAME"),
            "moduleName": settings.raw("PRODUCT_MODULE_NAME") or settings.raw("PRODUCT_NAME"),
            "deploymentTarget": settings.raw("IPHONEOS_DEPLOYMENT_TARGET"),
            "swiftVersion": settings.raw("SWIFT_VERSION"),
            "cxxInterop": settings.raw("SWIFT_CPP_INTEROP_MODE") or settings.raw("SWIFT_ENABLE_CXX_INTEROP"),
            "modules": settings.raw("CLANG_ENABLE_MODULES"),
            "definesModule": settings.raw("DEFINES_MODULE"),
        }
        if verbose:
            entry["sourcesList"] = [
                file.absolute_path or file.path for file in target.source_files()
            ]
            entry["phases"] = [phase.isa for phase in target.build_phases]
        targets.append(entry)
    return {"project": project.path, "configuration": configuration, "targets": targets}


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("ios_dir")
    parser.add_argument("--configuration", default="Release")
    parser.add_argument("--json", action="store_true")
    parser.add_argument("--verbose", action="store_true")
    args = parser.parse_args()

    workspace = xcodeproj.Workspace.discover(args.ios_dir)
    paths = list(workspace.project_paths) if workspace else []
    if not paths:
        paths = sorted(
            os.path.join(args.ios_dir, name)
            for name in os.listdir(args.ios_dir)
            if name.endswith(".xcodeproj")
        )
    if not paths:
        print(f"error: no .xcworkspace or .xcodeproj under {args.ios_dir}", file=sys.stderr)
        return 1

    report = [describe(xcodeproj.Project(path), args.configuration, args.verbose) for path in paths]
    if args.json:
        json.dump(report, sys.stdout, indent=2)
        print()
        return 0

    for project in report:
        print(f"\n== {project['project']} ({project['configuration']}) ==")
        for target in project["targets"]:
            print(
                f"  {target['name']:<28} {target['type']:<16} "
                f"src={target['sources']:<4} res={target['resources']:<3} "
                f"module={target['moduleName'] or '-':<24} "
                f"ios={target['deploymentTarget'] or '-'}"
            )
            for phase in target["scriptPhases"]:
                print(f"      script: {phase['name'] or '(unnamed)'}")
            if target["dependencies"]:
                print(f"      deps: {', '.join(target['dependencies'])}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
