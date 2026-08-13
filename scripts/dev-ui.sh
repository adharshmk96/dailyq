#!/usr/bin/env bash
# Nuxt dev server. /api is proxied to the local API server (see nuxt.config.ts).
set -euo pipefail

cd "$(dirname "$0")/../dailyq-ui"

exec bun run dev
