#!/usr/bin/env bash
workspace=$(cd $(dirname $0) && pwd -P)
# eg: 
# ./goLint.sh # lint all files
# ./goLint.sh DIR/... --fix # lint files in DIR and autofix
docker run --rm -v "$workspace":/app -w /app golangci/golangci-lint golangci-lint run $1 $2
