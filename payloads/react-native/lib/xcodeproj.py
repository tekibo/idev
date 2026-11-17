"""Xcode project and workspace model.

This is the layer the rest of `expo/` is built on: it turns `project.pbxproj`
into targets, build phases, source files and settings. It deliberately knows
nothing about CocoaPods, Expo or React Native — those are conventions that sit
on top of this model.
"""

from __future__ import annotations

import os
from dataclasses import dataclass, field
from typing import Iterator

from . import openstep, xcconfig
from .pif import file_type_for


def pbx_bool(value) -> bool:
    """pbxproj booleans are bare tokens: "0"/"1", sometimes YES/NO."""
    if isinstance(value, bool):
        return value
    if isinstance(value, str):
        return value.strip().upper() in ("1", "YES", "TRUE")
    return bool(value)

# Build phases Xcode applies, in the order they matter to us.
PHASE_TYPES = {
    "PBXSourcesBuildPhase",
    "PBXResourcesBuildPhase",
    "PBXFrameworksBuildPhase",
    "PBXHeadersBuildPhase",
    "PBXCopyFilesBuildPhase",
    "PBXShellScriptBuildPhase",
    "PBXRunScriptBuildPhase",
}


@dataclass
class FileRef:
    id: str
    path: str
    name: str = ""
    source_tree: str = "<group>"
    file_type: str = ""
    absolute_path: str = ""

    @property
    def extension(self) -> str:
        return os.path.splitext(self.path)[1].lower()


@dataclass
class BuildFile:
    id: str
    file_ref: FileRef | None
    settings: dict = field(default_factory=dict)


@dataclass
class BuildPhase:
    id: str
    isa: str
    files: list[BuildFile] = field(default_factory=list)
    name: str = ""
    shell_script: str = ""
    input_paths: list[str] = field(default_factory=list)
    output_paths: list[str] = field(default_factory=list)
    input_file_lists: list[str] = field(default_factory=list)
    output_file_lists: list[str] = field(default_factory=list)
    always_out_of_date: bool = False
    run_only_for_deployment_postprocessing: bool = False


@dataclass
class BuildConfiguration:
    id: str
    name: str
    settings: dict = field(default_factory=dict)
    base_configuration: str = ""


@dataclass
class Target:
    id: str
    name: str
    isa: str
    product_name: str
    product_type: str
    product_reference: str
    build_phases: list[BuildPhase]
    dependencies: list["Target"]
    configurations: dict[str, BuildConfiguration]
    project: "Project"
    dependency_ids: list[str] = field(default_factory=list)

    @property
    def is_application(self) -> bool:
        return self.product_type.endswith(".application")

    @property
    def is_static_library(self) -> bool:
        return self.product_type.endswith((".library.static", ".static-library"))

    def settings(self, configuration: str) -> xcconfig.SettingsStack:
        """Settings for this target, highest priority first."""
        layers: list[xcconfig.Settings] = []
        config = self.configurations.get(configuration) or next(iter(self.configurations.values()), None)
        if config:
            layers.append(xcconfig.Settings(config.settings))
            if config.base_configuration and os.path.isfile(config.base_configuration):
                layers.append(xcconfig.load(config.base_configuration))
        layers.extend(self.project.layers(configuration))
        return xcconfig.SettingsStack(
            layers,
            defaults={
                "TARGET_NAME": self.name,
                "PRODUCT_NAME": self.product_name or self.name,
                "CONFIGURATION": configuration,
                "SRCROOT": self.project.source_root,
                "PROJECT_DIR": self.project.source_root,
                "PODS_ROOT": os.path.join(self.project.source_root, "Pods"),
                "PLATFORM_NAME": "iphoneos",
            },
        )

    def source_files(self) -> list[FileRef]:
        result: list[FileRef] = []
        for phase in self.build_phases:
            if phase.isa != "PBXSourcesBuildPhase":
                continue
            for build_file in phase.files:
                if build_file.file_ref:
                    result.append(build_file.file_ref)
        return result

    def resource_files(self) -> list[FileRef]:
        result: list[FileRef] = []
        for phase in self.build_phases:
            if phase.isa != "PBXResourcesBuildPhase":
                continue
            for build_file in phase.files:
                if build_file.file_ref:
                    result.append(build_file.file_ref)
        return result

    def script_phases(self) -> list[BuildPhase]:
        return [phase for phase in self.build_phases if phase.isa in ("PBXShellScriptBuildPhase", "PBXRunScriptBuildPhase")]


