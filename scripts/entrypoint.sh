#!/usr/bin/env bash

# Dynamically load GPG key and passphrase, if provided
if [ -n "$GPG_KEY" ]; then
    gpg-agent --daemon --default-cache-ttl 21600

    if [ "$GPG_KEY" != "--"* ]; then
        GPG_KEY="$(echo $GPG_KEY | base64 -d)"
    fi

    # Import and activate the key ready for signing
    echo -e "$GPG_KEY" | gpg --batch --import --no-tty
    KEYID=$(gpg --with-colons --list-keys | awk -F: '/^pub/ { print $5 }')
    GPG_UID=$(gpg --with-colons --list-keys | awk -F: '/^uid/ { print $10 }')

    # Configure GPG by signing a temporary file
    echo "hello" > temp.txt
    gpg --batch --detach-sig -v --pinentry-mode loopback --passphrase "$GPG_PASSPHRASE" --no-tty temp.txt
    rm temp.*

    # Configure git globally to support GPG signing on both commits and tags
    git config --global user.signingkey "$KEYID"
    git config --global commit.gpgsign true
    git config --global tag.gpgsign true
    git config --global user.name "$(echo $GPG_UID | cut -d ' ' -f1)"
    git config --global user.email "$(echo $GPG_UID | cut -d ' ' -f2 | tr -d '<>')"
fi

exec uplift $@