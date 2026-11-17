"""`.xcconfig` parsing and build-setting resolution.

CocoaPods expresses almost everything through xcconfig files, so this is the
authoritative source for flags. The rule here is to expand values the way
Xcode does and never invent defaults.
"""

from __future__ import annotations

import os
import re
from dataclasses import dataclass, field

_VAR = re.compile(r"\$[({]([A-Za-z_][A-Za-z0-9_]*)(?::([^)}]+))?[)}]")


class Settings(dict):
    """A build-settings dictionary that remembers inherited list values."""

    def merged(self, other: "Settings") -> "Settings":
        result = Settings(self)
        for key, value in other.items():
            if value == "$(inherited)" and key in result:
                continue
            result[key] = value
        return result


def load(path: str, _seen: set[str] | None = None) -> Settings:
    """Read an xcconfig file, following `#include` directives."""
    seen = _seen or set()
    absolute = os.path.abspath(path)
    if absolute in seen or not os.path.isfile(absolute):
        return Settings()
    seen.add(absolute)

    settings = Settings()
    with open(absolute, "r", encoding="utf-8", errors="replace") as handle:
        raw_lines = handle.readlines()

    lines: list[str] = []
    pending = ""
    for line in raw_lines:
        stripped = line.rstrip("\n")
        if stripped.endswith("\\"):
            pending += stripped[:-1]
            continue
        lines.append(pending + stripped)
        pending = ""
    if pending:
        lines.append(pending)

    for line in lines:
        text = line.strip()
        if not text:
            continue
        if text.startswith("#include?"):
            target = _unquote(text[len("#include?"):].strip())
            settings.update(load(os.path.join(os.path.dirname(absolute), target), seen))
            continue
        if text.startswith("#include"):
            target = _unquote(text[len("#include"):].strip())
            settings.update(load(os.path.join(os.path.dirname(absolute), target), seen))
            continue
        if text.startswith("//"):
            continue
        if "=" not in text:
            continue
        key, _, value = text.partition("=")
        key = key.strip()
        if not key or key.startswith("//"):
            continue
        settings[key] = value.strip()
    return settings


def _unquote(value: str) -> str:
    if len(value) >= 2 and value[0] == value[-1] and value[0] in "\"'":
        return value[1:-1]
    return value


def expand(value: str, env: dict[str, str], _depth: int = 0) -> str:
    """Substitute `$(VAR)` / `${VAR}` using `env`, recursively."""
    if not isinstance(value, str) or not value:
        return value
    if _depth > 32:
        return value

    def replace(match: re.Match[str]) -> str:
        name, modifier = match.group(1), match.group(2)
        if name == "inherited":
            return ""
        if name not in env:
            # Leave unknown macros for the build system to resolve, the way
            # Xcode does; collapsing them to "" produces paths like /Release/X.
            return match.group(0)
        resolved = env[name]
        if modifier:
            for operation in modifier.split(":"):
                if operation == "lower":
                    resolved = resolved.lower()
                elif operation == "upper":
                    resolved = resolved.upper()
                elif operation == "rfc1034identifier":
                    resolved = re.sub(r"[^A-Za-z0-9.-]", "-", resolved)
                elif operation.startswith("default="):
                    resolved = resolved or operation[len("default="):]
        return resolved

    previous = None
    current = value
    while previous != current and _depth < 32:
        previous = current
        current = _VAR.sub(replace, current)
        _depth += 1
    return current


def split(value: str) -> list[str]:
    """Split a settings list value, honouring quotes."""
    tokens: list[str] = []
    current = ""
    quote = ""
    for char in value:
        if quote:
            if char == quote:
                quote = ""
            else:
                current += char
        elif char in "\"'":
            quote = char
        elif char.isspace():
            if current:
                tokens.append(current)
                current = ""
        else:
            current += char
    if current:
        tokens.append(current)
    return [token for token in tokens if token and token != "$(inherited)"]


@dataclass
class SettingsStack:
    """Resolve settings across project, xcconfig and target layers.

    Xcode's precedence is: target build settings > target xcconfig >
    project build settings > project xcconfig. `layers` is ordered from
    highest to lowest priority.
    """

    layers: list[Settings] = field(default_factory=list)
    defaults: dict[str, str] = field(default_factory=dict)

    def _env(self) -> dict[str, str]:
        env = dict(self.defaults)
        for layer in reversed(self.layers):  # lowest priority first
            for key, value in layer.items():
                if isinstance(value, list):
                    value = " ".join(str(item) for item in value)
                env[key] = value
        return env

    def raw(self, key: str, default: str = "") -> str:
        """Resolve one setting across layers, honouring `$(inherited)`.

        Xcode merges lower-priority values where a higher layer writes
        `$(inherited)`. Without that, a target setting as small as
        `OTHER_LDFLAGS = $(inherited) -lc++` hides the Pods xcconfig that lists
        every pod library.
        """
        accumulated = ""
        found = False
        for layer in reversed(self.layers):  # lowest priority first
            if key not in layer:
                continue
            value = layer[key]
            if isinstance(value, list):
                value = " ".join(str(item) for item in value)
            if value in ("", "$(inherited)"):
                continue
            found = True
            if "$(inherited)" in value:
                value = value.replace("$(inherited)", accumulated).replace("${inherited}", accumulated)
            accumulated = expand(value, self._env()).strip()
        return accumulated if found else default

    def resolve(self, extra: dict[str, str] | None = None) -> dict[str, str]:
        env = dict(extra or {})
        # Seed every key so nested references resolve.
        for layer in reversed(self.layers):
            for key, value in layer.items():
                env.setdefault(key, value)
        for _ in range(8):
            changed = False
            for key, value in list(env.items()):
                expanded = expand(value, env)
                if expanded != value:
                    env[key] = expanded
                    changed = True
            if not changed:
                break
        return env
