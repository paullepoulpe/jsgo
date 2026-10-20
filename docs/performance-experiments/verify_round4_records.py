#!/usr/bin/env python3
"""Verify a Round4 content-addressed snapshot; optionally restore verified data.

Usage: python3 verify_round4_records.py ARCHIVE [--restore NEW_DIRECTORY]
Restoration never executes archived content and refuses an existing destination.
Only manifest paths are restored, after every blob has been verified. Executable
binaries listed as excluded in the manifest are deliberately not reconstructed.
"""

import argparse
import gzip
import hashlib
import json
import os
from pathlib import Path
import re
import shutil
import stat
import tarfile
import tempfile

FORMAT = "jsgo-round4-content-addressed-v1"
CHUNK = 1024 * 1024
MANIFEST_LIMIT = 64 * CHUNK
HEX = re.compile(r"[0-9a-f]{64}\Z")
DIR_FLAGS = os.O_RDONLY | os.O_DIRECTORY | os.O_NOFOLLOW


def relative_path(value, allow_root=False):
    if allow_root and value == ".":
        return value
    if not isinstance(value, str) or not value or "\\" in value or "\0" in value:
        raise ValueError("invalid relative path")
    parts = value.split("/")
    if any(p in ("", ".", "..") for p in parts) or ":" in parts[0]:
        raise ValueError("unsafe relative path: " + repr(value))
    return value


def open_root(path):
    """Open a directory without following symlinks in any path component."""
    path = Path(os.path.abspath(path))
    fd = os.open(path.anchor, DIR_FLAGS)
    try:
        for part in path.parts[1:]:
            next_fd = os.open(part, DIR_FLAGS, dir_fd=fd)
            os.close(fd)
            fd = next_fd
        return fd
    except BaseException:
        os.close(fd)
        raise


def open_relative(root_fd, path, flags=os.O_RDONLY, mode=0o600):
    relative_path(path)
    parts = path.split("/")
    parent = os.dup(root_fd)
    try:
        for part in parts[:-1]:
            child = os.open(part, DIR_FLAGS, dir_fd=parent)
            os.close(parent)
            parent = child
        return os.open(parts[-1], flags | os.O_NOFOLLOW, mode, dir_fd=parent)
    finally:
        os.close(parent)


def signature(s):
    return (s.st_dev, s.st_ino, s.st_mode, s.st_size, s.st_mtime_ns, s.st_ctime_ns)


def canonical_json(value):
    return (json.dumps(value, sort_keys=True, indent=2, ensure_ascii=True) + "\n").encode("utf-8")


def decode_manifest(data):
    def unique_object(pairs):
        result = {}
        for key, value in pairs:
            if key in result:
                raise ValueError("duplicate manifest key")
            result[key] = value
        return result
    manifest = json.loads(data, object_pairs_hook=unique_object)
    if not isinstance(manifest, dict) or manifest.get("format") != FORMAT or manifest.get("hash_algorithm") != "sha256":
        raise ValueError("unsupported manifest format")
    if set(manifest) != {"format", "hash_algorithm", "source_root", "include_prefixes", "directories", "files", "excluded_files"}:
        raise ValueError("unexpected manifest fields")
    if not isinstance(manifest["source_root"], str):
        raise ValueError("invalid source root label")
    prefixes = manifest["include_prefixes"]
    if not isinstance(prefixes, list) or not prefixes or not all(isinstance(p, str) for p in prefixes) or prefixes != sorted(set(prefixes)):
        raise ValueError("invalid selection prefixes")
    for prefix in prefixes:
        relative_path(prefix, allow_root=True)
    def selected(path):
        return any(p == "." or path == p or path.startswith(p + "/") for p in prefixes)
    dirs = manifest["directories"]
    if not isinstance(dirs, list) or not all(isinstance(p, str) for p in dirs) or dirs != sorted(set(dirs)):
        raise ValueError("invalid directory list")
    for path in dirs:
        relative_path(path)
        if not selected(path):
            raise ValueError("directory outside selection")
    seen, blobs = set(dirs), {}
    for field in ("files", "excluded_files"):
        rows = manifest[field]
        if not isinstance(rows, list):
            raise ValueError("invalid file list")
        paths = []
        for row in rows:
            fields = {"path", "size", "sha256", "mode"}
            if field == "excluded_files":
                fields.add("reason")
            if not isinstance(row, dict) or set(row) != fields:
                raise ValueError("invalid file record")
            path = relative_path(row["path"])
            if path in seen or not selected(path):
                raise ValueError("duplicate path or file outside selection")
            seen.add(path)
            paths.append(path)
            if type(row["size"]) is not int or row["size"] < 0 or type(row["mode"]) is not int or not 0 <= row["mode"] <= 0o777:
                raise ValueError("invalid file size or mode")
            if not isinstance(row["sha256"], str) or not HEX.fullmatch(row["sha256"]):
                raise ValueError("invalid blob digest")
            if field == "files":
                old = blobs.setdefault(row["sha256"], row["size"])
                if old != row["size"]:
                    raise ValueError("conflicting blob sizes")
            elif row["reason"] not in {"ELF", "PE executable", "Mach-O executable", "static archive", "executable binary"}:
                raise ValueError("invalid exclusion reason")
        if paths != sorted(paths):
            raise ValueError("file records are not ordered")
    # A file must never also be an ancestor directory of another entry.
    file_paths = {r["path"] for f in ("files", "excluded_files") for r in manifest[f]}
    for path in seen:
        parts = path.split("/")
        if any("/".join(parts[:i]) in file_paths for i in range(1, len(parts))):
            raise ValueError("file/directory path conflict")
    return manifest, blobs


