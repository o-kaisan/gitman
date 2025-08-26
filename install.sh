#!/usr/bin/env bash
set -euo pipefail

REPO="o-kai/gitman"
INSTALL_DIR="/usr/local/bin"
BINARY_PREFIX="gitman"

# 最新リリースのタグを取得
LATEST_TAG=$(curl -s https://api.github.com/$REPO/releases/latest | grep -Po '"tag_name": "\K.*?(?=")')
if [ -z "$LATEST_TAG" ]; then
  echo "❌ 最新リリースが見つかりません"
  exit 1
fi
echo "👉 最新リリース: $LATEST_TAG"

# OS/ARCH 判定
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  aarch64 | arm64) ARCH="arm64" ;;
  *) echo "❌ 未対応のアーキテクチャ: $ARCH" && exit 1 ;;
esac

BINARY_NAME="${BINARY_PREFIX}-${LATEST_TAG}"
TAR_NAME="${BINARY_NAME}-${OS}-${ARCH}.tar.gz"
URL="https://github.com/$REPO/releases/download/$LATEST_TAG/$TAR_NAME"

# ダウンロード
TMP_DIR=$(mktemp -d)
echo "⬇️  ダウンロード中: $URL"
curl -sL "$URL" -o "$TMP_DIR/$TAR_NAME"

# 展開 & インストール
tar -xzf "$TMP_DIR/$TAR_NAME" -C "$TMP_DIR"
sudo mv "$TMP_DIR/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"

# 確認
echo "✅ インストール完了: $INSTALL_DIR/$BINARY_NAME"
"$INSTALL_DIR/$BINARY_NAME" --version || echo "⚠️ バージョン確認に失敗しました"
