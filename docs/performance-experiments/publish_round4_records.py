#!/usr/bin/env python3
"""Freeze explicitly selected COMPLETED Round4 result paths, append-only.

Usage: python3 publish_round4_records.py --root /workspace/jsgo-results/round4 \
  --include-prefix COMPLETED_DATASET --include-prefix ANOTHER_DATASET \
  --output NEW_ARCHIVE.tar.gz

Every selected regular file is preserved byte-for-byte, including JS, maps and
profiles. Executable binaries/static archives are excluded with hashes/sizes.
The archive contains canonical manifest.json plus SHA256-deduplicated blobs.
No timestamps are embedded in archive headers. An identical selection produces
identical bytes. Existing output paths are never overwritten. No archived code
is executed. Run --self-test for low-resource synthetic security/roundtrip tests.
"""

import argparse
import gzip
import hashlib
import io
import json
import os
from pathlib import Path
import re
import stat
import tarfile
import tempfile
import unittest

import verify_round4_records as records

# Scan bytes without printing the matching credential. This is a conservative
# tripwire, not a claim to recognize every possible secret representation.
CREDENTIALS = [
    ("private key", re.compile(rb"-----BEGIN (?:RSA |EC |DSA |OPENSSH |ENCRYPTED )?PRIVATE KEY-----")),
    ("AWS access key", re.compile(rb"\b(?:AKIA|ASIA)[A-Z0-9]{16}\b")),
    ("GitHub token", re.compile(rb"\b(?:gh[pousr]_[A-Za-z0-9]{30,}|github_pat_[A-Za-z0-9_]{40,})")),
    ("Slack token", re.compile(rb"\bxox[baprs]-[A-Za-z0-9-]{20,}")),
    ("API secret token", re.compile(rb"\bsk-(?:proj-|svcacct-)?[A-Za-z0-9_-]{24,}")),
    ("bearer credential", re.compile(rb"(?i)\bauthorization\s*[:=]\s*[\"']?bearer\s+[A-Za-z0-9_./+=-]{16,}")),
    ("assigned credential", re.compile(rb"(?i)[\"']?\b(?:api[_-]?key|access[_-]?token|client[_-]?secret|password)[\"']?\s*[:=]\s*[\"'][A-Za-z0-9_./+=-]{20,}")),
]


def credential_scan(chunk, path):
    for label, pattern in CREDENTIALS:
        if pattern.search(chunk):
            raise ValueError("credential pattern (" + label + ") in " + repr(path))


def binary_reason(sample, mode):
    if sample.startswith(b"\x7fELF"):
        return "ELF"
    if sample.startswith(b"MZ"):
        return "PE executable"
    if sample[:4] in {b"\xfe\xed\xfa\xce", b"\xce\xfa\xed\xfe", b"\xfe\xed\xfa\xcf", b"\xcf\xfa\xed\xfe", b"\xca\xfe\xba\xbe", b"\xbe\xba\xfe\xca"}:
        return "Mach-O executable"
    if sample.startswith(b"!<arch>\n"):
        return "static archive"
    if mode & 0o111 and b"\0" in sample and not sample.startswith(b"#!"):
        return "executable binary"
    return None


def inventory(root_fd, prefixes):
    """Inspect selected paths via directory FDs; reject links/special files."""
    result = {}
    def visit(path, fd):
        current = os.fstat(fd)
        if not (stat.S_ISDIR(current.st_mode) or stat.S_ISREG(current.st_mode)):
            raise ValueError("selected path is not a regular file/directory: " + repr(path))
        if path != ".":
            records.relative_path(path)
        previous = result.setdefault(path, current)
        if records.signature(previous) != records.signature(current):
            raise ValueError("selection changed during inventory")
        if stat.S_ISDIR(current.st_mode):
            with os.scandir(fd) as children:
                names = sorted(entry.name for entry in children)
            for name in names:
                child_path = name if path == "." else path + "/" + name
                records.relative_path(child_path)
                # O_NONBLOCK avoids hanging on a FIFO before fstat rejects it.
                child = os.open(name, os.O_RDONLY | os.O_NOFOLLOW | os.O_NONBLOCK, dir_fd=fd)
                try:
                    visit(child_path, child)
                finally:
                    os.close(child)
    for prefix in prefixes:
        fd = os.dup(root_fd) if prefix == "." else records.open_relative(root_fd, prefix, os.O_RDONLY | os.O_NONBLOCK)
        try:
            visit(prefix, fd)
        finally:
            os.close(fd)
    return result


