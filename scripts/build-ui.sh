#!/usr/bin/env bash
# Generate the static UI and copy it into the Go embed directory.
set -euo pipefail

cd "$(dirname "$0")/.."

UI_DIR="dailyq-ui"
DIST_DIR="dailyq-api/internal/web/dist"

echo "==> installing ui dependencies"
(cd "$UI_DIR" && bun install --frozen-lockfile)

echo "==> generating static ui"
(cd "$UI_DIR" && bun run generate)

echo "==> embedding build into $DIST_DIR"
rm -rf "$DIST_DIR"
mkdir -p "$DIST_DIR"
touch "$DIST_DIR/.keep"
cp -R "$UI_DIR/.output/public/." "$DIST_DIR/"
find "$DIST_DIR" -name .DS_Store -delete

echo "==> ui ready"
