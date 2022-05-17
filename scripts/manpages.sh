#!/bin/sh

# Borrowed from: https://raw.githubusercontent.com/goreleaser/goreleaser/main/scripts/manpages.sh
set -e
rm -rf manpages
mkdir manpages
go run ./cmd/uplift/... man | gzip -c -9 > manpages/uplift.1.gz