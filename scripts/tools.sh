#!/usr/bin/env bash

set -e  # exit on failure

## Install golint.
mkdir -p "$GOBIN"
echo "Contents of $GOBIN:"
ls -l "$GOBIN" || true

if ! command -v golint > /dev/null; then
  rm -rf "${GOBIN}/golint"
  echo "Installing 'golint'..."
  GO111MODULE=off go get -u golang.org/x/lint/golint
fi
echo "golint: $(command -v golint)"

set +e
