#!/usr/bin/env bash
# link.sh - Symlink dotfiles from this repo into $HOME

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DOTFILES_DIR="$SCRIPT_DIR/dotfiles"

# What does symlink mean?
# A symlink is a file that points to another file. When something reads the symlink,
# it actually reads the target file. Useful for keeping dotfiles in a Git repo
# while having the OS look for them at $HOME.

link_file() {
    local src="$1"   # source file (in this repo)
    local dst="$2"   # destination (where the OS expects it)

    if [[ -L "$dst" ]]; then
        echo "[ OK ] Already symlinked: $dst"
        return
    fi

    if [[ -e "$dst" ]]; then
        local backup="${dst}.backup.$(date +%s)"
        echo "[INFO] Backing up existing $dst -> $backup"
        mv "$dst" "$backup"
    fi

    ln -s "$src" "$dst"
    echo "[ OK ] Linked: $dst -> $src"
}

link_file "$DOTFILES_DIR/.zshrc"            "$HOME/.zshrc"
link_file "$DOTFILES_DIR/.gitconfig"        "$HOME/.gitconfig"
link_file "$DOTFILES_DIR/.gitignore_global" "$HOME/.gitignore_global"

# Tell git about the global gitignore
git config --global core.excludesfile "$HOME/.gitignore_global"

echo
echo "Done. Run 'source ~/.zshrc' or open a new terminal."
