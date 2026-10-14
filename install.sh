#!/usr/bin/env bash
set -e

REPO="hamimLohani/snake"
INSTALL_DIR="/usr/local/bin"

if [ ! -w "$INSTALL_DIR" ]; then
    INSTALL_DIR="$HOME/.local/bin"
    mkdir -p "$INSTALL_DIR"
fi

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$ARCH" in
    x86_64|amd64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

case "$OS" in
    linux|darwin) ;;
    *) echo "Unsupported operating system: $OS"; exit 1 ;;
esac

BINARY="snake-${OS}-${ARCH}"
URL="https://github.com/${REPO}/releases/latest/download/${BINARY}"

echo "==> Downloading Snake CLI for ${OS}/${ARCH}..."
curl -fsSL "$URL" -o "${INSTALL_DIR}/snake"
chmod +x "${INSTALL_DIR}/snake"

echo "==> Snake CLI installed successfully to ${INSTALL_DIR}/snake!"
echo "    Run 'snake' to start playing."
