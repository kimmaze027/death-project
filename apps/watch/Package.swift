// swift-tools-version: 6.1
import PackageDescription

let package = Package(
    name: "WatchCore",
    platforms: [
        .watchOS(.v10),
        .macOS(.v15),
    ],
    products: [
        .library(
            name: "WatchCore",
            targets: ["WatchCore"]
        )
    ],
    targets: [
        .target(
            name: "WatchCore"
        ),
        .testTarget(
            name: "WatchCoreTests",
            dependencies: ["WatchCore"]
        )
    ]
)
