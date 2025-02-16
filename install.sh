#!/bin/bash

REPO="Kamaliev/gossh"
VERSION="latest"
INSTALL_DIR="/usr/local/bin"

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

if [[ "$OS" == "darwin" ]]; then
  OS_NAME="mac"
else
  OS_NAME=$OS
fi

if [[ "$ARCH" == "x86_64" ]]; then
    ARCH="amd64"
elif [[ "$ARCH" == "arm64" || "$ARCH" == "aarch64" ]]; then
    ARCH="arm64"
else
    echo "Unsupported architecture: $ARCH"
    exit 1
fi

URL="https://github.com/$REPO/releases/$VERSION/download/gossh-${OS_NAME}-${ARCH}"
echo "⬇️  Скачивание $URL..."
curl -L -o gossh "$URL"

chmod +x gossh
sudo mv gossh "$INSTALL_DIR/gossh"

if [[ "$OS_NAME" == "mac" ]]; then sudo xattr -d com.apple.quarantine "$INSTALL_DIR/gossh"; fi

echo "✅ Установка завершена! Теперь можно использовать команду: gossh"
