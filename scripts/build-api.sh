#!/usr/bin/env bash
# Compile the server, embedding whatever is currently in internal/web/dist.
set -euo pipefail

cd "$(dirname "$0")/.."

echo "==> building server"
(cd dailyq-api && go build -o bin/dailyq .)

echo "==> binary at dailyq-api/bin/dailyq"
