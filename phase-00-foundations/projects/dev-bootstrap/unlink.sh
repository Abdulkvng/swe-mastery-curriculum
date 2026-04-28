#!/usr/bin/env bash
# unlink.sh - Remove dotfile symlinks (useful for testing / teardown)

set -euo pipefail

unlink_file() {
    local target="$1"
    if [[ -L "$target" ]]; then
        rm "$target"
        echo "[ OK ] Removed symlink: $target"
    else
        echo "[INFO] Not a symlink, skipping: $target"
    fi
}

unlink_file "$HOME/.zshrc"
unlink_file "$HOME/.gitconfig"
unlink_file "$HOME/.gitignore_global"

echo "Done. Restore from .backup files if needed."
