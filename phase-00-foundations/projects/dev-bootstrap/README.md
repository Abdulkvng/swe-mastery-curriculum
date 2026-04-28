# dev-bootstrap

> Set up a fresh Mac for SWE work in 15 minutes. Includes dotfiles and CI.

## What this is

A bash project that:
1. Installs Homebrew + a vetted set of dev tools.
2. Configures Git globally.
3. Generates an SSH key and prints the public key for you to add to GitHub/GitLab.
4. Symlinks dotfiles (`.zshrc`, `.gitconfig`) so they're version-controlled.
5. Has a GitHub Actions workflow that lints all bash scripts via `shellcheck`.

## Files in this project

```
dev-bootstrap/
├── README.md              <- this
├── setup.sh               <- main installer
├── link.sh                <- symlinks dotfiles into $HOME
├── unlink.sh              <- removes the symlinks (for testing)
├── dotfiles/
│   ├── .zshrc
│   ├── .gitconfig
│   └── .gitignore_global
└── .github/
    └── workflows/
        └── ci.yml         <- shellcheck on every push
```

## Run it

```bash
git clone git@github.com:<you>/dev-bootstrap.git
cd dev-bootstrap
./setup.sh
./link.sh
source ~/.zshrc
```

## What you should learn from this project

- `set -euo pipefail` and why every script needs it
- Idempotency: running the script twice should be safe
- Symlinks: `ln -sf` and why we use them for dotfiles
- Detecting Apple Silicon vs Intel and branching accordingly
- Reading a real CI YAML file and understanding what each line does
