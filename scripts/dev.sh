#!/bin/zsh
set -euo pipefail
export ENV=dev
export ADDR=:8080
exec go run ./cmd/server
