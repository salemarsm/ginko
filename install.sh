#!/usr/bin/env bash
# Install ginko from GitHub Releases.
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/salemarsm/llm-memory/main/install.sh | bash
#   GINKO_VERSION=v0.2.1 bash install.sh
#   GINKO_INSTALL_DIR=/usr/local/bin bash install.sh

set -euo pipefail

REPO="salemarsm/llm-memory"
INSTALL_DIR="${GINKO_INSTALL_DIR:-$HOME/.local/bin}"
VERSION="${GINKO_VERSION:-}"

# --- OS detection ---
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS" in
  linux|darwin) ;;
  *) echo "error: unsupported OS: $OS" >&2; exit 1 ;;
esac

# --- Arch detection ---
ARCH=$(uname -m)
case "$ARCH" in
  x86_64|amd64)   ARCH="amd64" ;;
  aarch64|arm64)  ARCH="arm64" ;;
  *) echo "error: unsupported architecture: $ARCH" >&2; exit 1 ;;
esac

# --- Resolve version ---
if [ -z "$VERSION" ]; then
  echo "Fetching latest release..."
  VERSION=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" \
    | grep '"tag_name"' \
    | sed 's/.*"tag_name": *"\([^"]*\)".*/\1/')
fi
[ -z "$VERSION" ] && { echo "error: could not resolve version (no releases published yet?)" >&2; exit 1; }

echo "Installing ginko ${VERSION} (${OS}/${ARCH}) → ${INSTALL_DIR}"

# --- Download ---
ARCHIVE="ginko_${VERSION#v}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/$REPO/releases/download/$VERSION/$ARCHIVE"
TMP=$(mktemp -d)
trap 'rm -rf "$TMP"' EXIT

if ! curl -fsSL "$URL" -o "$TMP/ginko.tar.gz"; then
  echo "" >&2
  echo "error: download failed: $URL" >&2
  echo "Check that the release exists: https://github.com/$REPO/releases" >&2
  exit 1
fi

tar -xzf "$TMP/ginko.tar.gz" -C "$TMP"

# --- Install binaries ---
mkdir -p "$INSTALL_DIR"
installed=()
for bin in ginko llm-memory memctl memmcp memserver; do
  if [ -f "$TMP/$bin" ]; then
    install -m 755 "$TMP/$bin" "$INSTALL_DIR/$bin"
    installed+=("$bin")
  fi
done

echo "Installed: ${installed[*]}"
echo "Location:  $INSTALL_DIR"

# --- Checksum verification (if checksums.txt present) ---
if [ -f "$TMP/checksums.txt" ]; then
  if command -v sha256sum >/dev/null 2>&1; then
    (cd "$TMP" && sha256sum -c checksums.txt --ignore-missing --quiet 2>/dev/null) \
      && echo "Checksums verified." \
      || echo "warning: checksum mismatch — verify manually" >&2
  fi
fi

# --- PATH hint ---
case ":${PATH}:" in
  *":${INSTALL_DIR}:"*) ;;
  *)
    echo ""
    echo "Add ${INSTALL_DIR} to your PATH:"
    if [ -n "${ZSH_VERSION:-}" ] || [ "$(basename "${SHELL:-}")" = "zsh" ]; then
      echo "  echo 'export PATH=\"\$HOME/.local/bin:\$PATH\"' >> ~/.zshrc && source ~/.zshrc"
    else
      echo "  echo 'export PATH=\"\$HOME/.local/bin:\$PATH\"' >> ~/.bashrc && source ~/.bashrc"
    fi
    ;;
esac

echo ""
echo "Next step:"
echo "  ginko setup claude-code"
