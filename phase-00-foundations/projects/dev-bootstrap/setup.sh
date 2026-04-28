#!/usr/bin/env bash
# setup.sh - Bootstrap a fresh Mac for SWE work
#
# Idempotent: safe to run multiple times.
# Mac (Apple Silicon and Intel) supported.

set -euo pipefail

# ----- Pretty output helpers -----
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

info()  { echo -e "${BLUE}[INFO]${NC} $*"; }
ok()    { echo -e "${GREEN}[ OK ]${NC} $*"; }
warn()  { echo -e "${YELLOW}[WARN]${NC} $*"; }
err()   { echo -e "${RED}[ERR ]${NC} $*" >&2; }

# ----- Detect platform -----
ARCH="$(uname -m)"  # arm64 on Apple Silicon, x86_64 on Intel
OS="$(uname -s)"    # Darwin on Mac

if [[ "$OS" != "Darwin" ]]; then
    err "This script is for macOS only. You're on: $OS"
    exit 1
fi

info "Detected: $OS $ARCH"

# ----- Homebrew -----
install_brew() {
    if command -v brew &>/dev/null; then
        ok "Homebrew already installed: $(brew --version | head -1)"
        return
    fi

    info "Installing Homebrew..."
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

    # Add brew to PATH for current session
    if [[ "$ARCH" == "arm64" ]]; then
        eval "$(/opt/homebrew/bin/brew shellenv)"
    else
        eval "$(/usr/local/bin/brew shellenv)"
    fi

    ok "Homebrew installed"
}

# ----- Brew packages -----
# Lists: keep them sorted, easy to scan.

CORE_PACKAGES=(
    git
    gh           # GitHub CLI
    glab         # GitLab CLI
    jq           # JSON processor
    ripgrep      # rg - fast grep
    fd           # fast find
    bat          # cat with colors
    fzf          # fuzzy finder
    tmux         # terminal multiplexer
    htop         # process viewer
    tree         # directory tree printer
    wget
    httpie       # better curl
    tldr         # short man pages
    shellcheck   # bash linter (used in CI)
)

LANGUAGES=(
    go
    rust
    node
    python@3.12
    scala
    openjdk
)

CONTAINERS=(
    kubectl
    k3d
    helm
)

DATABASES=(
    postgresql@16
    redis
)

CASKS=(
    iterm2
    docker
    visual-studio-code
    rectangle    # window manager
    raycast      # spotlight replacement
)

install_brew_packages() {
    info "Updating brew..."
    brew update

    info "Installing core packages..."
    brew install "${CORE_PACKAGES[@]}"

    info "Installing languages..."
    brew install "${LANGUAGES[@]}"

    info "Installing container tools..."
    brew install "${CONTAINERS[@]}"

    info "Installing databases (not started)..."
    brew install "${DATABASES[@]}"

    info "Installing GUI apps..."
    brew install --cask "${CASKS[@]}" || warn "Some casks may require manual install (already in /Applications?)"
}

# ----- Git config -----
configure_git() {
    info "Configuring Git globally..."

    # Ask for name/email if not already set
    local current_name
    local current_email
    current_name="$(git config --global user.name || echo '')"
    current_email="$(git config --global user.email || echo '')"

    if [[ -z "$current_name" ]]; then
        read -rp "Git name (e.g. Kvng): " name
        git config --global user.name "$name"
    else
        ok "Git name already set: $current_name"
    fi

    if [[ -z "$current_email" ]]; then
        read -rp "Git email: " email
        git config --global user.email "$email"
    else
        ok "Git email already set: $current_email"
    fi

    # Sensible defaults
    git config --global init.defaultBranch main
    git config --global pull.rebase false
    git config --global push.autoSetupRemote true
    git config --global core.editor "code --wait"
    git config --global rerere.enabled true   # remember conflict resolutions

    # If git-delta is installed, use it for diffs
    if command -v delta &>/dev/null; then
        git config --global core.pager delta
        git config --global interactive.diffFilter "delta --color-only"
        git config --global delta.navigate true
        git config --global merge.conflictstyle diff3
        git config --global diff.colorMoved default
    fi

    ok "Git configured"
}

# ----- SSH key -----
setup_ssh_key() {
    local key_path="$HOME/.ssh/id_ed25519"

    if [[ -f "$key_path" ]]; then
        ok "SSH key already exists at $key_path"
    else
        info "Generating SSH key..."
        local email
        email="$(git config --global user.email)"
        ssh-keygen -t ed25519 -C "$email" -f "$key_path" -N ""
        ok "SSH key created"
    fi

    # Start ssh-agent and add the key
    eval "$(ssh-agent -s)" >/dev/null
    ssh-add --apple-use-keychain "$key_path" 2>/dev/null || ssh-add "$key_path"

    info "Your public SSH key (add this to GitHub & GitLab):"
    echo "----------------------------------------"
    cat "${key_path}.pub"
    echo "----------------------------------------"
    pbcopy < "${key_path}.pub" && ok "Copied to clipboard."
}

# ----- Main -----
main() {
    info "Starting bootstrap..."
    install_brew
    install_brew_packages
    configure_git
    setup_ssh_key
    ok "Bootstrap complete!"
    info "Next steps:"
    echo "  1. Run ./link.sh to symlink dotfiles."
    echo "  2. Open a new terminal (or 'source ~/.zshrc')."
    echo "  3. Add your SSH key to GitHub: https://github.com/settings/keys"
    echo "  4. Add your SSH key to GitLab: https://gitlab.com/-/profile/keys"
}

main "$@"
