"""Emit a PIF document from the Xcode project model.

PIF is what Swift Build consumes: a JSON array of

    { "signature": "PROJECT@v11_P1", "type": "project", "contents": { ... } }

Key names and object shapes are copied from upstream
`SWBProtocol/PIFKeyConstants.swift` and
`SWBProjectModel/IDE/IDESwiftPackageExtensions.swift`. The loader stores every
object as plain JSON (`SWBCore/ProjectModel/PIFLoader.swift`), so no
MessagePack is involved.

This module only translates. It makes no decision about how anything compiles.
"""

from __future__ import annotations

import json
import os

FILE_TYPES = {
    ".swift": "sourcecode.swift",
    ".m": "sourcecode.c.objc",
    ".mm": "sourcecode.cpp.objcpp",
    ".c": "sourcecode.c.c",
    ".cc": "sourcecode.cpp.cpp",
    ".cpp": "sourcecode.cpp.cpp",
    ".h": "sourcecode.c.h",
    ".hpp": "sourcecode.cpp.h",
    ".hh": "sourcecode.cpp.h",
    ".plist": "text.plist.xml",
    ".json": "text.json",
    ".js": "sourcecode.javascript",
    ".png": "image.png",
    ".jpg": "image.jpeg",
    ".jpeg": "image.jpeg",
    ".xcassets": "folder.assetcatalog",
    ".storyboard": "file.storyboard",
    ".xib": "file.xib",
    ".strings": "text.plist.strings",
    ".stringsdict": "text.plist.stringsdict",
    ".xcstrings": "text.json.xcstrings",
    ".framework": "wrapper.framework",
    ".xcframework": "wrapper.xcframework",
    ".a": "archive.ar",
    ".dylib": "compiled.mach-o.dylib",
    ".bundle": "wrapper.plug-in",
    ".xcprivacy": "text.xml",
    ".entitlements": "text.plist.entitlements",
    ".xcconfig": "text.xcconfig",
    ".modulemap": "sourcecode.modulemap",
    ".ttf": "file",
    ".otf": "file",
}

PHASE_TYPES = {
    "PBXSourcesBuildPhase": "com.apple.buildphase.sources",
    "PBXResourcesBuildPhase": "com.apple.buildphase.resources",
    "PBXFrameworksBuildPhase": "com.apple.buildphase.frameworks",
    "PBXHeadersBuildPhase": "com.apple.buildphase.headers",
    "PBXCopyFilesBuildPhase": "com.apple.buildphase.copy-files",
    "PBXShellScriptBuildPhase": "com.apple.buildphase.shell-script",
}


def file_type_for(path: str, recorded: str = "") -> str:
    if recorded:
        return recorded
    return FILE_TYPES.get(os.path.splitext(path)[1].lower(), "file")


