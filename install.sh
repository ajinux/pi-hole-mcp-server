#!/bin/sh
set -e

REPO="ajinux/pi-hole-mcp-server"
VERSION="v0.1.0"
BIN_NAME="pihole-mcp"
INSTALL_DIR="/usr/local/bin"

ARCH=$(uname -m)
OS=$(uname -s | tr '[:upper:]' '[:lower:]')

case "$ARCH" in
  x86_64) TARGET="linux-amd64" ;;
  aarch64) TARGET="linux-arm64" ;;
  armv7l|armv6l) TARGET="linux-armv7" ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

URL="https://github.com/$REPO/releases/download/$VERSION/$BIN_NAME-$TARGET"

echo "Downloading $BIN_NAME ($TARGET)..."
curl -fsSL "$URL" -o "/tmp/$BIN_NAME"

chmod +x "/tmp/$BIN_NAME"
sudo mv "/tmp/$BIN_NAME" "$INSTALL_DIR/$BIN_NAME"

echo "Installed $BIN_NAME to $INSTALL_DIR/$BIN_NAME"
