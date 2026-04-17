#!/usr/bin/env sh
set -e

REPO="luismasuarez/wg-tui"
BIN="wg-tui"
DEST="${HOME}/.local/bin"

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) echo "Arquitectura no soportada: $ARCH"; exit 1 ;;
esac

LATEST=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | cut -d'"' -f4)
URL="https://github.com/${REPO}/releases/download/${LATEST}/${BIN}-${OS}-${ARCH}"

echo "Instalando wg-tui ${LATEST}..."
mkdir -p "$DEST"
curl -fsSL "$URL" -o "${DEST}/${BIN}"
chmod +x "${DEST}/${BIN}"
echo "Instalado en ${DEST}/${BIN}"
echo "Asegúrate de que ${DEST} esté en tu PATH."
