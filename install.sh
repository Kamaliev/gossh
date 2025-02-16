#!/bin/bash

REPO="Kamaliev/gossh"
VERSION="latest"
INSTALL_DIR="/usr/local/bin"

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

if [[ "$ARCH" == "x86_64" ]]; then
    ARCH="amd64"
elif [[ "$ARCH" == "arm64" || "$ARCH" == "aarch64" ]]; then
    ARCH="arm64"
else
    echo "Unsupported architecture: $ARCH"
    exit 1
fi

URL="https://github.com/$REPO/releases/$VERSION/download/gossh-${OS}-${ARCH}"
echo "⬇️  Скачивание $URL..."
curl -L -o gossh "$URL"

chmod +x gossh
sudo mv gossh "$INSTALL_DIR/gossh"

echo "✅ Установка завершена! Теперь можно использовать команду: gossh"
