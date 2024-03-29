#!/bin/sh

# Borrowed from: https://raw.githubusercontent.com/goreleaser/goreleaser/main/scripts/completions.sh
set -e
rm -rf completions
mkdir completions

# Directly invoke uplift and generate the shell completion scripts
for SH in bash zsh fish; do
	go run ./cmd/uplift/... completion "${SH}" > "completions/uplift.${SH}"
done