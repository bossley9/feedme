#!/bin/sh

EXE=$1
env="./.env"

get_env() {
  local cmd_args=""
  while read line; do
    key="${line%%=*}"
    value="${line##*=}"
    cmd_args="${cmd_args} -X main.${key}=${value}"
  done < "$env"
  echo $cmd_args
}

args=""

if [ -f "$env" ]; then
  go build -ldflags "$(get_env)" -o "$EXE" ./cmd/*.go
else
  echo "warning: env file not found. Continuing with defaults."
  go build -o "$EXE" ./cmd/*.go
fi
