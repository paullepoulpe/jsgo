#!/usr/bin/env bash
# apply_prep.sh: Deterministically applies PRs 1-30 patches to .cache/closure-src/
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
CLOSURE_SRC="$REPO_ROOT/.cache/closure-src"
CLOSURE_SRC_JAR="$REPO_ROOT/.cache/closure-src.jar"
PATCHES_DIR="$SCRIPT_DIR/patches"
STAMP_FILE="$CLOSURE_SRC/.prep_applied"

usage() {
  echo "Usage: $0 [--reset] [--check] [--verify]"
  echo "  --reset   Restore pristine .cache/closure-src from .cache/closure-src.jar before patching"
  echo "  --check   Dry-run check against pristine .cache/closure-src.jar: verify all patches apply cleanly"
  echo "  --verify  Run OpenJDK 27 full compilation and Java CLI parity verification after patching"
  exit 1
}

RESET=0
CHECK_ONLY=0
VERIFY=0

while [[ $# -gt 0 ]]; do
  case "$1" in
    --reset) RESET=1; shift ;;
    --check) CHECK_ONLY=1; shift ;;
    --verify) VERIFY=1; shift ;;
    -h|--help) usage ;;
    *) echo "Unknown option: $1"; usage ;;
  esac
done

if [[ ! -f "$CLOSURE_SRC_JAR" ]]; then
  echo "ERROR: Missing upstream source jar: $CLOSURE_SRC_JAR" >&2
  exit 1
fi

PATCH_FILES=()
while IFS= read -r line; do
  PATCH_FILES+=("$line")
done < <(ls "$PATCHES_DIR"/*.patch | sort)

if [[ ${#PATCH_FILES[@]} -eq 0 ]]; then
  echo "ERROR: No .patch files found in $PATCHES_DIR" >&2
  exit 1
fi
echo "==> Found ${#PATCH_FILES[@]} patch files."

compute_prep_fingerprint() {
  shasum -a 256 "${PATCH_FILES[@]}" "$SCRIPT_DIR"/com/google/javascript/jscomp/*.java | shasum -a 256 | awk '{print $1}'
}

apply_patches_to_dir() {
  local target_dir="$1"
  for patch_file in "${PATCH_FILES[@]}"; do
    local patch_name
    patch_name="$(basename "$patch_file")"
    echo "==> Applying $patch_name..."
    local out
    if ! out="$(patch -d "$target_dir" -p1 -f -N < "$patch_file" 2>&1)"; then
      echo "ERROR: Failed to apply $patch_name:" >&2
      echo "$out" >&2
      exit 1
    fi
    local patched_count
    patched_count="$(printf '%s\n' "$out" | grep -c "^patching file " || true)"
    if [[ "$patched_count" -eq 0 ]]; then
      echo "ERROR: Patch $patch_name modified 0 files!" >&2
      echo "$out" >&2
      exit 1
    fi
  done

  # Ensure pure-Java prep helper classes in transpiled/prep/com/ are synced into target_dir
  if [[ -d "$SCRIPT_DIR/com" ]]; then
    cp -R "$SCRIPT_DIR/com/." "$target_dir/com/"
  fi
}

if [[ $CHECK_ONLY -eq 1 ]]; then
  echo "==> Running dry-run check of all ${#PATCH_FILES[@]} patches against pristine closure-src.jar..."
  CHECK_DIR="$(mktemp -d "${TMPDIR:-/tmp}/closure-check.XXXXXX")"
  trap 'rm -rf "$CHECK_DIR"' EXIT
  unzip -qo "$CLOSURE_SRC_JAR" -d "$CHECK_DIR"
  apply_patches_to_dir "$CHECK_DIR"
  echo "==> All ${#PATCH_FILES[@]} patches apply cleanly (dry-run passed)."
  exit 0
fi

EXPECTED_FP="$(compute_prep_fingerprint)"

if [[ $RESET -eq 1 ]]; then
  echo "==> Resetting .cache/closure-src to pristine upstream state..."
  rm -rf "$CLOSURE_SRC"
  mkdir -p "$CLOSURE_SRC"
  unzip -qo "$CLOSURE_SRC_JAR" -d "$CLOSURE_SRC"
  echo "==> Reset complete."
fi

if [[ ! -d "$CLOSURE_SRC/com" ]]; then
  echo "==> Extracting pristine .cache/closure-src from $CLOSURE_SRC_JAR..."
  mkdir -p "$CLOSURE_SRC"
  unzip -qo "$CLOSURE_SRC_JAR" -d "$CLOSURE_SRC"
fi

if [[ -f "$STAMP_FILE" ]]; then
  CURRENT_FP="$(cat "$STAMP_FILE")"
  if [[ "$CURRENT_FP" == "$EXPECTED_FP" ]]; then
    echo "==> Prep patches already applied and verified up-to-date (.prep_applied = ${CURRENT_FP:0:12}...)."
  else
    echo "ERROR: $CLOSURE_SRC has stale prep patches applied (.prep_applied mismatch). Re-run with --reset." >&2
    exit 1
  fi
else
  echo "==> Applying all ${#PATCH_FILES[@]} patches to $CLOSURE_SRC..."
  apply_patches_to_dir "$CLOSURE_SRC"
  printf '%s\n' "$EXPECTED_FP" > "$STAMP_FILE"
  echo "==> All ${#PATCH_FILES[@]} patches applied successfully (stamp ${EXPECTED_FP:0:12}...)."
fi

if [[ $VERIFY -eq 1 ]]; then
  "$SCRIPT_DIR/verify_prep.sh"
fi