class Pif:
    """Accumulates top-level PIF objects and serializes them."""

    def __init__(self) -> None:
        self.objects: list[dict] = []

    def add(self, kind: str, signature: str, contents: dict) -> str:
        self.objects.append({"signature": signature, "type": kind, "contents": contents})
        return signature

    def dumps(self) -> str:
        return json.dumps(self.objects, indent=1)

    # -- top-level objects ----------------------------------------------

    def workspace(self, guid: str, name: str, path: str, projects: list[str]) -> str:
        return self.add(
            "workspace",
            f"WORKSPACE@v11_{guid}",
            {"guid": guid, "name": name, "path": path, "projects": projects},
        )

    def project(
        self,
        guid: str,
        name: str,
        project_path: str,
        project_dir: str,
        group_tree: dict,
        build_configurations: list[dict],
        targets: list[str],
        default_configuration: str,
    ) -> str:
        return self.add(
            "project",
            f"PROJECT@v11_{guid}",
            {
                "guid": guid,
                "projectName": name,
                "projectIsPackage": "false",
                "path": project_path,
                "projectDirectory": project_dir,
                "groupTree": group_tree,
                "buildConfigurations": build_configurations,
                "defaultConfigurationName": default_configuration,
                "targets": targets,
                "developmentRegion": "en",
                "knownRegions": ["en", "Base"],
            },
        )

    def aggregate_target(
        self,
        guid: str,
        name: str,
        build_configurations: list[dict],
        build_phases: list[dict],
        dependencies: list[str],
    ) -> str:
        return self.add(
            "target",
            f"TARGET@v11_{guid}",
            {
                "type": "aggregate",
                "guid": guid,
                "name": name,
                "dependencies": [
                    {"guid": dependency, "platformFilters": [], "buildConfigurationFilters": []}
                    for dependency in dependencies
                ],
                "buildRules": [],
                "customTasks": [],
                "buildPhases": build_phases,
                "buildConfigurations": build_configurations,
            },
        )

    def target(
        self,
        guid: str,
        name: str,
        product_name: str,
        product_type: str,
        product_reference_guid: str,
        build_configurations: list[dict],
        build_phases: list[dict],
        dependencies: list[str],
    ) -> str:
        return self.add(
            "target",
            f"TARGET@v11_{guid}",
            {
                "type": "standard",
                "guid": guid,
                "name": name,
                "productTypeIdentifier": product_type,
                "productReference": {
                    "type": "file",
                    "guid": product_reference_guid,
                    "name": product_name,
                },
                "dependencies": [
                    {"guid": dependency, "platformFilters": [], "buildConfigurationFilters": []}
                    for dependency in dependencies
                ],
                "buildRules": [],
                "customTasks": [],
                "buildPhases": build_phases,
                "buildConfigurations": build_configurations,
                "approvedByUser": "true",
            },
        )


# -- inline objects -----------------------------------------------------


def build_configuration(guid: str, name: str, settings: dict[str, str]) -> dict:
    return {
        "guid": guid,
        "name": name,
        "buildSettings": dict(settings),
        "impartedBuildProperties": {"buildSettings": {}},
    }


def build_file(guid: str, file_guid: str) -> dict:
    return {
        "guid": guid,
        "fileReference": file_guid,
        "platformFilters": [],
        "codeSignOnCopy": "false",
        "removeHeadersOnCopy": "false",
    }


def group(guid: str, name: str, path: str, source_tree: str, children: list[dict]) -> dict:
    return {
        "type": "group",
        "guid": guid,
        "name": name,
        "path": path,
        "sourceTree": source_tree,
        "children": children,
    }


def file_reference(guid: str, name: str, path: str, source_tree: str, file_type: str) -> dict:
    return {
        "type": "file",
        "guid": guid,
        "name": name,
        "path": path,
        "sourceTree": source_tree,
        "fileType": file_type,
    }


def build_phase(phase, guid: str, files: list[dict]) -> dict:
    """Translate one Xcode build phase into its PIF form."""
    kind = PHASE_TYPES.get(phase.isa, "com.apple.buildphase.sources")
    if phase.isa == "PBXShellScriptBuildPhase":
        return {
            "type": kind,
            "guid": guid,
            "name": phase.name,
            "shellPath": "/bin/sh",
            # Xcode runs script phases from the project root; Swift Build does
            # not, and CocoaPods scripts use paths relative to it.
            "scriptContents": 'cd "${SRCROOT}" || exit 1\n' + phase.shell_script,
            "buildFiles": files,
            "inputFilePaths": list(phase.input_paths),
            "inputFileListPaths": list(phase.input_file_lists),
            "outputFilePaths": list(phase.output_paths),
            "outputFileListPaths": list(phase.output_file_lists),
            "emitEnvironment": "false",
            "sandboxingOverride": "basedOnBuildSetting",
            "alwaysRunForInstallHdrs": "false",
            # Xcode runs script phases every build; Swift Build otherwise
            # treats a phase without declared outputs as up to date.
            "alwaysOutOfDate": "true",
            "runOnlyForDeploymentPostprocessing": (
                "true" if phase.run_only_for_deployment_postprocessing else "false"
            ),
            "originalObjectID": phase.id,
        }
    return {"type": kind, "guid": guid, "buildFiles": files}
