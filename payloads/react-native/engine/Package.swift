// swift-tools-version:6.2
import Foundation
import PackageDescription

// The engine entry point. `swbuild`'s console always creates its session with
// `developerPath: nil`, so it can never register Apple platforms. This package
// calls the same library with the Darwin SDK's Developer directory instead.
//
// Point SWIFT_BUILD_PATH at a swift-build checkout when it is not vendored.
let swiftBuildPath = ProcessInfo.processInfo.environment["SWIFT_BUILD_PATH"] ?? "../vendor/swift-build"

let package = Package(
    name: "expo-build",
    dependencies: [
        .package(name: "SwiftBuild", path: swiftBuildPath)
    ],
    targets: [
        .executableTarget(
            name: "expo-build",
            dependencies: [
                .product(name: "SwiftBuild", package: "SwiftBuild"),
                .product(name: "SWBUtil", package: "SwiftBuild")
            ]
        )
    ],
    swiftLanguageModes: [.v5]
)
