# ~/.zshrc - shell config
# This file is sourced every time you open a new interactive shell.

# ===== Homebrew =====
if [[ -f /opt/homebrew/bin/brew ]]; then
    eval "$(/opt/homebrew/bin/brew shellenv)"     # Apple Silicon
elif [[ -f /usr/local/bin/brew ]]; then
    eval "$(/usr/local/bin/brew shellenv)"        # Intel
fi

# ===== History =====
HISTSIZE=50000
SAVEHIST=50000
HISTFILE=~/.zsh_history
setopt HIST_IGNORE_DUPS       # don't record duplicates
setopt HIST_IGNORE_SPACE      # cmds starting with space aren't recorded (for secrets)
setopt SHARE_HISTORY          # share history between concurrent shells
setopt EXTENDED_HISTORY       # save timestamp + duration

# ===== Editor =====
export EDITOR="code --wait"
export VISUAL="$EDITOR"

# ===== Path additions =====
export GOPATH="$HOME/go"
export PATH="$PATH:$GOPATH/bin"
export PATH="$PATH:$HOME/.cargo/bin"
export PATH="$PATH:$HOME/.local/bin"

# ===== Aliases =====
# General
alias ll='ls -lah'
alias la='ls -lAh'
alias ..='cd ..'
alias ...='cd ../..'
alias c='clear'
alias h='history'
alias reload='source ~/.zshrc'

# Git
alias gs='git status'
alias gd='git diff'
alias gds='git diff --staged'
alias gco='git checkout'
alias gsw='git switch'
alias gp='git push'
alias gl='git pull'
alias gc='git commit'
alias gca='git commit --amend'
alias gap='git add -p'
alias glog='git log --oneline --graph --all -20'

# Docker
alias d='docker'
alias dc='docker compose'
alias dps='docker ps'
alias dpa='docker ps -a'

# Kubernetes
alias k='kubectl'
alias kgp='kubectl get pods'
alias kgs='kubectl get svc'
alias kgd='kubectl get deployments'
alias kctx='kubectl config use-context'

# Modern replacements
alias cat='bat --paging=never'
alias find='fd'
alias grep='rg'

# ===== Completions =====
# Initialize zsh completion system
autoload -Uz compinit && compinit

# Enable git tab-completion (brew install zsh-completions for more)

# ===== Prompt =====
# A simple, fast prompt. For fancier, install starship: brew install starship && starship init zsh
autoload -Uz vcs_info
precmd() { vcs_info }
zstyle ':vcs_info:git:*' formats ' (%b)'
setopt PROMPT_SUBST
PROMPT='%F{cyan}%~%f%F{yellow}${vcs_info_msg_0_}%f $ '

# ===== fzf =====
[ -f ~/.fzf.zsh ] && source ~/.fzf.zsh

# ===== Functions =====

# mkdir + cd in one
mkcd() {
    mkdir -p "$1" && cd "$1"
}

# Find process on port
port() {
    lsof -i ":$1"
}

# Kill process on port
killport() {
    lsof -ti ":$1" | xargs kill -9
}