def _add(tar, name, size, source):
    member = tarfile.TarInfo(name)
    member.size, member.mode = size, 0o444
    member.uid = member.gid = member.mtime = 0
    member.uname = member.gname = ""
    tar.addfile(member, source)


def publish(root, prefixes, output):
    root, output = Path(os.path.abspath(root)), Path(os.path.abspath(output))
    prefixes = sorted(set(records.relative_path(p, allow_root=True) for p in prefixes))
    if not prefixes:
        raise ValueError("at least one explicit include prefix is required")
    # A snapshot must not ingest its own staging/output tree.
    if output == root or root in output.parents:
        raise ValueError("output must be outside the result root")
    output_parent = records.open_root(output.parent)
    root_fd = None
    temp_name = None
    temp_created = False
    try:
        root_fd = records.open_root(root)
        try:
            os.stat(output.name, dir_fd=output_parent, follow_symlinks=False)
        except FileNotFoundError:
            pass
        else:
            raise FileExistsError("output already exists: " + str(output))
        before = inventory(root_fd, prefixes)
        with tempfile.TemporaryDirectory(prefix="round4-snapshot-") as staging:
            staging = Path(staging)
            files, excluded = [], []
            for path, observed in sorted(before.items()):
                if stat.S_ISDIR(observed.st_mode):
                    continue
                fd = records.open_relative(root_fd, path, os.O_RDONLY | os.O_NONBLOCK)
                with os.fdopen(fd, "rb") as source, tempfile.TemporaryFile(dir=staging) as copied:
                    if records.signature(observed) != records.signature(os.fstat(source.fileno())):
                        raise ValueError("file changed before snapshot: " + repr(path))
                    hasher, size, tail, reason = hashlib.sha256(), 0, b"", None
                    while chunk := source.read(records.CHUNK):
                        if size == 0:
                            reason = binary_reason(chunk[:4096], observed.st_mode)
                        if reason is None:
                            credential_scan(tail + chunk, path)
                            tail = (tail + chunk)[-4096:]
                            copied.write(chunk)
                        hasher.update(chunk)
                        size += len(chunk)
                    if size != observed.st_size or records.signature(observed) != records.signature(os.fstat(source.fileno())):
                        raise ValueError("file changed during snapshot: " + repr(path))
                    digest = hasher.hexdigest()
                    row = {"path": path, "size": size, "sha256": digest, "mode": stat.S_IMODE(observed.st_mode) & 0o777}
                    if reason is not None:
                        row["reason"] = reason
                        excluded.append(row)
                    else:
                        files.append(row)
                        blob = staging / digest
                        if not blob.exists():
                            copied.seek(0)
                            with blob.open("xb") as target:
                                while chunk := copied.read(records.CHUNK):
                                    target.write(chunk)
            after = inventory(root_fd, prefixes)
            if set(before) != set(after) or any(records.signature(s) != records.signature(after[p]) for p, s in before.items()):
                raise ValueError("selection changed during snapshot")
            manifest = {"format": records.FORMAT, "hash_algorithm": "sha256", "source_root": str(root),
                        "include_prefixes": prefixes, "directories": sorted(p for p, s in before.items() if p != "." and stat.S_ISDIR(s.st_mode)),
                        "files": files, "excluded_files": excluded}
            data = records.canonical_json(manifest)
            if len(data) > records.MANIFEST_LIMIT:
                raise ValueError("manifest exceeds verifier size limit")
            _, blobs = records.decode_manifest(data)
            # Create using the anchored output directory. Hard-link publication
            # is atomic and fails if *any* destination entry already exists.
            temp_name = ".round4-publish-" + os.urandom(12).hex()
            archive_fd = os.open(temp_name, os.O_WRONLY | os.O_CREAT | os.O_EXCL | os.O_NOFOLLOW, 0o600, dir_fd=output_parent)
            temp_created = True
            with os.fdopen(archive_fd, "wb") as raw:
                with gzip.GzipFile(filename="", fileobj=raw, mode="wb", mtime=0, compresslevel=9) as gz:
                    with tarfile.open(fileobj=gz, mode="w|", format=tarfile.USTAR_FORMAT) as tar:
                        _add(tar, "manifest.json", len(data), io.BytesIO(data))
                        for digest, size in sorted(blobs.items()):
                            with (staging / digest).open("rb") as source:
                                _add(tar, "blobs/sha256/" + digest, size, source)
                raw.flush()
                os.fsync(raw.fileno())
                os.fchmod(raw.fileno(), 0o444)
            final = inventory(root_fd, prefixes)
            if set(before) != set(final) or any(records.signature(s) != records.signature(final[p]) for p, s in before.items()):
                raise ValueError("selection changed while building archive")
            os.link(temp_name, output.name, src_dir_fd=output_parent, dst_dir_fd=output_parent, follow_symlinks=False)
            os.fsync(output_parent)
            return {"archive": str(output), "files": len(files), "excluded_files": len(excluded), "unique_blobs": len(blobs),
                    "source_bytes": sum(row["size"] for row in files), "blob_bytes": sum(blobs.values())}
    finally:
        if temp_created:
            os.unlink(temp_name, dir_fd=output_parent)
        if root_fd is not None:
            os.close(root_fd)
        os.close(output_parent)