class Project:
    """One `project.pbxproj`."""

    def __init__(self, path: str) -> None:
        self.path = os.path.abspath(path)
        if os.path.isdir(self.path):
            self.path = os.path.join(self.path, "project.pbxproj")
        with open(self.path, "r", encoding="utf-8", errors="replace") as handle:
            text = handle.read()
        root = openstep.loads(text)
        self.objects: dict[str, dict] = root.get("objects", {})
        self.root_id: str = root.get("rootObject", "")
        self.root = self.objects.get(self.root_id, {})
        self.source_root = os.path.dirname(os.path.dirname(self.path))
        self.project_dir = os.path.join(self.source_root, self.root.get("projectDirPath", ""))
        self._paths: dict[str, str] = {}
        self._resolve_group_paths()
        self.targets = [self._load_target(target_id) for target_id in self.root.get("targets", [])]
        by_id = {target.id: target for target in self.targets}
        for target in self.targets:
            target.dependencies = [by_id[dep] for dep in target.dependency_ids if dep in by_id]

    # -- group tree -----------------------------------------------------

    def group_tree(self) -> dict:
        """The `mainGroup` subtree in PIF form."""
        main_group = self.root.get("mainGroup")
        if not main_group:
            return {"type": "group", "guid": f"{self.root_id}::MAINGROUP", "name": "",
                    "path": "", "sourceTree": "<group>", "children": []}
        return self._group_node(main_group, root=True)

    def _group_node(self, group_id: str, root: bool = False) -> dict:
        group = self.objects.get(group_id, {})
        node = {
            "type": "group",
            "guid": group_id,
            "name": group.get("name", "") if not root else "",
            "path": group.get("path", ""),
            "sourceTree": group.get("sourceTree", "<group>"),
            "children": [],
        }
        for child_id in group.get("children", []):
            child = self.objects.get(child_id, {})
            isa = child.get("isa", "")
            if isa in ("PBXGroup", "PBXVariantGroup", "XCVersionGroup"):
                node["children"].append(self._group_node(child_id))
            elif isa == "PBXFileReference":
                node["children"].append(
                    {
                        "type": "file",
                        "guid": child_id,
                        "name": child.get("name", child.get("path", "")),
                        "path": child.get("path", ""),
                        "sourceTree": child.get("sourceTree", "<group>"),
                        "fileType": child.get("lastKnownFileType", "")
                        or file_type_for(child.get("path", "")),
                    }
                )
        return node

    def _resolve_group_paths(self) -> None:
        main_group = self.root.get("mainGroup")
        if main_group:
            self._walk_group(main_group, self.project_dir)

    def _walk_group(self, group_id: str, parent: str) -> None:
        group = self.objects.get(group_id)
        if not group:
            return
        path = self._join(parent, group.get("path", ""), group.get("sourceTree", "<group>"))
        for child_id in group.get("children", []):
            child = self.objects.get(child_id)
            if not child:
                continue
            if child.get("isa") in ("PBXGroup", "PBXVariantGroup", "XCVersionGroup"):
                self._walk_group(child_id, path)
            else:
                self._paths[child_id] = self._join(
                    path, child.get("path", child.get("name", "")), child.get("sourceTree", "<group>")
                )

    def _join(self, parent: str, path: str, source_tree: str) -> str:
        if not path:
            return parent
        if source_tree == "<absolute>":
            return path
        if source_tree in ("SOURCE_ROOT", "SRCROOT"):
            return os.path.normpath(os.path.join(self.source_root, path))
        if source_tree == "BUILT_PRODUCTS_DIR":
            return os.path.join("$(BUILT_PRODUCTS_DIR)", path)
        return os.path.normpath(os.path.join(parent, path))

    def file_ref(self, file_id: str) -> FileRef | None:
        obj = self.objects.get(file_id)
        if not obj:
            return None
        if obj.get("isa") not in ("PBXFileReference", "PBXVariantGroup", "XCVersionGroup"):
            # A build file may point at a group; take its first child.
            children = obj.get("children") or []
            return self.file_ref(children[0]) if children else None
        path = obj.get("path", "")
        source_tree = obj.get("sourceTree", "<group>")
        absolute = self._paths.get(file_id, "")
        return FileRef(
            id=file_id,
            path=path,
            name=obj.get("name", ""),
            source_tree=source_tree,
            file_type=obj.get("lastKnownFileType", ""),
            absolute_path=absolute,
        )

    # -- targets --------------------------------------------------------

    def _load_target(self, target_id: str) -> Target:
        obj = self.objects.get(target_id, {})
        phases = [self._load_phase(phase_id) for phase_id in obj.get("buildPhases", [])]
        configurations = self._load_configurations(obj.get("buildConfigurationList", ""))
        dependencies: list[Target] = []
        dependency_ids: list[str] = []
        for dependency_id in obj.get("dependencies", []):
            dependency = self.objects.get(dependency_id, {})
            proxied = dependency.get("targetProxy", "")
            remote = self.objects.get(self.objects.get(proxied, {}).get("remoteGlobalIDString", ""), {})
            remote_id = self.objects.get(proxied, {}).get("remoteGlobalIDString", "")
            if remote_id:
                dependency_ids.append(remote_id)
            name = remote.get("name") or remote.get("productName") or ""
            for candidate in getattr(self, "targets", []):
                if candidate.name == name:
                    dependencies.append(candidate)
        return Target(
            id=target_id,
            name=obj.get("name", ""),
            isa=obj.get("isa", ""),
            product_name=obj.get("productName", obj.get("name", "")),
            product_type=obj.get("productType", ""),
            product_reference=obj.get("productReference", ""),
            build_phases=[phase for phase in phases if phase],
            dependencies=dependencies,
            configurations=configurations,
            project=self,
            dependency_ids=dependency_ids,
        )

    def _load_phase(self, phase_id: str) -> BuildPhase | None:
        obj = self.objects.get(phase_id)
        if not obj:
            return None
        phase = BuildPhase(
            id=phase_id,
            isa=obj.get("isa", ""),
            name=obj.get("name", ""),
            shell_script=obj.get("shellScript", ""),
            input_paths=list(obj.get("inputPaths", [])),
            output_paths=list(obj.get("outputPaths", [])),
            input_file_lists=list(obj.get("inputFileListPaths", [])),
            output_file_lists=list(obj.get("outputFileListPaths", [])),
            always_out_of_date=pbx_bool(obj.get("alwaysOutOfDate")),
            run_only_for_deployment_postprocessing=pbx_bool(
                obj.get("runOnlyForDeploymentPostprocessing")
            ),
        )
        for build_file_id in obj.get("files", []):
            build_file = self.objects.get(build_file_id, {})
            file_id = build_file.get("fileRef", "")
            phase.files.append(
                BuildFile(
                    id=build_file_id,
                    file_ref=self.file_ref(file_id),
                    settings=build_file.get("settings", {}),
                )
            )
        return phase

    def _load_configurations(self, list_id: str) -> dict[str, BuildConfiguration]:
        configuration_list = self.objects.get(list_id, {})
        result: dict[str, BuildConfiguration] = {}
        for config_id in configuration_list.get("buildConfigurations", []):
            obj = self.objects.get(config_id, {})
            base = obj.get("baseConfigurationReference", "")
            base_path = self._paths.get(base, "")
            result[obj.get("name", "")] = BuildConfiguration(
                id=config_id,
                name=obj.get("name", ""),
                settings=dict(obj.get("buildSettings", {})),
                base_configuration=base_path,
            )
        return result

    # -- project-level settings ----------------------------------------

    def layers(self, configuration: str) -> list[xcconfig.Settings]:
        configuration_list = self.objects.get(self.root.get("buildConfigurationList", ""), {})
        for config_id in configuration_list.get("buildConfigurations", []):
            obj = self.objects.get(config_id, {})
            if obj.get("name") != configuration:
                continue
            layers = [xcconfig.Settings(obj.get("buildSettings", {}))]
            base_path = self._paths.get(obj.get("baseConfigurationReference", ""), "")
            if base_path and os.path.isfile(base_path):
                layers.append(xcconfig.load(base_path))
            return layers
        return []

    def target(self, name: str) -> Target | None:
        for target in self.targets:
            if target.name == name:
                return target
        return None

    def application_targets(self) -> Iterator[Target]:
        for target in self.targets:
            if target.is_application:
                yield target


class Workspace:
    """`.xcworkspace` → the projects it references."""

    _FILE_REF = "<FileRef location ="

    def __init__(self, path: str) -> None:
        self.path = os.path.abspath(path)
        if os.path.isdir(self.path):
            self.path = os.path.join(self.path, "contents.xcworkspacedata")
        self.directory = os.path.dirname(os.path.dirname(self.path))
        self.project_paths: list[str] = []
        if os.path.isfile(self.path):
            import re

            with open(self.path, "r", encoding="utf-8", errors="replace") as handle:
                text = handle.read()
            for match in re.finditer(r'location\s*=\s*"([^"]+)"', text):
                location = match.group(1)
                if not location.startswith("group:"):
                    continue
                relative = location[len("group:"):]
                candidate = os.path.normpath(os.path.join(self.directory, relative))
                if os.path.exists(candidate):
                    self.project_paths.append(candidate)

    @classmethod
    def discover(cls, ios_dir: str) -> "Workspace | None":
        for name in sorted(os.listdir(ios_dir)):
            if name.endswith(".xcworkspace"):
                return cls(os.path.join(ios_dir, name))
        return None