def _mkdirs(root_fd, path):
    parent = os.dup(root_fd)
    try:
        for part in path.split("/") if path else ():
            try:
                os.mkdir(part, 0o700, dir_fd=parent)
            except FileExistsError:
                pass
            child = os.open(part, DIR_FLAGS, dir_fd=parent)
            os.close(parent)
            parent = child
    finally:
        os.close(parent)


def verify_archive(archive, restore=None):
    archive = Path(os.path.abspath(archive))
    parent_fd = open_root(archive.parent)
    try:
        fd = os.open(archive.name, os.O_RDONLY | os.O_NOFOLLOW, dir_fd=parent_fd)
    finally:
        os.close(parent_fd)
    with os.fdopen(fd, "rb") as raw, tempfile.TemporaryDirectory(prefix="round4-verify-") as staging:
        initial = os.fstat(raw.fileno())
        if not stat.S_ISREG(initial.st_mode):
            raise ValueError("archive must be a regular file")
        archive_hash = hashlib.sha256()
        while chunk := raw.read(CHUNK):
            archive_hash.update(chunk)
        raw.seek(0)
        manifest, expected, found = None, None, set()
        with gzip.GzipFile(fileobj=raw, mode="rb") as gz, tarfile.open(fileobj=gz, mode="r|") as tar:
            for member in tar:
                name = relative_path(member.name)
                if not member.isfile() or member.linkname or member.pax_headers:
                    raise ValueError("archive member is not a plain regular file")
                if manifest is None:
                    if name != "manifest.json" or member.size > MANIFEST_LIMIT:
                        raise ValueError("first member must be a bounded manifest.json")
                    data = tar.extractfile(member).read()
                    manifest, expected = decode_manifest(data)
                    if data != canonical_json(manifest):
                        raise ValueError("manifest is not canonical JSON")
                    continue
                match = re.fullmatch(r"blobs/sha256/([0-9a-f]{64})", name)
                if not match:
                    raise ValueError("unexpected archive member: " + name)
                digest = match.group(1)
                if digest in found or digest not in expected or member.size != expected[digest]:
                    raise ValueError("duplicate, unreferenced or incorrectly sized blob")
                found.add(digest)
                hasher = hashlib.sha256()
                source = tar.extractfile(member)
                target = open(Path(staging) / digest, "xb") if restore is not None else None
                try:
                    while chunk := source.read(CHUNK):
                        hasher.update(chunk)
                        if target is not None:
                            target.write(chunk)
                finally:
                    if target is not None:
                        target.close()
                if hasher.hexdigest() != digest:
                    raise ValueError("blob SHA256 mismatch")
            # Consume through gzip EOF/CRC and reject concatenated hidden data.
            padding = 0
            while chunk := tar.fileobj.read(CHUNK):
                padding += len(chunk)
                if padding > CHUNK or any(chunk):
                    raise ValueError("unexpected trailing archive content")
        if manifest is None or found != set(expected):
            raise ValueError("missing manifest or blobs")
        if signature(initial) != signature(os.fstat(raw.fileno())):
            raise ValueError("archive changed during verification")
        if restore is not None:
            destination = Path(os.path.abspath(restore))
            parent_fd = open_root(destination.parent)
            try:
                os.mkdir(destination.name, 0o700, dir_fd=parent_fd)  # exclusive; never replace
                restore_fd = os.open(destination.name, DIR_FLAGS, dir_fd=parent_fd)
            finally:
                os.close(parent_fd)
            try:
                for path in manifest["directories"]:
                    _mkdirs(restore_fd, path)
                for row in manifest["files"]:
                    _mkdirs(restore_fd, row["path"].rpartition("/")[0])
                    out_fd = open_relative(restore_fd, row["path"], os.O_WRONLY | os.O_CREAT | os.O_EXCL)
                    with os.fdopen(out_fd, "wb") as target, open(Path(staging) / row["sha256"], "rb") as source:
                        shutil.copyfileobj(source, target, CHUNK)
                        os.fchmod(target.fileno(), row["mode"])
            finally:
                os.close(restore_fd)
        return {"archive": str(archive), "archive_sha256": archive_hash.hexdigest(), "files": len(manifest["files"]),
                "excluded_files": len(manifest["excluded_files"]), "unique_blobs": len(expected),
                "source_bytes": sum(row["size"] for row in manifest["files"]), "blob_bytes": sum(expected.values())}


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("archive", type=Path)
    parser.add_argument("--restore", type=Path, help="new directory; existing paths are refused")
    args = parser.parse_args()
    try:
        print(json.dumps(verify_archive(args.archive, args.restore), sort_keys=True, indent=2))
    except (OSError, ValueError, tarfile.TarError, EOFError) as exc:
        parser.exit(1, "verification failed: " + str(exc) + "\n")


if __name__ == "__main__":
    main()