class SyntheticTests(unittest.TestCase):
    def setUp(self):
        self.temporary = tempfile.TemporaryDirectory(prefix="round4-record-test-")
        self.base = Path(self.temporary.name)
        self.root = self.base / "source"
        (self.root / "completed" / "empty-dir").mkdir(parents=True)
        (self.root / "completed" / "empty").write_bytes(b"")
        (self.root / "completed" / "a.js").write_bytes(b"compiled JS\n")
        (self.root / "completed" / "a.js.map").write_bytes(b"compiled JS\n")
        (self.root / "completed" / "cpu.pprof").write_bytes(b"\x00\x01profile")
        (self.root / "completed" / "compiler").write_bytes(b"\x7fELFfake-binary")
        (self.root / "unfinished").write_bytes(b"must not be included")
        self.archive = self.base / "one.tar.gz"

    def tearDown(self):
        self.temporary.cleanup()

    def test_roundtrip_dedup_exclusions_and_determinism(self):
        result = publish(self.root, ["completed"], self.archive)
        self.assertEqual((result["files"], result["excluded_files"], result["unique_blobs"]), (4, 1, 3))
        restored = self.base / "restored"
        verified = records.verify_archive(self.archive, restored)
        self.assertEqual(verified["files"], 4)
        for source in (self.root / "completed").iterdir():
            target = restored / "completed" / source.name
            if source.is_dir():
                self.assertTrue(target.is_dir())
            elif source.name == "compiler":
                self.assertFalse(target.exists())
            else:
                self.assertEqual(source.read_bytes(), target.read_bytes())
        self.assertFalse((restored / "unfinished").exists())
        second = self.base / "two.tar.gz"
        publish(self.root, ["completed"], second)
        self.assertEqual(self.archive.read_bytes(), second.read_bytes())
        with self.assertRaises(FileExistsError):
            publish(self.root, ["completed"], self.archive)
        with self.assertRaises(FileExistsError):
            records.verify_archive(self.archive, restored)

    def test_links_traversal_and_credentials_rejected(self):
        for path in ("../outside", "/absolute", "completed/../unfinished", "a\\b", "a//b", "C:/drive"):
            with self.assertRaises(ValueError):
                publish(self.root, [path], self.archive)
        link = self.root / "completed" / "linked"
        link.symlink_to(self.root / "unfinished")
        with self.assertRaises(OSError):
            publish(self.root, ["completed"], self.archive)
        link.unlink()
        (self.root / "completed" / "secret").write_bytes(b"-----BEGIN PRIVATE KEY-----")
        with self.assertRaisesRegex(ValueError, "credential pattern"):
            publish(self.root, ["completed"], self.archive)
        self.assertFalse(self.archive.exists())

    def test_boundary_secret_and_existing_symlink_output_rejected(self):
        secret = b"-----BEGIN PRIVATE KEY-----"
        (self.root / "completed" / "boundary-secret").write_bytes(b" " * (records.CHUNK - 8) + secret)
        with self.assertRaisesRegex(ValueError, "credential pattern"):
            publish(self.root, ["completed"], self.archive)
        (self.root / "completed" / "boundary-secret").unlink()
        destination = self.base / "untouched"
        destination.write_bytes(b"untouched")
        self.archive.symlink_to(destination)
        with self.assertRaises(FileExistsError):
            publish(self.root, ["completed"], self.archive)
        self.assertEqual(destination.read_bytes(), b"untouched")

    def test_tree_changes_while_archiving_rejected(self):
        original_add = _add
        changed = False
        def mutate(*args):
            nonlocal changed
            original_add(*args)
            if not changed:
                (self.root / "completed" / "new-result").write_bytes(b"concurrent writer")
                changed = True
        globals()["_add"] = mutate
        try:
            with self.assertRaisesRegex(ValueError, "changed while building"):
                publish(self.root, ["completed"], self.archive)
        finally:
            globals()["_add"] = original_add
        self.assertFalse(self.archive.exists())

    def test_changed_snapshot_rejected(self):
        original_scan = credential_scan
        def mutate(data, path):
            original_scan(data, path)
            if path.endswith("a.js"):
                (self.root / path).write_bytes(b"changed during reading")
        globals()["credential_scan"] = mutate
        try:
            with self.assertRaisesRegex(ValueError, "changed"):
                publish(self.root, ["completed"], self.archive)
        finally:
            globals()["credential_scan"] = original_scan
        self.assertFalse(self.archive.exists())

    def test_malicious_members_manifests_and_corrupt_blobs_rejected(self):
        publish(self.root, ["completed"], self.archive)
        with tarfile.open(self.archive, "r:gz") as tar:
            manifest = json.loads(tar.extractfile("manifest.json").read())
            contents = [(m.name, tar.extractfile(m).read()) for m in tar if m.name != "manifest.json"]
        def bad_archive(name, payload, members=()):
            path = self.base / name
            with tarfile.open(path, "w:gz") as tar:
                _add(tar, "manifest.json", len(payload), io.BytesIO(payload))
                for member_name, content in members:
                    _add(tar, member_name, len(content), io.BytesIO(content))
            return path
        for i, member in enumerate(("../escape", "/escape", "blobs/../escape")):
            bad = bad_archive("traversal" + str(i), records.canonical_json(manifest), [(member, b"bad")])
            with self.assertRaises(ValueError):
                records.verify_archive(bad, self.base / ("restore" + str(i)))
            self.assertFalse((self.base / ("restore" + str(i))).exists())
        manifest["files"][0]["path"] = "../escape"
        with self.assertRaises(ValueError):
            records.verify_archive(bad_archive("manifest-escape", records.canonical_json(manifest)))
        with tarfile.open(self.archive, "r:gz") as tar:
            data = tar.extractfile("manifest.json").read()
        corrupt = [(name, bytes([body[0] ^ 1]) + body[1:] if body else body) for name, body in contents]
        with self.assertRaisesRegex(ValueError, "SHA256"):
            records.verify_archive(bad_archive("corrupt", data, corrupt))
        for kind in (tarfile.SYMTYPE, tarfile.LNKTYPE, tarfile.DIRTYPE):
            bad = self.base / ("member-type-" + kind.decode())
            with tarfile.open(bad, "w:gz") as tar:
                _add(tar, "manifest.json", len(data), io.BytesIO(data))
                member = tarfile.TarInfo("blobs/sha256/" + "0" * 64)
                member.type, member.linkname = kind, "../escape"
                tar.addfile(member)
            with self.assertRaises(ValueError):
                records.verify_archive(bad)
        duplicate = bad_archive("duplicate-blob", data, contents + contents[:1])
        with self.assertRaisesRegex(ValueError, "duplicate"):
            records.verify_archive(duplicate)
        trailing = self.base / "trailing-gzip"
        trailing.write_bytes(self.archive.read_bytes() + gzip.compress(b"hidden second stream", mtime=0))
        with self.assertRaisesRegex(ValueError, "trailing"):
            records.verify_archive(trailing)


def main():
    parser = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    parser.add_argument("--root", type=Path)
    parser.add_argument("--include-prefix", action="append", default=[], help="explicit completed relative file/directory; repeatable")
    parser.add_argument("--output", type=Path)
    parser.add_argument("--self-test", action="store_true")
    args = parser.parse_args()
    if args.self_test:
        result = unittest.TextTestRunner(verbosity=2).run(unittest.defaultTestLoader.loadTestsFromTestCase(SyntheticTests))
        raise SystemExit(not result.wasSuccessful())
    if args.root is None or args.output is None or not args.include_prefix:
        parser.error("--root, --output and at least one --include-prefix are required")
    try:
        print(json.dumps(publish(args.root, args.include_prefix, args.output), sort_keys=True, indent=2))
    except (OSError, ValueError, tarfile.TarError) as exc:
        parser.exit(1, "publication failed: " + str(exc) + "\n")


if __name__ == "__main__":
    main()
