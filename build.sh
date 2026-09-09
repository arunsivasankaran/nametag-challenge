#!/usr/bin/env bash
set -euo pipefail

OUT_DIR="dist"
mkdir -p "$OUT_DIR"

TARGETS=(
  "windows amd64"
  "linux amd64"
  "linux arm64"
  "darwin amd64"
  "darwin arm64"
)

for target in "${TARGETS[@]}"; do
  set -- $target
  GOOS=$1
  GOARCH=$2

  output_name="nametag-${GOOS}-${GOARCH}"
  if [[ "$GOOS" == "windows" ]]; then
    output_name+=".exe"
  fi

  echo "Building ${GOOS}/${GOARCH} -> ${OUT_DIR}/${output_name}"
  GOOS="$GOOS" GOARCH="$GOARCH" go build -o "${OUT_DIR}/${output_name}" .
done

echo "Build complete. Artifacts are in ${OUT_DIR}/"
