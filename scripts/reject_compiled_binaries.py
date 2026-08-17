#!/usr/bin/env python3
"""Reject compiled executables from entering git history.

At pre-commit, inspect the filenames supplied by prek against the working tree.
At pre-push, ignore the net-diff filenames and inspect every newly reachable
blob. PRE_COMMIT_FROM_REF supplies the existing-history boundary when present;
if prek omits it for an orphan-history push, scan the full ancestry of
PRE_COMMIT_TO_REF. This catches an executable added in one pushed commit and
deleted in a later one, and never consults a jj working tree that may differ
from the pushed ref.
"""
import os
import string
import subprocess
import sys

NULL_REF = "0" * 40

EXECUTABLE_MAGICS = {
    b"\x7fELF": "ELF executable",
    b"\xfe\xed\xfa\xce": "Mach-O executable (32-bit)",
    b"\xfe\xed\xfa\xcf": "Mach-O executable (64-bit)",
    b"\xce\xfa\xed\xfe": "Mach-O executable (32-bit, LE)",
    b"\xcf\xfa\xed\xfe": "Mach-O executable (64-bit, LE)",
    b"\xca\xfe\xba\xbe": "Mach-O universal binary",
    b"\xbe\xba\xfe\xca": "Mach-O universal binary (LE)",
}


class GitInspectionError(RuntimeError):
    pass


def classify_content(content: bytes) -> str | None:
    """Return a human label if content starts with executable magic bytes."""
    head = content[:4]
    for magic, label in EXECUTABLE_MAGICS.items():
        if head.startswith(magic):
            return label
    if head[:2] == b"MZ":
        return "PE/DOS executable"
    return None


def classify_worktree_path(path: str) -> tuple[str, str | None]:
    """Classify one pre-commit path from the working tree."""
    try:
        with open(path, "rb") as fh:
            label = classify_content(fh.read(4))
    except (OSError, IsADirectoryError):
        label = None
    return path, label


def run_git(args: list[str], *, input_text: str | None = None) -> str:
    """Run Git or raise so a broken history check fails closed."""
    try:
        result = subprocess.run(
            ["git", *args],
            input=input_text,
            capture_output=True,
            text=True,
            timeout=30,
        )
    except (OSError, subprocess.TimeoutExpired) as exc:
        raise GitInspectionError(f"git {' '.join(args)} failed: {exc}") from exc
    if result.returncode != 0:
        detail = result.stderr.strip() or f"exit {result.returncode}"
        raise GitInspectionError(f"git {' '.join(args)} failed: {detail}")
    return result.stdout


def is_null_ref(ref: str) -> bool:
    return bool(ref) and set(ref) == {"0"}


def newly_reachable_objects(from_ref: str, to_ref: str) -> list[tuple[str, str]]:
    """Return (object id, path hint) for objects pushed by from_ref..to_ref."""
    if is_null_ref(to_ref):  # deleting a remote ref uploads no objects
        return []

    args = ["rev-list", "--objects", to_ref]
    if not is_null_ref(from_ref):
        args.extend(["--not", from_ref])
    else:
        # A new remote ref has no from_ref. Exclude objects known to be present
        # under that remote's tracking refs; without a remote name, scan the
        # full ancestry of to_ref rather than silently skipping it.
        remote = os.environ.get("PRE_COMMIT_REMOTE_NAME")
        if remote:
            args.extend(["--not", f"--remotes={remote}"])

    objects: dict[str, str] = {}
    for line in run_git(args).splitlines():
        oid, separator, path = line.partition(" ")
        if len(oid) not in (40, 64) or any(ch not in string.hexdigits for ch in oid):
            continue
        objects.setdefault(oid, path if separator else "")
    return list(objects.items())


def object_types(objects: list[tuple[str, str]]) -> dict[str, str]:
    if not objects:
        return {}
    payload = "".join(f"{oid}\n" for oid, _ in objects)
    output = run_git(
        ["cat-file", "--batch-check=%(objectname) %(objecttype)"],
        input_text=payload,
    )
    types: dict[str, str] = {}
    for line in output.splitlines():
        oid, _, kind = line.partition(" ")
        if not kind or kind == "missing":
            raise GitInspectionError(f"git cat-file could not resolve object {oid}")
        types[oid] = kind
    return types


def read_blob_head(oid: str) -> bytes:
    """Read only the first four bytes of a known blob without buffering it."""
    try:
        proc = subprocess.Popen(
            ["git", "cat-file", "blob", oid],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
        )
        assert proc.stdout is not None
        head = proc.stdout.read(4)
        returncode = proc.poll()
        if returncode is None:
            proc.terminate()
            proc.wait(timeout=5)
        elif returncode != 0:
            assert proc.stderr is not None
            detail = proc.stderr.read().decode(errors="replace").strip()
            raise GitInspectionError(f"git cat-file blob {oid} failed: {detail}")
        return head
    except (OSError, subprocess.TimeoutExpired) as exc:
        raise GitInspectionError(f"git cat-file blob {oid} failed: {exc}") from exc


def classify_pushed_blobs(from_ref: str, to_ref: str) -> list[tuple[str, str]]:
    objects = newly_reachable_objects(from_ref, to_ref)
    types = object_types(objects)
    offenders = []
    for oid, path in objects:
        if types.get(oid) != "blob":
            continue
        if label := classify_content(read_blob_head(oid)):
            display = path or f"blob {oid}"
            offenders.append((f"{display} [blob {oid[:12]}]", label))
    return offenders


def report(offenders: list[tuple[str, str]], action: str) -> int:
    if not offenders:
        return 0
    sys.stderr.write(f"refusing to {action} compiled executable binaries:\n")
    for path, label in offenders:
        sys.stderr.write(f"  {path}  ({label})\n")
    if action == "commit":
        sys.stderr.write(
            "these look like build artifacts — unstage them with "
            "git rm --cached <file> and add an ignore rule.\n"
        )
    else:
        sys.stderr.write(
            "rewrite or remove the offending commit before pushing; deleting "
            "the file in a later commit does not remove its blob from history.\n"
        )
    return 1


def main(argv: list[str]) -> int:
    from_ref = os.environ.get("PRE_COMMIT_FROM_REF")
    to_ref = os.environ.get("PRE_COMMIT_TO_REF")

    if from_ref or to_ref:
        if not to_ref:
            sys.stderr.write(
                "reject_compiled_binaries: incomplete pre-push ref context "
                "(PRE_COMMIT_TO_REF is missing)\n"
            )
            return 2
        try:
            # Some pre-push adapters omit FROM_REF for orphan-history pushes.
            # Treat that as a new remote ref rather than falling back to the
            # working-tree pre-commit path and silently bypassing history scan.
            return report(
                classify_pushed_blobs(from_ref or NULL_REF, to_ref), "push"
            )
        except GitInspectionError as exc:
            sys.stderr.write(f"reject_compiled_binaries: {exc}\n")
            return 2

    offenders = [
        (path, label)
        for path in argv
        if (label := classify_worktree_path(path)[1]) is not None
    ]
    return report(offenders, "commit")


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))
