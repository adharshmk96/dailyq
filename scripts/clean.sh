#!/usr/bin/env bash
# Remove build outputs and reset the embed directory to its placeholder.
set -euo pipefail

cd "$(dirname "$0")/.."

DIST_DIR="dailyq-api/internal/web/dist"

rm -rf dailyq-ui/.output dailyq-ui/.nuxt dailyq-api/bin
rm -rf "$DIST_DIR"
mkdir -p "$DIST_DIR"
touch "$DIST_DIR/.keep"

echo "==> cleaned"
