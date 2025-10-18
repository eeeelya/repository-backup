#!/usr/bin/env bash

CLI_NAME="git-repository-backup"

OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

URL="https://github.com/eeeelya/repository-backup/releases/download/latest/${CLI_NAME}-${OS}-${ARCH}"

echo "Installing $CLI_NAME version 'latest' for $OS/$ARCH..."

curl -L -o /tmp/$CLI_NAME "$URL"
chmod +x /tmp/$CLI_NAME

sudo mv /tmp/$CLI_NAME /usr/local/bin/$CLI_NAME

echo "$CLI_NAME installed successfully!"
$CLI_NAME --version
