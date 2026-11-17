// Entry point for the Expo build engine.
//
// `swbuild`'s console hardcodes `developerPath: nil` when it creates a session,
// so PlatformRegistry never sees the Apple platforms and every iOS build fails
// with "unable to resolve run destination platform". This calls the same
// library with an explicit developer directory.

import Foundation
import SwiftBuild
import SWBUtil

struct Options {
    var pif = ""
    var developerDir = ""
    var target = ""
    var buildAll = false
    var resourcePaths: [String] = []
}

func parseArguments() -> Options {
    var options = Options()
    var iterator = CommandLine.arguments.dropFirst().makeIterator()
    while let argument = iterator.next() {
        switch argument {
        case "--pif": options.pif = iterator.next() ?? ""
        case "--developer-dir": options.developerDir = iterator.next() ?? ""
        case "--target": options.target = iterator.next() ?? ""
        case "--resource-path": options.resourcePaths.append(iterator.next() ?? "")
        case "--all-targets": options.buildAll = true
        default:
            FileHandle.standardError.write(Data("expo-build: unknown argument \(argument)\n".utf8))
        }
    }
    return options
}

struct UnsupportedExternalTool: Error {}

let defaultToolchainIdentifier = "com.apple.dt.toolchain.XcodeDefault"

/// The Darwin SDK artifact bundle ships a toolchain without `Info.plist`, so
/// ToolchainRegistry cannot identify it as the default one. Build a small
/// overlay that points at the real Platforms/Toolchains and adds that file,
/// instead of writing into the installed SDK.
func prepareDeveloperDirectory(_ original: String) throws -> String {
    let fm = FileManager.default
    let toolchain = "\(original)/Toolchains/XcodeDefault.xctoolchain"
    if fm.fileExists(atPath: "\(toolchain)/Info.plist") {
        return original
    }

    let overlay = "\(NSTemporaryDirectory())expo-developer-\(getpid())"
    try? fm.removeItem(atPath: overlay)
    try fm.createDirectory(atPath: "\(overlay)/Toolchains", withIntermediateDirectories: true)
    try fm.createDirectory(
        atPath: "\(overlay)/Toolchains/XcodeDefault.xctoolchain/usr",
        withIntermediateDirectories: true
    )
    try fm.createSymbolicLink(atPath: "\(overlay)/Platforms", withDestinationPath: "\(original)/Platforms")

    // Mirror the toolchain's usr, but with a real bin directory so Darwin
    // tool names can be supplied without touching the installed SDK.
    let overlayUsr = "\(overlay)/Toolchains/XcodeDefault.xctoolchain/usr"
    let realUsr = "\(toolchain)/usr"
    try fm.createDirectory(atPath: "\(overlayUsr)/bin", withIntermediateDirectories: true)
    for entry in (try? fm.contentsOfDirectory(atPath: realUsr)) ?? [] {
        let target = "\(realUsr)/\(entry)"
        let link = "\(overlayUsr)/\(entry)"
        if entry == "bin" {
            for tool in (try? fm.contentsOfDirectory(atPath: target)) ?? [] {
                try? fm.createSymbolicLink(atPath: "\(link)/\(tool)", withDestinationPath: "\(target)/\(tool)")
            }
        } else {
            try fm.createSymbolicLink(atPath: link, withDestinationPath: target)
        }
    }

    // Apple's libtool creates static archives; Linux ships GNU libtool.
    let libtool = """
        #!/bin/sh
        if [ "$1" = "--version" ] || [ "$1" = "-V" ] || [ "$1" = "-v" ]; then
            echo "Apple Inc. version cctools-1000.4.5"
            exit 0
        fi
        exec /usr/bin/llvm-libtool-darwin "$@"
        """
    let libtoolPath = "\(overlayUsr)/bin/libtool"
    try libtool.write(toFile: libtoolPath, atomically: true, encoding: .utf8)
    try fm.setAttributes([.posixPermissions: 0o755], ofItemAtPath: libtoolPath)

    let info: [String: Any] = [
        "Identifier": defaultToolchainIdentifier,
        "Aliases": ["default"],
        "DisplayName": "Xcode Default",
        "Version": "27.0",
    ]
    let data = try PropertyListSerialization.data(fromPropertyList: info, format: .xml, options: 0)
    try data.write(to: URL(fileURLWithPath: "\(overlay)/Toolchains/XcodeDefault.xctoolchain/Info.plist"))

    return overlay
}

