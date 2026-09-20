#!/bin/bash
set -euo pipefail

NAME="socksserver"
VERSION=$(git describe --tags --always 2>/dev/null || echo "0.3")
BUILD_DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ')
LDFLAGS="-s -w -X main.version=${VERSION} -X main.buildDate=${BUILD_DATE}"

rm -rf pack pack.zip
mkdir -p pack

TARGETS=${1:-"all"}

if [ "$TARGETS" = "fast" ]; then
  build_list="linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64"
else
  build_list=$(go tool dist list)
fi

TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

for line in $build_list; do
  os=$(echo "$line" | cut -d'/' -f1)
  arch=$(echo "$line" | cut -d'/' -f2)

  # Filter unsupported platforms for server binaries
  if [ "$os" = "android" ] || [ "$os" = "ios" ] || [ "$arch" = "wasm" ] || [ "$os" = "js" ]; then
    continue
  fi

  echo "Building ${NAME} for ${os}/${arch}..."
  bin_name="${NAME}"
  if [ "$os" = "windows" ]; then
    bin_name="${NAME}.exe"
  fi

  out_path="${TMP_DIR}/${bin_name}"
  if ! CGO_ENABLED=0 GOOS="$os" GOARCH="$arch" go build -trimpath -ldflags="${LDFLAGS}" -o "${out_path}" .; then
    echo "Warning: ${os}/${arch} build failed, skipping"
    continue
  fi

  zip_name="${NAME}_${os}_${arch}.zip"
  (cd "$TMP_DIR" && zip -q "${zip_name}" "${bin_name}")
  mv "${TMP_DIR}/${zip_name}" pack/
  rm -f "${out_path}"
  echo "Done ${os}/${arch} -> pack/${zip_name}"
done

zip -q -r pack.zip pack/
echo "Packaging complete: pack.zip created."
