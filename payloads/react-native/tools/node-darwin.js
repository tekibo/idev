// Make Node-based Apple build tooling believe it runs on macOS.
//
// Expo's iOS run flow gates on `process.platform === "darwin"`, and its
// autolinking scan has the same gate. On Linux only the gate needs Darwin;
// native helper binaries still have to be the Linux builds. The platform is
// therefore a writable value so the scanner can briefly see the real host.
//
// Loading this module through NODE_OPTIONS enables those code paths without
// patching any package.

const fs = require("fs");
const path = require("path");
const Module = require("module");

// `process.platform` is forced to Darwin for the CLI gates. Probe the host
// directly so native helpers and Metro's file watcher use the real platform:
// /proc is Linux-only, and the process report's header does not expose the OS
// on every Node build.
function hostPlatform() {
  try {
    if (fs.existsSync("/proc/self/exe")) {
      return "linux";
    }
  } catch {
    // Fall through to the process platform.
  }
  return process.platform;
}

let reportedPlatform = "darwin";
Object.defineProperty(process, "platform", {
  get: () => reportedPlatform,
  set: (value) => {
    reportedPlatform = value;
  },
  enumerable: true,
  configurable: true,
});

// Metro's NativeWatcher is selected when os.platform() is Darwin. It uses the
// macOS recursive fs.watch, which does not report nested changes on Linux.
// Keep the CLI gate on process.platform and report the host to os.platform.
const os = require("os");
const host = hostPlatform();
if (host !== "darwin") {
  os.platform = () => host;
}

// React Native's install path reads the built bundle identifier through
// `/usr/libexec/PlistBuddy`, an absolute macOS path. Answer that one call
// from the built plist so the CLI can reach the manager-backed install.
const childProcess = require("child_process");
const originalExecFileSync = childProcess.execFileSync;

function plistValue(pathname, key) {
  const script = [
    "import plistlib, sys",
    "with open(sys.argv[1], 'rb') as handle:",
    "    value = plistlib.load(handle)",
    "for part in sys.argv[2].split('.'):",
    "    value = value[part]",
    "print(value)",
  ].join("\n");
  const result = childProcess.spawnSync(
    "python3",
    ["-c", script, pathname, key],
    { encoding: "utf8" },
  );
  if (result.error) {
    throw result.error;
  }
  if (result.status !== 0) {
    throw new Error(result.stderr.trim() || "PlistBuddy failed");
  }
  return result.stdout.trim();
}

childProcess.execFileSync = function (file, args, options) {
  if (String(file) === "/usr/libexec/PlistBuddy" && Array.isArray(args)) {
    const command = args.find(
      (arg) => typeof arg === "string" && arg.startsWith("Print:"),
    );
    const plist = args.length > 1 ? args[args.length - 1] : "";
    if (command && plist) {
      const value = plistValue(plist, command.slice("Print:".length));
      if (options && options.encoding) {
        return value;
      }
      return Buffer.from(value);
    }
  }
  return originalExecFileSync.apply(this, arguments);
};

const originalLoad = Module._load;

function wrapHostScanner(scanner) {
  if (!scanner || typeof scanner.scanModules !== "function") {
    return;
  }
  if (scanner.scanModules.__idevHostScanner) {
    return;
  }
  if (hostPlatform() !== "linux") {
    return;
  }
  const scanModules = scanner.scanModules;
  const wrapped = function (...args) {
    const previous = process.platform;
    // expo-modules-macros chooses its scanner binary from process.platform.
    // The Expo gate above needs "darwin"; this one call needs the host.
    process.platform = "linux";
    let outcome;
    try {
      outcome = scanModules.apply(this, args);
    } catch (error) {
      process.platform = previous;
      throw error;
    }
    if (outcome && typeof outcome.then === "function") {
      return outcome.finally(() => {
        process.platform = previous;
      });
    }
    process.platform = previous;
    return outcome;
  };
  wrapped.__idevHostScanner = true;
  scanner.scanModules = wrapped;
}