final class PlanningDelegate: SWBPlanningOperationDelegate {
    func provisioningTaskInputs(targetGUID: String, provisioningSourceData: SWBProvisioningTaskInputsSourceData) async -> SWBProvisioningTaskInputs {
        // The pipeline stays unsigned; xtool is the only signer.
        SWBProvisioningTaskInputs()
    }

    func executeExternalTool(commandLine: [String], workingDirectory: String?, environment: [String: String]) async throws -> SWBExternalToolResult {
        // `.deferred` lets the build system run the tool itself; returning a
        // result here would mean the client executed it.
        .deferred
    }
}

@main
struct ExpoBuild {
    static func main() async {
        let options = parseArguments()
        guard !options.pif.isEmpty, !options.developerDir.isEmpty else {
            FileHandle.standardError.write(Data("expo-build: --pif and --developer-dir are required\n".utf8))
            exit(2)
        }

        var failed = false
        let stdout = FileHandle.standardOutput

        func emit(_ message: SwiftBuildMessage) throws {
            let stream = OutputByteStream()
            try stream.write(JSONEncoder(outputFormatting: [.sortedKeys, .withoutEscapingSlashes]).encode(message))
            stream.write("\n")
            let payload = stream.bytes
            if let text = String(data: Data(payload.bytes), encoding: .utf8), text.contains("\"kind\":\"error\"") {
                failed = true
            }
            try stdout.write(contentsOf: Data("\(payload.count - 1)\n".utf8))
            try stdout.write(contentsOf: Data(payload.bytes))
        }

        do {
            let service = try await SWBBuildService(connectionMode: .default)
            defer { Task { await service.close() } }

            var resourcePaths = options.resourcePaths
            let developerDir = try prepareDeveloperDirectory(options.developerDir)
            if resourcePaths.isEmpty {
                let developer = URL(fileURLWithPath: developerDir)
                resourcePaths = [
                    developer.appendingPathComponent("Platforms").path,
                    developer.appendingPathComponent("Toolchains/XcodeDefault.xctoolchain").path,
                ]
            }

            let (result, diagnostics) = await service.createSession(
                name: "expo",
                developerPath: developerDir,
                resourceSearchPaths: resourcePaths,
                cachePath: nil,
                inferiorProductsPath: nil,
                // Script phases (Expo's configure script, CocoaPods' embed
                // helpers) spawn node and other tools; Swift Build otherwise
                // gives them a minimal environment in which those are missing.
                environment: [
                    "PATH": ProcessInfo.processInfo.environment["PATH"]
                        ?? "/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin"
                ]
            )
            for diagnostic in diagnostics {
                FileHandle.standardError.write(Data("\(diagnostic.kind.rawValue): \(diagnostic.message)\n".utf8))
            }
            let session: SWBBuildServiceSession
            switch result {
            case .success(let created): session = created
            case .failure(let error):
                FileHandle.standardError.write(Data("expo-build: could not create session: \(error)\n".utf8))
                exit(1)
            }

            try await session.loadWorkspace(containerPath: options.pif)

            let parametersPath = URL(fileURLWithPath: options.pif)
                .deletingLastPathComponent()
                .appendingPathComponent("build-parameters.json")
            let parameters = try JSONDecoder().decode(
                SWBBuildParameters.self,
                from: Path(parametersPath.path),
                fs: localFS
            )

            let info = try await session.workspaceInfo()
            var configured: [SWBConfiguredTarget] = []
            if options.buildAll {
                configured = info.targetInfos.map { SWBConfiguredTarget(guid: $0.guid) }
            } else {
                let matches = info.targetInfos.filter { $0.targetName == options.target }
                guard !matches.isEmpty else {
                    let names = info.targetInfos.map(\.targetName).sorted().joined(separator: ", ")
                    FileHandle.standardError.write(Data("expo-build: no target '\(options.target)'. Known: \(names)\n".utf8))
                    exit(2)
                }
                configured = matches.map { SWBConfiguredTarget(guid: $0.guid) }
            }

            var request = SWBBuildRequest()
            request.parameters = parameters
            request.configuredTargets = configured
            request.useParallelTargets = true
            request.useImplicitDependencies = false
            request.useDryRun = false
            request.hideShellScriptEnvironment = false
            request.showNonLoggedProgress = true

            try await session.setSystemInfo(.default())
            try await session.setUserInfo(.default)

            let operation = try await session.createBuildOperation(request: request, delegate: PlanningDelegate())
            let stream = try await operation.start()
            for await message in stream {
                try emit(message)
            }
            try? await session.close()
            exit(failed ? 1 : 0)
        } catch {
            FileHandle.standardError.write(Data("expo-build: \(error)\n".utf8))
            exit(1)
        }
    }
}
