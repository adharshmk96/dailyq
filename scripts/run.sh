#!/usr/bin/env bash
# Serve the embedded UI and the API from the built binary.
set -euo pipefail

cd "$(dirname "$0")/.."

exec dailyq-api/bin/dailyq "${@:-serve}"
