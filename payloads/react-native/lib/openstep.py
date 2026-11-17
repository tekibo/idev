"""Parser for the OpenStep property list used by project.pbxproj.

Xcode project files are not JSON and not XML: they are old-style plists with
unquoted tokens, `/* comments */` and `// comments`. This module parses them
into plain Python data and writes them back out, so the rest of `expo/` can
treat a project as data instead of text.
"""

from __future__ import annotations

import re
from typing import Any

_PUNCT = "{}()=;,"
_BARE_STOP = set(' \t\r\n{}()=;,"')


class ParseError(ValueError):
    pass


class _Parser:
    def __init__(self, text: str) -> None:
        self.text = text
        self.pos = 0
        self.length = len(text)

    def error(self, message: str) -> ParseError:
        line = self.text.count("\n", 0, self.pos) + 1
        return ParseError(f"line {line}: {message}")

    def skip_space(self) -> None:
        while self.pos < self.length:
            char = self.text[self.pos]
            if char in " \t\r\n":
                self.pos += 1
            elif self.text.startswith("/*", self.pos):
                end = self.text.find("*/", self.pos + 2)
                self.pos = self.length if end < 0 else end + 2
            elif self.text.startswith("//", self.pos):
                end = self.text.find("\n", self.pos + 2)
                self.pos = self.length if end < 0 else end + 1
            else:
                return

    def parse(self) -> Any:
        self.skip_space()
        value = self.value()
        self.skip_space()
        if self.pos != self.length:
            raise self.error("trailing content")
        return value

    def value(self) -> Any:
        self.skip_space()
        if self.pos >= self.length:
            raise self.error("unexpected end of input")
        char = self.text[self.pos]
        if char == "{":
            return self.dictionary()
        if char == "(":
            return self.array()
        if char == '"':
            return self.quoted()
        if char in _PUNCT:
            raise self.error(f"unexpected '{char}'")
        return self.bare()

    def dictionary(self) -> dict[str, Any]:
        self.pos += 1  # {
        result: dict[str, Any] = {}
        while True:
            self.skip_space()
            if self.pos >= self.length:
                raise self.error("unterminated dictionary")
            if self.text[self.pos] == "}":
                self.pos += 1
                return result
            key = self.value()
            if not isinstance(key, str):
                raise self.error("dictionary key must be a string")
            self.skip_space()
            if self.pos >= self.length or self.text[self.pos] != "=":
                raise self.error("expected '=' after key")
            self.pos += 1
            result[key] = self.value()
            self.skip_space()
            if self.pos < self.length and self.text[self.pos] == ";":
                self.pos += 1

    def array(self) -> list[Any]:
        self.pos += 1  # (
        result: list[Any] = []
        while True:
            self.skip_space()
            if self.pos >= self.length:
                raise self.error("unterminated array")
            if self.text[self.pos] == ")":
                self.pos += 1
                return result
            result.append(self.value())
            self.skip_space()
            if self.pos < self.length and self.text[self.pos] == ",":
                self.pos += 1

    def quoted(self) -> str:
        self.pos += 1  # opening quote
        chunks: list[str] = []
        while self.pos < self.length:
            char = self.text[self.pos]
            if char == "\\":
                self.pos += 1
                if self.pos >= self.length:
                    raise self.error("unterminated escape")
                escaped = self.text[self.pos]
                chunks.append({"n": "\n", "t": "\t", "r": "\r"}.get(escaped, escaped))
                self.pos += 1
            elif char == '"':
                self.pos += 1
                return "".join(chunks)
            else:
                chunks.append(char)
                self.pos += 1
        raise self.error("unterminated string")

    def bare(self) -> str:
        start = self.pos
        while self.pos < self.length and self.text[self.pos] not in _BARE_STOP:
            if self.text.startswith("/*", self.pos) or self.text.startswith("//", self.pos):
                break
            self.pos += 1
        if self.pos == start:
            raise self.error("expected a value")
        return self.text[start:self.pos]


def loads(text: str) -> Any:
    return _Parser(text).parse()


_NEEDS_QUOTES = re.compile(r'[^A-Za-z0-9_.$/:-]')


def dumps(value: Any, indent: int = 0) -> str:
    """Serialize back to OpenStep form, matching Xcode's layout closely enough
    that a round-trip diff stays small."""
    pad = "\t" * indent
    if isinstance(value, dict):
        if not value:
            return "{}"
        lines = ["{"]
        for key in sorted(value):
            body = dumps(value[key], indent + 1)
            lines.append(f"{pad}\t{_quote(key)} = {body};")
        lines.append(pad + "}")
        return "\n".join(lines)
    if isinstance(value, list):
        if not value:
            return "()"
        lines = ["("]
        for item in value:
            lines.append(f"{pad}\t{dumps(item, indent + 1)},")
        lines.append(pad + ")")
        return "\n".join(lines)
    if isinstance(value, str):
        return _quote(value)
    if isinstance(value, bool):
        return "1" if value else "0"
    return str(value)


def _quote(token: str) -> str:
    if token and not _NEEDS_QUOTES.search(token):
        return token
    escaped = token.replace("\\", "\\\\").replace('"', '\\"').replace("\n", "\\n")
    return f'"{escaped}"'
