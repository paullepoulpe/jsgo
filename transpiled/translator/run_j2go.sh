#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
JAVA_SRC_DIR="${SCRIPT_DIR}/java"
BUILD_DIR="${J2GO_BUILD_DIR:-${TMPDIR:-/tmp}/j2go_build_classes}"

if [[ -x "/opt/homebrew/opt/openjdk/bin/javac" ]]; then
  JAVAC="/opt/homebrew/opt/openjdk/bin/javac"
  JAVA="/opt/homebrew/opt/openjdk/bin/java"
else
  JAVAC="$(command -v javac)"
  JAVA="$(command -v java)"
fi

mkdir -p "${BUILD_DIR}"

needs_compile=0
if [[ ! -f "${BUILD_DIR}/j2go/Main.class" ]]; then
  needs_compile=1
else
  for f in "${JAVA_SRC_DIR}"/j2go/*.java; do
    if [[ "${f}" -nt "${BUILD_DIR}/j2go/Main.class" ]]; then
      needs_compile=1
      break
    fi
  done
fi

if [[ "${needs_compile}" -eq 1 ]]; then
  rm -rf "${BUILD_DIR}/j2go"
  "${JAVAC}" --release 21 -encoding UTF-8 -d "${BUILD_DIR}" "${JAVA_SRC_DIR}"/j2go/*.java
fi

exec "${JAVA}" -Xmx4g -cp "${BUILD_DIR}" j2go.Main "$@"
