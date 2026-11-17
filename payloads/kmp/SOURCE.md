# Vendored source

The Kotlin Multiplatform scripts in this directory are vendored from the
pinned apple-dev commit in `internal/platform/config/sources.lock`.

- apple-dev commit: `d0f4980cb39c43ee75506811435f1bacf0b10816`

The script uses `$IDEV_HOME/cache/kmp`, requires `DARWIN_SDK_BUNDLE`, and keeps
Gradle's Kotlin/Native data under the same private cache.

Follow `docs/UPSTREAM.md` to port upstream changes. The user-facing surface
does not show the upstream names.
