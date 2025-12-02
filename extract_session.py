#!/usr/bin/env python3
"""
Extract a specified cookie from a Mozilla Firefox profile.
Compatible with Python 3.13, strict type checking, and argparse CLI.
"""

import configparser
import sqlite3
import pathlib
import platform
import argparse
import sys
import typing as T


def get_default_firefox_profile() -> pathlib.Path | None:
    """
    Locate the default Firefox profile directory.

    :return: Path object pointing to the default profile or None if not found.
    """
    system: str = platform.system()
    if system == "Windows":
        ini_path = (
            pathlib.Path.home()
            / "AppData"
            / "Roaming"
            / "Mozilla"
            / "Firefox"
            / "profiles.ini"
        )
    elif system == "Darwin":  # macOS
        ini_path = (
            pathlib.Path.home()
            / "Library"
            / "Application Support"
            / "Firefox"
            / "profiles.ini"
        )
    elif system == "Linux":
        ini_path = (
            pathlib.Path.home() / ".mozilla" / "firefox" / "profiles.ini"
        )
    else:
        print(f"Unsupported OS: {system}", file=sys.stderr)
        return None

    if not ini_path.exists():
        print(f"profiles.ini not found at {ini_path}", file=sys.stderr)
        return None

    # Parse profiles.ini
    config = configparser.ConfigParser()
    for __ in config.read(ini_path):
        for section_name in config.sections():  # pyright: ignore[reportAny]
            section = config[section_name]
            if "default" in section and section["default"] == "1":
                path_str = section["path"]
                is_relative = (
                    config.get(section_name, "isrelative", fallback="1") == "1"  # pyright: ignore[reportAny]
                )
                profile_path = (
                    ini_path.parent / path_str
                    if is_relative
                    else pathlib.Path(path_str)
                )
                return profile_path.resolve()

    print("Default profile not found.", file=sys.stderr)
    return None


def get_firefox_cookie(
    profile_path: pathlib.Path, cookie_name: str, host: str | None = None
) -> str | None:
    """
    Extract a specified cookie from a Firefox profile.

    :param profile_path: Path to the Firefox profile directory.
    :param cookie_name: Name of the cookie to extract.
    :param host: Optional domain filter (e.g., 'example.com').
    :return: Dictionary with cookie details or None if not found.
    """
    db_path = profile_path / "cookies.sqlite"

    if not db_path.exists():
        print(
            f"Error: cookies.sqlite not found in {profile_path}",
            file=sys.stderr,
        )
        return None

    try:
        conn: sqlite3.Connection = sqlite3.connect(
            f"file:{db_path}?mode=ro&immutable=1", uri=True
        )
        cursor: sqlite3.Cursor = conn.cursor()

        query: str = "SELECT value FROM moz_cookies WHERE name = ?"
        params: list[str] = [cookie_name]

        if host:
            query += " AND host LIKE ?"
            params.append(f"%{host}%")

        __ = cursor.execute(query, params)
        result = T.cast(tuple[str] | None, cursor.fetchone())

        conn.close()

        if result:
            return result[0]
        else:
            print("Cookie not found.", file=sys.stderr)
            return None

    except sqlite3.Error as e:
        print(f"SQLite error: {e}", file=sys.stderr)
        return None


def parse_args() -> argparse.Namespace:
    """
    Parse command-line arguments using argparse.
    """
    parser = argparse.ArgumentParser(
        description="Extract a specified cookie from a Firefox profile."
    )
    default_path = get_default_firefox_profile()
    if (default_path) is not None:
        __ = parser.add_argument(
            "profile_path",
            nargs="?",
            type=pathlib.Path,
            help="Path to the Firefox profile directory (e.g., /path/to/profile).",
            default=default_path,
        )
    else:
        __ = parser.add_argument(
            "profile_path",
            type=pathlib.Path,
            help="Path to the Firefox profile directory (e.g., /path/to/profile).",
        )
    return parser.parse_args()


if __name__ == "__main__":
    args = parse_args()
    cookie: str | None = get_firefox_cookie(
        T.cast(pathlib.Path, args.profile_path),
        "session",
        "adventofcode.com",
    )
    if cookie:
        print(cookie)
