#!/bin/bash
#
# install.sh — one-click installer for mini_launch
# Usage: curl -fsSL https://raw.githubusercontent.com/simpossible/mini_launch/main/install.sh | bash
#

set -euo pipefail

REPO="simpossible/mini_launch"
BINARY="mini_launch"

# ---------- colors ----------
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

info()  { echo -e "${GREEN}[INFO]${NC} $*"; }
warn()  { echo -e "${YELLOW}[WARN]${NC} $*"; }
error() { echo -e "${RED}[ERROR]${NC} $*" >&2; exit 1; }

# ---------- detect OS ----------
detect_os() {
    local os="$(uname -s | tr '[:upper:]' '[:lower:]')"
    case "$os" in
        linux)  echo "linux" ;;
        darwin) echo "darwin" ;;
        *)      error "Unsupported OS: $os" ;;
    esac
}

# ---------- detect architecture ----------
detect_arch() {
    local arch="$(uname -m)"
    case "$arch" in
        x86_64|amd64) echo "amd64" ;;
        aarch64|arm64) echo "arm64" ;;
        *)             error "Unsupported architecture: $arch" ;;
    esac
}

# ---------- get latest version from GitHub ----------
get_latest_version() {
    local url="https://api.github.com/repos/${REPO}/releases/latest"
    local version

    version=$(curl -fsSL "$url" 2>/dev/null | grep '"tag_name"' | head -1 | sed -E 's/.*"([^"]+)".*/\1/')

    if [ -z "$version" ]; then
        error "Failed to get latest version from GitHub. Check your network connection."
    fi

    echo "$version"
}

# ---------- get currently installed version ----------
get_installed_version() {
    if command -v "$BINARY" &>/dev/null; then
        "$BINARY" --version 2>/dev/null | head -1 | grep -oE '[0-9]+\.[0-9]+\.[0-9]+' || echo ""
    fi
}

# ---------- determine install directory ----------
get_install_dir() {
    if [ -w "/usr/local/bin" ]; then
        echo "/usr/local/bin"
    else
        echo "${HOME}/.local/bin"
    fi
}

# ---------- download and install ----------
install() {
    local os="$1"
    local arch="$2"
    local version="$3"
    local install_dir="$4"

    local stripped_version="${version#v}"
    local archive_name="${BINARY}_${stripped_version}_${os}_$( [ "$arch" = "amd64" ] && echo "x86_64" || echo "$arch" ).tar.gz"
    local download_url="https://github.com/${REPO}/releases/download/${version}/${archive_name}"

    local tmp_dir
    tmp_dir=$(mktemp -d)
    trap 'rm -rf "$tmp_dir"' EXIT

    info "Downloading ${archive_name}..."
    curl -fsSL -o "${tmp_dir}/${archive_name}" "$download_url" || error "Download failed from ${download_url}"

    info "Extracting..."
    tar -xzf "${tmp_dir}/${archive_name}" -C "$tmp_dir"

    # Find the binary in the extracted files
    local binary_path
    binary_path=$(find "$tmp_dir" -name "$BINARY" -type f | head -1)
    [ -z "$binary_path" ] && error "Binary not found in archive"

    # Ensure install directory exists
    mkdir -p "$install_dir"

    # Install binary
    cp "$binary_path" "${install_dir}/${BINARY}"
    chmod +x "${install_dir}/${BINARY}"

    # Check if install_dir is in PATH
    if ! echo "$PATH" | tr ':' '\n' | grep -q "^${install_dir}$"; then
        warn "${install_dir} is not in your PATH."
        echo ""
        echo "Add it by running:"
        echo "  echo 'export PATH=\"${install_dir}:\$PATH\"' >> ~/.bashrc"
        echo "  source ~/.bashrc"
        echo ""
    fi
}

# ---------- main ----------
main() {
    echo ""
    echo "  mini_launch installer"
    echo "  ---------------------"
    echo ""

    # Detect platform
    local os arch
    os=$(detect_os)
    arch=$(detect_arch)
    info "Platform: ${os}/${arch}"

    # Get latest version
    local latest_version
    latest_version=$(get_latest_version)
    info "Latest version: ${latest_version}"

    # Check installed version
    local installed_version
    installed_version=$(get_installed_version)

    if [ -n "$installed_version" ]; then
        if [ "$installed_version" = "${latest_version#v}" ]; then
            info "mini_launch ${installed_version} is already installed and up to date."
            exit 0
        else
            warn "Installed version: ${installed_version}, upgrading to ${latest_version}"
        fi
    fi

    # Determine install location
    local install_dir
    install_dir=$(get_install_dir)

    # Download and install
    install "$os" "$arch" "$latest_version" "$install_dir"

    # Verify
    local installed_binary="${install_dir}/${BINARY}"
    if [ -x "$installed_binary" ]; then
        local new_version
        new_version=$("$installed_binary" --version 2>/dev/null | head -1 || echo "")
        echo ""
        info "mini_launch installed successfully!"
        info "Location: ${installed_binary}"
        [ -n "$new_version" ] && info "Version: ${new_version}"
        echo ""
        info "Quick start:"
        echo "  mkdir -p \$HOME/servers/myapp && cd \$HOME/servers/myapp"
        echo "  mini_launch initial"
        echo "  mini_launch start"
        echo ""
    else
        error "Installation failed — binary not found at ${installed_binary}"
    fi
}

main "$@"
