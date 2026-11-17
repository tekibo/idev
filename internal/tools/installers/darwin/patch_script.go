package darwin

const patchScript = `
import json
import os
import plistlib
import re
import sys

bundle = sys.argv[1]
swift_root = sys.argv[2]

for platform_name in ("iPhoneOS.platform", "MacOSX.platform"):
    path = os.path.join(bundle, "Developer/Platforms", platform_name, "Info.plist")
    if not os.path.isfile(path):
        continue
    with open(path, "rb") as handle:
        data = plistlib.load(handle)
    data.setdefault("DefaultProperties", {})["STRINGS_FILE_INPUT_ENCODING"] = "utf-8"
    with open(path + ".tmp", "wb") as handle:
        plistlib.dump(data, handle)
    os.replace(path + ".tmp", path)

toolset = os.path.join(bundle, "toolset-swb.json")
with open(toolset) as handle:
    data = json.load(handle)
data.setdefault("swiftCompiler", {})["extraCLIOptions"] = ["-Xfrontend", "-enable-cross-import-overlays"]
swift_sdk_lib = "/".join(("", "usr", "lib", "swift"))
data.setdefault("linker", {})["extraCLIOptions"] = ["-lswiftCore", "-L" + swift_sdk_lib, "-all_load"]
with open(toolset + ".tmp", "w") as handle:
    json.dump(data, handle, indent=4)
os.replace(toolset + ".tmp", toolset)

for platform_name in ("iPhoneOS.platform", "MacOSX.platform"):
    sdks = os.path.join(bundle, "Developer/Platforms", platform_name, "Developer/SDKs")
    if not os.path.isdir(sdks):
        continue
    for name in sorted(os.listdir(sdks)):
        if os.path.islink(os.path.join(sdks, name)):
            continue
        json_path = os.path.join(sdks, name, "SDKSettings.json")
        if os.path.isfile(json_path):
            with open(json_path) as handle:
                settings = json.load(handle)
            props = settings.setdefault("DefaultProperties", {})
            props["STRINGS_FILE_INPUT_ENCODING"] = "utf-8"
            props["STRINGS_FILE_OUTPUT_ENCODING"] = "binary"
            props.setdefault("GENERATE_INTERMEDIATE_TEXT_BASED_STUBS", "NO")
            props.setdefault("GENERATE_TEXT_BASED_STUBS", "NO")
            with open(json_path, "w") as handle:
                json.dump(settings, handle)
        plist_path = os.path.join(sdks, name, "SDKSettings.plist")
        if os.path.isfile(plist_path):
            with open(plist_path, "rb") as handle:
                settings = plistlib.load(handle)
            props = settings.setdefault("DefaultProperties", {})
            props["GENERATE_INTERMEDIATE_TEXT_BASED_STUBS"] = "NO"
            props["GENERATE_TEXT_BASED_STUBS"] = "NO"
            with open(plist_path + ".tmp", "wb") as handle:
                plistlib.dump(settings, handle)
            os.replace(plist_path + ".tmp", plist_path)

specs = (
    ("usr/share/pm/SwiftBuild_SWBUniversalPlatform.bundle/CopyStringsFile.xcspec", "STRINGS_FILE_INPUT_ENCODING", "utf-8"),
    ("usr/share/pm/SwiftBuild_SWBCore.bundle/CoreBuildSystem.xcspec", "STRINGS_FILE_OUTPUT_ENCODING", "binary"),
    ("usr/share/pm/SwiftBuild_SWBCore.bundle/NativeBuildSystem.xcspec", "STRINGS_FILE_OUTPUT_ENCODING", "binary"),
)
for relative, setting, value in specs:
    path = os.path.join(swift_root, relative)
    if not os.path.isfile(path):
        continue
    with open(path) as handle:
        lines = handle.read().split("\n")
    for index, line in enumerate(lines):
        if f'Name = "{setting}";' in line or f'Name = {setting};' in line:
            for offset in range(index + 1, min(index + 12, len(lines))):
                if "DefaultValue =" in lines[offset]:
                    lines[offset] = re.sub(r'"[^"]*"\s*;', f'"{value}";', lines[offset])
                    break
    with open(path, "w") as handle:
        handle.write("\n".join(lines))
`
