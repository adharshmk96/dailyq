#!/usr/bin/env bash
# Run the API server from source, without rebuilding the UI.
set -euo pipefail

cd "$(dirname "$0")/../dailyq-api"

exec go run . "${@:-serve}"
