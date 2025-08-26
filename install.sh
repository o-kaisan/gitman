#!/usr/bin/env bash
set -euo pipefail

REPO="o-kaisan/gitman"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="gitman"

# Get the latest release tag
LATEST_TAG=$(curl -s https://api.github.com/repos/$REPO/releases/latest | grep -Po '"tag_name": "\K.*?(?=")')
if [ -z "$LATEST_TAG" ]; then
  echo "❌ No latest release found"
  exit 1
fi
echo "👉 Latest release: $LATEST_TAG"

# Detect OS/ARCH
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  aarch64 | arm64) ARCH="arm64" ;;
  *) echo "❌ Unsupported architecture: $ARCH" && exit 1 ;;
esac

TAR_NAME="${BINARY_NAME}-${LATEST_TAG}-${OS}-${ARCH}.tar.gz"
URL="https://github.com/$REPO/releases/download/$LATEST_TAG/$TAR_NAME"

# Download with progress bar
TMP_DIR=$(mktemp -d)
echo "⬇️  Downloading: $URL"
curl -L --progress-bar "$URL" -o "$TMP_DIR/$TAR_NAME"

# Extract & install
echo "📦 Extracting package..."
tar -xzf "$TMP_DIR/$TAR_NAME" -C "$TMP_DIR"

# The extracted binary will have a tag name, rename it to 'gitman'
EXTRACTED_BINARY="$TMP_DIR/${BINARY_NAME}-${LATEST_TAG}"
if [ ! -f "$EXTRACTED_BINARY" ]; then
  echo "❌ Extracted binary not found: $EXTRACTED_BINARY"
  exit 1
fi

sudo mv "$EXTRACTED_BINARY" "$INSTALL_DIR/$BINARY_NAME"
sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"

# Verify
echo "✅ Installed: $INSTALL_DIR/$BINARY_NAME"
"$INSTALL_DIR/$BINARY_NAME" --version || echo "⚠️ Failed to verify version"

# Cleanup
rm -rf "$TMP_DIR"

echo "✅ Installation completed"