Module._load = function (request, parent, isMain) {
  // lightningcss selects its native binding from process.platform. The
  // Darwin gate would make it require lightningcss.darwin-x64.node, which
  // does not exist on Linux. Load it under the real host platform.
  if (request === "lightningcss") {
    const previous = process.platform;
    process.platform = "linux";
    try {
      return originalLoad.apply(this, arguments);
    } finally {
      process.platform = previous;
    }
  }
  // fb-dotslash selects its binary from os.platform(), which this shim
  // reports as process.platform. The Expo gate needs "darwin"; the DotSlash
  // package must still load its Linux binary on the Linux host.
  if (request === "fb-dotslash") {
    const previous = process.platform;
    process.platform = "linux";
    try {
      return originalLoad.apply(this, arguments);
    } finally {
      process.platform = previous;
    }
  }
  const result = originalLoad.apply(this, arguments);
  try {
    // React Native 20 computes its simulator fallback before it checks a
    // physical-device flag. On Linux there are no simulators, so the eager
    // lookup throws even when --device is present. The fallback is unused in
    // that branch; return a placeholder only for the explicit device case.
    if (typeof request === "string" && request.includes("getDestinationSimulator")) {
      const original = result.getDestinationSimulator;
      if (typeof original === "function" && !original.__idevDeviceFallback) {
        const wrapped = function (args, fallbackSimulators) {
          try {
            return original(args, fallbackSimulators);
          } catch (error) {
            if (args && (args.device || args.udid)) {
              return {
                udid: "idev-device-fallback",
                name: "placeholder",
                state: "Shutdown",
                version: "iOS 17.0",
              };
            }
            throw error;
          }
        };
        wrapped.__idevDeviceFallback = true;
        return { ...result, getDestinationSimulator: wrapped };
      }
    }

    // Expo's iOS run flow asks the local keychain for a signing identity
    // before it invokes xcodebuild. This pipeline never signs in xcodebuild:
    // the build stays unsigned and xtool signs at install.
    if (typeof request === "string" && request.includes("configureCodeSigning")) {
      return {
        ...result,
        ensureDeviceIsCodeSignedForDeploymentAsync: async () => null,
      };
    }

    // The primary device install path is AppleDevice/installation_proxy,
    // which sends the unsigned build and fails verification. The manager owns
    // the install/launch surface: when IDEV_BIN is present, hand the resolved
    // .app to `idev install-app`, which calls xtool (the only signer) and then
    // launches the installed app. The devicectl fallback keeps the old
    // container flow working while the manager is introduced.
    if (typeof request === "string" && request.includes("installOnDeviceAsync")) {
      return {
        ...result,
        installOnDeviceAsync: async (props) => {
          if (process.env.IDEV_BUILD_ONLY === "1") {
            return;
          }
          const idev = process.env.IDEV_BIN;
          if (idev) {
            const { spawnSync } = require("child_process");
            const args = [
              "install-app",
              "--bundle", props.bundle,
              "--bundle-id", props.bundleIdentifier,
            ];
            if (props.udid) {
              args.push("--udid", props.udid);
            }
            if (process.env.IDEV_DEV_SERVER_URL) {
              args.push("--dev-server-url", process.env.IDEV_DEV_SERVER_URL);
            }
            const outcome = spawnSync(idev, args, {
              stdio: "inherit",
              env: process.env,
            });
            if (outcome.error) {
              throw outcome.error;
            }
            if (outcome.status !== 0) {
              throw new Error(`idev install-app failed with status ${outcome.status}`);
            }
            return;
          }

          const resolved = Module._resolveFilename(request, parent, isMain);
          const target = path.resolve(
            path.dirname(resolved),
            "../../../start/platforms/ios/devicectl.js",
          );
          const devicectl = originalLoad(target, parent, false);
          return devicectl.installAndLaunchAppAsync(props);
        },
      };
    }

    // Expo's autolinking scan is gated on Darwin but invokes a native scanner.
    // Let that single call select the Linux scanner binary.
    if (typeof request === "string" && request.includes("expo-modules-macros")) {
      wrapHostScanner(result);
    }
  } catch {
    // Fall through to the unmodified module if the shape is unexpected.
  }
  return result;
};
