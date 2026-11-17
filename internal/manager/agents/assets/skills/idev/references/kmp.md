# Kotlin Multiplatform

```bash
idev setup --flow kmp
idev kmp /path/to/gradle/root
idev kmp . --module :shared --dynamic --xcframework
idev kmp . --target ios_simulator_arm64
```

`idev kmp` links the iOS framework from the klibs Gradle builds on Linux. Options:

| Option | Meaning |
|---|---|
| `--module :path` | One Gradle module. Default: every project declaring a framework. |
| `--target ios_arm64` | Device slice. Default. |
| `--target ios_simulator_arm64` | Simulator slice. |
| `--out DIR` | Output. Default: `<module>/build/omarchy-kmp/<target>`. |
| `--dynamic` | Dynamic framework. |
| `--xcframework` | Wrap the result as an xcframework. |

The first run downloads the Kotlin/Native iOS runtime from Maven Central (about
300 MB) and, on x86_64, LLVM and libffi for the host (about 580 MB extracted).
Both stay in `$IDEV_HOME/cache/kmp`. The x86_64 archives are checked against
SHA-256 pins in the script; `KMP_X64_SUMS` can point at a SHA256SUMS file you
trust for another Kotlin/Native version. An archive with no digest is not
fetched.

After linking, declare the framework as a binary target in the converted xtool
package and build the Swift app normally. The framework is signed with the app
by `idev sign-dev`/xtool.

## Rules

- A JDK is required on `PATH`. For a build with Android targets, point the
  Android SDK `sdk.dir` at the SDK.
- Do not skip the cache; do not delete another agent's cache mid-link.
- `idev kmp` is a passthrough; use `idev help kmp` for the current flags.
