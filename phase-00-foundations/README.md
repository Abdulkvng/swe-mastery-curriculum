# Phase 0 — Foundations & Environment

> Before you can build anything, your tools need to be sharp. This phase gets your Mac set up the way real engineers' machines are set up, teaches you Linux/bash properly (not just memorized commands), gets Git into your bones, and gives you a working CI/CD pipeline you'll reuse for every later phase.

**Time:** 2–4 days if you're new to most of this; 1 day if you're just filling gaps.

**You'll know you're done when:** you can clone a repo, branch, edit, commit, push, open a PR, and watch it run through CI — all without thinking about it.

---

## Table of contents

1. [What does this even mean? — "Foundations"](#what-foundations-means)
2. [Module 0.1 — Mac setup for SWE work](#module-01--mac-setup)
3. [Module 0.2 — The shell, deeply](#module-02--the-shell-deeply)
4. [Module 0.3 — Linux/Unix philosophy & core tools](#module-03--linuxunix-philosophy)
5. [Module 0.4 — Git, internals included](#module-04--git-internals-included)
6. [Module 0.5 — GitHub vs GitLab workflows](#module-05--github-vs-gitlab)
7. [Module 0.6 — CI/CD: what it is, what it isn't](#module-06--cicd)
8. [🛠️ Project: `dev-bootstrap`](#project-dev-bootstrap)
9. [Exercises](#exercises)
10. [What you should now know](#what-you-should-now-know)

---

<a name="what-foundations-means"></a>
## 🧠 What does this even mean? — "Foundations"

When engineers say "foundations," they don't mean "the basics you can skip if you already know them." They mean **the layer beneath everything else** — the stuff that, if you don't have it, makes every later thing 3x harder.

A jazz pianist who doesn't know scales can still play. They just play badly, and they can't ever play *fast*. Same with engineers: you can ship code without knowing what `chmod +x` does — until the day a script won't run and you have no idea why, and you waste 40 minutes Googling.

This phase isn't about memorization. It's about **building the muscle memory** that lets you do the boring stuff without thinking, so your brain is free for the hard stuff.

---

<a name="module-01--mac-setup"></a>
## Module 0.1 — Mac setup for SWE work

### Why this matters

You have an M3 Pro. Apple Silicon (ARM64) is great for performance and battery, but some tools need attention to install correctly. Setting up your Mac *once*, properly, saves hundreds of "why isn't this working" moments later.

### Step 1: Install Homebrew

> 📖 **Definition — Homebrew:** The de facto package manager for macOS. A "package manager" is a program that installs other programs. Instead of manually downloading apps, you run `brew install <thing>`.

```bash
# Run this in Terminal (the built-in app at /Applications/Utilities/Terminal.app)
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

After installing, Homebrew will tell you to add it to your PATH. On Apple Silicon, the command is usually:

```bash
echo 'eval "$(/opt/homebrew/bin/brew shellenv)"' >> ~/.zprofile
eval "$(/opt/homebrew/bin/brew shellenv)"
```

> 📖 **Definition — PATH:** An environment variable (a system-wide setting) that lists directories the shell searches when you type a command. When you type `git`, the shell looks in each PATH directory in order until it finds an executable named `git`.

Verify:

```bash
brew --version
# Homebrew 4.x.x
```

### Step 2: Install the core SWE toolkit

```bash
# Languages we'll use
brew install go            # Go - for backend, infra, k8s ecosystem
brew install rust          # Rust - for systems, networking, performance
brew install node          # Node.js - for TypeScript/React
brew install python@3.12   # Python - for ML, scripting
brew install scala         # Scala - for Spark
brew install openjdk       # Java JDK - Scala/Spark need this

# Build tools
brew install cmake make pkg-config

# Containers & orchestration
brew install docker         # Container runtime CLI
brew install --cask docker  # Docker Desktop (the GUI + daemon)
brew install kubectl        # Kubernetes CLI
brew install k3d            # Lightweight k8s in Docker (we'll use this for capstone)
brew install helm           # k8s package manager

# Databases
brew install postgresql@16
brew install redis

# Networking & debugging
brew install wireshark      # Packet inspection (we'll use it in Phase 1)
brew install nmap           # Network scanner
brew install httpie         # Friendlier curl
brew install jq             # JSON processor (used everywhere)

# Git tooling
brew install gh             # GitHub CLI
brew install glab           # GitLab CLI
brew install git-delta      # Pretty git diffs

# Quality-of-life
brew install ripgrep        # Faster grep (rg)
brew install fd             # Faster find
brew install bat            # Cat with syntax highlighting
brew install fzf            # Fuzzy finder
brew install tmux           # Terminal multiplexer
brew install neovim         # Or stick with VS Code
```

Each of these will reappear in later phases. You're not learning them all now — you're just installing them so they're ready.

### Step 3: Install a real terminal and shell

macOS Terminal works, but iTerm2 + zsh + Oh-My-Zsh is the standard real-engineer stack:

```bash
brew install --cask iterm2
# zsh is already the default on modern macOS
# install Oh-My-Zsh:
sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"
```

> 📖 **Definition — Terminal vs Shell:** A *terminal* is the window (Terminal.app, iTerm2). A *shell* is the program running inside it that interprets your commands (zsh, bash). They're different things that get conflated.

### Step 4: Set up your `.zshrc`

Your `.zshrc` is the file the shell reads every time you open a new terminal. It's where you put aliases, functions, and PATH tweaks.

Open it: `code ~/.zshrc` (or `nvim ~/.zshrc`).

Add these helpful lines to the end:

```bash
# Custom aliases - shortcuts for commands you type often
alias ll='ls -lah'
alias gs='git status'
alias gd='git diff'
alias gco='git checkout'
alias gp='git push'
alias gl='git pull'
alias k='kubectl'
alias d='docker'
alias dc='docker compose'

# Better history
HISTSIZE=50000
SAVEHIST=50000
setopt HIST_IGNORE_DUPS
setopt SHARE_HISTORY

# Make less work nicely with colors
export LESS='-R'

# Go path
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin

# Rust path (cargo installs here)
export PATH=$PATH:$HOME/.cargo/bin
```

Reload it:

```bash
source ~/.zshrc
```

> 📖 **Definition — `source`:** A shell command that runs another file *in the current shell*, so any variables/functions defined there are available now. (Just running `./file.sh` runs it in a subshell — changes don't persist.)

### Step 5: SSH keys for GitHub/GitLab

> 📖 **Definition — SSH key:** A pair of cryptographic keys (public + private) used to prove your identity to a remote server without typing a password. The public key sits on the server; the private key stays on your laptop. You'll learn the math behind this in Phase 1.

```bash
# Generate a key (use Ed25519, the modern algorithm)
ssh-keygen -t ed25519 -C "your_email@example.com"
# Press Enter to accept the default location (~/.ssh/id_ed25519)
# Set a passphrase (optional but recommended)

# Add to ssh-agent so you don't retype passphrase
eval "$(ssh-agent -s)"
ssh-add --apple-use-keychain ~/.ssh/id_ed25519

# Copy your public key to clipboard
pbcopy < ~/.ssh/id_ed25519.pub
```

Now go to GitHub → Settings → SSH and GPG keys → New SSH key → paste. Same for GitLab.

Test:

```bash
ssh -T git@github.com
# "Hi <username>! You've successfully authenticated..."
```

---

<a name="module-02--the-shell-deeply"></a>
## Module 0.2 — The shell, deeply

Most tutorials show you commands. This section explains *what's actually happening* when you run them.

### What is a command, really?

When you type `ls`, here's what happens:

```
1. You press Enter.
2. The shell parses the line: command="ls", args=[]
3. The shell searches each directory in $PATH for an executable file named "ls"
4. It finds /bin/ls
5. The shell calls fork() to create a child process
6. The child calls execve() to replace itself with /bin/ls
7. /bin/ls runs, reads the current directory, prints output
8. /bin/ls exits with a status code (0 = success, non-zero = error)
9. The shell prints the prompt again
```

That `fork()` + `execve()` dance is how *every* program runs on Unix-like systems. We'll revisit this in Phase 6 (Concurrency & OS).

> 📖 **Definition — Process:** A running program, with its own memory and identity (a Process ID, or PID). When `ls` runs, it's a process. When it finishes, the process dies.

### The 20 commands you'll use 90% of the time

```bash
# Navigation
pwd                    # print working directory
cd ~/code              # change to directory (~ = home)
cd -                   # go back to previous directory
ls -lah                # list, long format, all (incl hidden), human sizes

# Files
touch file.txt         # create empty file (or update timestamp)
mkdir -p a/b/c         # make directory + parents as needed
cp src dst             # copy
cp -r srcdir dstdir    # copy directory recursively
mv old new             # move/rename
rm file                # delete file
rm -rf dir             # delete directory recursively (DANGER: no undo)

# Reading
cat file               # dump file to stdout
less file              # paginated view (q to quit, / to search)
head -n 20 file        # first 20 lines
tail -n 20 file        # last 20 lines
tail -f log.txt        # follow a file as it grows (great for logs)

# Searching
grep "pattern" file    # find lines matching pattern
grep -r "pattern" .    # recursive search
rg "pattern"           # ripgrep - much faster, respects .gitignore
find . -name "*.go"    # find files by name
fd "\.go$"             # fd - faster, simpler

# Permissions
chmod +x script.sh     # make executable
chmod 644 file         # rw-r--r-- (owner: rw, others: r)
chown user:group file  # change owner

# Process management
ps aux                 # all running processes
top                    # interactive process viewer
htop                   # nicer top (brew install htop)
kill <pid>             # send TERM signal (graceful)
kill -9 <pid>          # send KILL signal (force)
```

### Pipes, redirects, and the Unix philosophy

> 📖 **Definition — Pipe (`|`):** Connects the stdout of one command to the stdin of another. `cmd1 | cmd2` means "feed cmd1's output as cmd2's input."

> 📖 **Definition — stdin/stdout/stderr:** Every process has three default streams: standard input (where it reads from), standard output (where it prints), standard error (where it prints errors). By default they're all connected to your terminal.

The Unix philosophy (Doug McIlroy, 1978):
> "Write programs that do one thing and do it well. Write programs to work together."

This is why pipes are magical:

```bash
# Find all .go files, count their lines
find . -name "*.go" | xargs wc -l | tail -1

# What's eating my disk?
du -sh ./* | sort -hr | head -10

# Show me unique IPs hitting my nginx log, top 10
cat access.log | awk '{print $1}' | sort | uniq -c | sort -rn | head -10

# Check if a port is open
nc -zv localhost 8080
```

Each of those is a small, single-purpose command. Composed, they're a search engine, a disk analyzer, a log aggregator, a network diagnostic.

Redirects:

```bash
cmd > file            # write stdout to file (overwrite)
cmd >> file           # write stdout to file (append)
cmd 2> errors.log     # write stderr to file
cmd > out.log 2>&1    # both stdout AND stderr to one file
cmd < input.txt       # feed file as stdin
cmd1 && cmd2          # run cmd2 only if cmd1 succeeds
cmd1 || cmd2          # run cmd2 only if cmd1 fails
cmd1 ; cmd2           # run cmd2 unconditionally after cmd1
```

### Shell scripting basics

A shell script is just a file of commands the shell runs in order. Save as `script.sh`:

```bash
#!/usr/bin/env bash
# ^ "shebang" - tells the OS which interpreter to use

set -euo pipefail
# -e : exit immediately if a command fails
# -u : exit if you reference an undefined variable
# -o pipefail : a pipe fails if ANY command in it fails (not just the last)
# THIS LINE IS GOLDEN. Put it in every script.

# Variables
NAME="Kvng"
COUNT=5

# String interpolation - use double quotes, ${VAR}
echo "Hello, ${NAME}!"

# Arithmetic
SUM=$((COUNT + 10))
echo "Sum: $SUM"

# If/else
if [[ -f "myfile.txt" ]]; then
    echo "File exists"
elif [[ -d "myfile.txt" ]]; then
    echo "It's a directory"
else
    echo "Nope"
fi

# For loop
for i in 1 2 3; do
    echo "Iteration $i"
done

# For loop over files
for f in *.go; do
    echo "Found: $f"
done

# Functions
greet() {
    local name="$1"  # 'local' = scoped to this function
    echo "Hello, $name"
}

greet "Kvng"

# Capture command output
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
echo "On branch: $CURRENT_BRANCH"
```

Make executable and run:

```bash
chmod +x script.sh
./script.sh
```

> 🎯 **Interview note:** Apple sometimes asks "write a bash one-liner to find duplicate files" or "what does `set -e` do?" Knowing the basics here is table stakes.

---

<a name="module-03--linuxunix-philosophy"></a>
## Module 0.3 — Linux/Unix philosophy & core tools

### Why this matters even though you have a Mac

macOS is BSD-derived Unix, so 95% of Linux skills transfer. But Datadog's servers, every Docker container you'll ever run, and Apple's data centers all run Linux. Knowing the Linux model isn't optional.

### The filesystem hierarchy

```
/                  # root of everything
├── bin/           # essential user binaries (ls, cat, ...)
├── sbin/          # system binaries (require root)
├── usr/           # user programs
│   ├── bin/       # most installed programs
│   └── local/     # locally installed (your stuff)
├── etc/           # config files
├── home/          # user home directories on Linux
│   └── kvng/
├── var/           # variable data (logs, databases)
│   └── log/
├── tmp/           # temporary files (often wiped on reboot)
├── proc/          # virtual filesystem of running processes (Linux)
└── dev/           # device files (disks, terminals as files)
```

On Mac, `/home` is `/Users` instead, but the rest is similar.

### "Everything is a file"

A core Unix idea: disks are files (`/dev/sda`), terminals are files (`/dev/tty`), even the kernel exposes processes as files (`/proc/1234/status`). This means tools like `cat`, `grep`, `dd` work on *anything*.

```bash
# Read CPU info on Linux (won't work on Mac)
cat /proc/cpuinfo

# On Mac, equivalent:
sysctl -a | grep machdep.cpu

# Read a file char-by-char from a serial device (embedded systems work)
cat /dev/ttyUSB0
```

This idea — **uniform interfaces** — is why Unix is everywhere 50 years later.

### Permissions, properly

```bash
ls -l
# -rw-r--r--  1 kvng staff  1234 Apr 27 10:00 file.txt
# │└┬┘└┬┘└┬┘  │  │   │
# │ │  │  │   │  │   └── group
# │ │  │  │   │  └────── owner
# │ │  │  │   └──── number of hard links
# │ │  │  └─────── others' permissions (read only)
# │ │  └────────── group's permissions (read only)
# │ └───────────── owner's permissions (read+write)
# └─────────────── file type (- = file, d = dir, l = symlink)
```

Numeric mode is base-8: `r=4, w=2, x=1`. So `chmod 755 file` = `rwxr-xr-x` (owner: all, group: read+exec, others: read+exec).

### Environment variables

> 📖 **Definition — Environment variable:** A key-value pair that the shell (and any program it launches) can read. Used for configuration, secrets, paths, etc.

```bash
# View all env vars
env

# Set one for this session
export API_KEY="secret123"

# Use it in a program
echo $API_KEY

# Set permanently: add to ~/.zshrc
```

### Process basics

```bash
# Run something in the background
sleep 100 &
# [1] 12345  (job number, PID)

# List background jobs in this shell
jobs

# Bring most recent back to foreground
fg

# See ALL processes
ps aux | head

# Kill a process
kill 12345        # send SIGTERM (polite "please exit")
kill -9 12345     # send SIGKILL (forced; process can't clean up)

# What's listening on port 8080?
lsof -i :8080
```

We'll go deep on processes vs threads in Phase 6.

---

<a name="module-04--git-internals-included"></a>
## Module 0.4 — Git, internals included

Most engineers use Git like a magic black box. They memorize 4 commands and panic when anything weird happens. Let's not be that.

### The mental model

> 📖 **Definition — Git:** A *content-addressable filesystem* with a version control UI on top. That sentence sounds scary. It just means: every file's content gets hashed (turned into a unique fingerprint), and Git stores files by their hash, not their name.

When you `git add` a file, Git:
1. Computes the SHA-1 hash of the file's contents.
2. Stores the file in `.git/objects/<first-2-chars>/<rest>` as a *blob*.
3. Adds an entry to the *index* (staging area): "this filename → this blob hash."

When you `git commit`:
1. Git creates a *tree* object (folder snapshot) referencing all the blobs.
2. Git creates a *commit* object referencing the tree, the parent commit, author, message.
3. The commit gets a hash too. Your branch (`main`) is just a pointer to the latest commit hash.

That's it. **Git is a graph of commits, where each commit points to a snapshot tree, where each tree points to blobs.** Branches are sticky notes on commits.

```
        c1 ─── c2 ─── c3 ◄── main
                       \
                        c4 ─── c5 ◄── feature-x
```

When you `git merge`, Git either fast-forwards `main` to `c5`, or creates a merge commit with two parents.

### The commands you'll actually use

```bash
# === Setup ===
git config --global user.name "Kvng"
git config --global user.email "you@example.com"
git config --global init.defaultBranch main
git config --global pull.rebase false  # default to merge on pull
git config --global core.editor "code --wait"  # use VSCode for commit messages

# === Daily flow ===
git status                   # what's changed?
git add file.go              # stage a file
git add .                    # stage everything
git add -p                   # stage interactively, hunk by hunk (LEARN THIS)
git commit -m "Add foo"      # commit with message
git commit                   # commit and open editor for longer message

git log --oneline            # recent commits, terse
git log --graph --oneline --all  # visualize all branches
git diff                     # what's changed but not staged
git diff --staged            # what's staged but not committed

# === Branching ===
git branch                   # list local branches
git branch -a                # all branches incl remote
git checkout -b feature-x    # create + switch to new branch
git switch feature-x         # modern alternative to checkout for branches
git switch main
git merge feature-x          # merge feature-x into current branch

# === Remotes ===
git remote -v                # show remotes
git push origin feature-x    # push branch to origin
git pull                     # fetch + merge
git fetch                    # download remote changes WITHOUT merging
git push --force-with-lease  # safer force-push (won't overwrite others' work)

# === Undo ===
git restore file.go          # discard unstaged changes to file
git restore --staged file.go # unstage but keep changes
git reset HEAD~1             # undo last commit, keep changes staged
git reset --hard HEAD~1      # undo last commit, DISCARD changes (DANGER)
git revert <hash>            # create a NEW commit that undoes <hash> (safe for shared branches)

# === Inspecting ===
git show <hash>              # see a specific commit
git blame file.go            # who wrote each line and when
git log -p file.go           # all changes to a file ever
git log --author=kvng        # commits by author

# === Stashing ===
git stash                    # park uncommitted changes
git stash pop                # bring them back
git stash list

# === Rewriting history (use carefully) ===
git commit --amend           # edit the last commit
git rebase -i HEAD~5         # interactive rebase, last 5 commits (squash, reorder, edit)
```

### `git rebase` vs `git merge` — the question

When you've branched off main and main has moved on:

```
        c1 ─── c2 ─── c3 ◄── main
                \
                 c4 ─── c5 ◄── feature-x
```

**Merge** creates a merge commit:
```
        c1 ─── c2 ─── c3 ─── M ◄── main (after merging feature-x)
                \           /
                 c4 ─── c5
```
- ✅ Preserves history exactly as it happened
- ❌ History gets cluttered with merge commits

**Rebase** rewrites your branch's commits onto main's tip:
```
        c1 ─── c2 ─── c3 ─── c4' ─── c5' ◄── feature-x
                                        ◄── main
```
- ✅ Linear, clean history
- ❌ Rewrites commit hashes — never rebase commits others have pulled

**Rule of thumb:** rebase your local feature branch *before* opening a PR. Once a PR is shared, only merge.

### `.gitignore`

Tell Git what NOT to track:

```gitignore
# Dependencies
node_modules/
vendor/
target/

# Build artifacts
*.exe
*.o
dist/
build/

# Secrets
.env
*.pem
*.key

# Editor cruft
.vscode/
.idea/
*.swp
.DS_Store

# Logs
*.log
```

> 🎯 **Interview note:** Apple commonly asks "describe your Git workflow" and "what's the difference between merge and rebase?" Be ready.

### `git add -p` — the move that separates you from juniors

Instead of `git add .` (stage everything), use `git add -p` (patch mode). Git walks you through each *hunk* of changes and asks "stage this? y/n/s/e/?" Letters:
- `y` — yes
- `n` — no
- `s` — split this hunk into smaller ones
- `e` — edit manually

This forces you to *review* your changes before committing. You'll catch print statements, debug code, half-baked thoughts. Use it.

---

<a name="module-05--github-vs-gitlab"></a>
## Module 0.5 — GitHub vs GitLab workflows

### What's the difference?

Both are platforms that host Git repos. Both have:
- Issues (bug tracker)
- Pull/Merge requests (code review)
- CI/CD pipelines
- Wikis
- Releases

GitHub:
- **Bigger ecosystem** (most open source lives here)
- "Pull requests" (PRs)
- Actions for CI/CD (YAML in `.github/workflows/`)
- Apple uses GitHub Enterprise

GitLab:
- **All-in-one DevOps** (built-in container registry, K8s integration, more enterprise features)
- "Merge requests" (MRs) — same idea, different name
- Pipelines defined in `.gitlab-ci.yml`
- **Datadog uses GitLab**

You should be fluent in both.

### Standard PR/MR workflow

```bash
# 1. Sync with main
git switch main
git pull

# 2. Branch off
git switch -c feature/add-auth

# 3. Code, commit (small, atomic commits)
git add -p
git commit -m "Add JWT validation middleware"
# more commits...

# 4. Push
git push -u origin feature/add-auth

# 5. Open PR/MR via web UI or CLI
gh pr create --title "Add JWT auth" --body "..."
# or
glab mr create

# 6. Code review happens. Address feedback:
git add -p
git commit --amend  # if it's a tiny fix to the last commit
# OR
git commit -m "Address review: handle expired tokens"
git push  # (may need --force-with-lease if you amended)

# 7. Once approved, merge (squash, rebase-merge, or merge commit — team chooses)

# 8. Clean up
git switch main
git pull
git branch -d feature/add-auth
```

### Code review etiquette (you'll do this daily at Datadog)

When reviewing:
- **Be kind. Be specific. Be helpful.** "Why this approach over X?" beats "this is wrong."
- Distinguish blocking from non-blocking. Prefix non-blocking with "nit:" or "(non-blocking):".
- Approve when you'd ship it; request changes when you wouldn't.

When receiving review:
- Don't take it personally. Code is not you.
- If you disagree, push back politely with reasoning.
- Resolve threads as you address them.

---

<a name="module-06--cicd"></a>
## Module 0.6 — CI/CD: what it is, what it isn't

### 🧠 What does CI/CD even mean?

**Continuous Integration (CI):** Every time anyone pushes code, an automated system pulls that code, builds it, runs the tests, and reports back. Goal: catch bugs in minutes, not weeks.

**Continuous Delivery (CD):** Once tests pass, the artifact is automatically *ready to deploy* — sitting in a registry waiting for someone to click "go."

**Continuous Deployment (also CD):** Goes further — auto-deploys to production. Datadog and Apple both deploy multiple times per day per service.

The "system" running all this is a **CI runner**: a server (or pool of servers) that executes pipeline jobs in containers.

### A pipeline, conceptually

```
┌────────┐   ┌──────┐   ┌──────┐   ┌────────┐   ┌────────┐
│  push  │ → │ lint │ → │ test │ → │  build │ → │ deploy │
└────────┘   └──────┘   └──────┘   └────────┘   └────────┘
                ↓ fail     ↓ fail
            block PR   block PR
```

Each stage runs in a fresh container. Stages can run in parallel (e.g., test on Linux + macOS simultaneously).

### GitHub Actions example (`.github/workflows/ci.yml`)

```yaml
name: CI

on:
  push:
    branches: [main]
  pull_request:

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4         # pull the code
      - uses: actions/setup-go@v5         # install Go
        with:
          go-version: '1.22'
      - name: Run tests
        run: go test ./...
      - name: Lint
        run: |
          go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
          golangci-lint run

  build-docker:
    needs: test                            # only runs if test job passes
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Build image
        run: docker build -t myapp:${{ github.sha }} .
```

### GitLab CI example (`.gitlab-ci.yml`) — what you'll see at Datadog

```yaml
stages:
  - test
  - build
  - deploy

variables:
  GO_VERSION: "1.22"

test:
  stage: test
  image: golang:${GO_VERSION}
  script:
    - go test ./...
    - go vet ./...

lint:
  stage: test
  image: golangci/golangci-lint:latest
  script:
    - golangci-lint run

build-image:
  stage: build
  image: docker:24
  services:
    - docker:24-dind   # docker-in-docker, lets us run docker inside a docker job
  script:
    - docker build -t $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA .
    - docker push $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA
  only:
    - main

deploy-staging:
  stage: deploy
  image: bitnami/kubectl
  script:
    - kubectl set image deployment/myapp myapp=$CI_REGISTRY_IMAGE:$CI_COMMIT_SHA
  environment: staging
  only:
    - main
```

Key concepts both share:
- **Jobs** run in isolated containers
- **Stages** group jobs and define order
- **Artifacts** (build outputs) can pass between jobs
- **Variables** can be set in the YAML or in the project's settings (for secrets)
- **Caching** speeds up repeat builds

### What "good CI" looks like

1. **Fast.** Under 5 minutes, ideally under 2. Slow CI = developers stop running it locally and just push to find out.
2. **Reliable.** No flaky tests. Flaky CI is worse than no CI because it teaches people to ignore failures.
3. **Strict in CI, loose locally.** Lint and format errors block merge but can be auto-fixed locally.
4. **Tests run on every PR.** No exceptions.
5. **Deploys are automated.** Manual deploy steps are where bugs hide.

---

<a name="project-dev-bootstrap"></a>
## 🛠️ Project: `dev-bootstrap`

A bash script + a Git repo that, when run on a fresh Mac, sets up your entire dev environment AND has a working CI pipeline.

### What you'll build

A repo (`dev-bootstrap`) containing:

1. `setup.sh` — installs Homebrew (if missing), all your toolchains, configures Git, generates SSH keys.
2. `dotfiles/` — your `.zshrc`, `.gitconfig`, `.gitignore_global`.
3. `link.sh` — symlinks the dotfiles into `$HOME`.
4. `.github/workflows/ci.yml` — runs `shellcheck` (linter for bash) on every push.
5. `README.md` — usage instructions.

### Step-by-step build

See `projects/dev-bootstrap/` in this folder. The full code is in there — you should:
1. Read the code.
2. Type it out yourself (don't copy-paste — fingers learn).
3. Run it.
4. Push to a new GitHub repo.
5. Watch CI run.
6. Break it on purpose (introduce a bash error) and watch CI fail.
7. Fix it, push again, watch CI pass.

By the end you should be able to:
- Wipe and re-set-up your machine in 15 minutes
- Onboard a friend the same way
- Read any GitLab CI / GitHub Actions YAML and understand it

---

<a name="exercises"></a>
## Exercises

These build real muscle memory. Do all of them.

1. **Bash one-liner challenges:**
   - Find the 5 largest files in your `~/Downloads`.
   - Count how many `.go` files exist in any directory below `$HOME` (use `find` then `fd`, compare speed).
   - Print the unique IPs from a hypothetical `access.log` file (sort + uniq + awk).
   - Recursively rename all `.txt` files to `.md` in a directory.
   - Write a one-liner that, given a directory, prints each subdirectory's total size sorted descending.

2. **Git scenarios:**
   - Make 3 commits, then squash them into 1 with `git rebase -i`.
   - Make a commit with a typo in the message, fix it with `--amend`.
   - Accidentally commit a secret to a file. Remove it from history (Google: `git filter-repo`).
   - Have a merge conflict on purpose: branch A and branch B both edit the same line. Merge B into A, resolve.
   - Find which commit introduced a specific line of code using `git blame` and `git log -S`.

3. **CI exercises:**
   - Write a GitHub Actions workflow that runs on PRs, prints the PR title, and posts a comment on the PR with the output.
   - Write a GitLab CI pipeline with 3 stages where the third stage runs only on the `main` branch.
   - Add a workflow that fails if any committed file has the string `TODO` in it.

---

<a name="what-you-should-now-know"></a>
## ✅ What you should now know

After this phase, you can:

- [ ] Set up a fresh Mac for SWE work in under 30 minutes
- [ ] Read and write bash scripts confidently (with `set -euo pipefail`)
- [ ] Use pipes, redirects, and the standard core tools fluently
- [ ] Explain what happens when you press Enter on a command (fork/exec)
- [ ] Explain Git's data model (blobs, trees, commits, refs)
- [ ] Choose merge vs rebase appropriately
- [ ] Use `git add -p` for clean, atomic commits
- [ ] Open and review PRs/MRs on both GitHub and GitLab
- [ ] Write a CI pipeline that lints, tests, builds, and deploys
- [ ] Diagnose a failing CI pipeline

If any of those feel shaky, redo that module before moving on. **Phase 1 builds directly on this** — you'll write a Rust networking project with Git workflow + CI from day one.

---

**Next:** [Phase 1 — Networking & Protocols](../phase-01-networking/README.md)
