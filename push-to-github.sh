#!/usr/bin/env bash
# push-to-github.sh
#
# Initializes git, makes the first commit, and pushes to your GitHub.
# Run this from inside the swe-mastery-curriculum/ directory on your laptop.
#
# Prereqs:
#   1. You've created the empty repo on github.com (https://github.com/new):
#        Name: gap_knowledge
#        Visibility: your choice (private if you want, public if you want it on your resume)
#        DO NOT initialize with README/LICENSE — this script provides those
#   2. SSH key added to GitHub (you set this up in Phase 0).
#
# Usage:
#   chmod +x push-to-github.sh
#   ./push-to-github.sh

set -euo pipefail

GITHUB_USER="${GITHUB_USER:-Akvng}"          # change if your username is different
REPO_NAME="${REPO_NAME:-gap_knowledge}"

# Sanity check: are we in the curriculum dir?
if [[ ! -f "README.md" ]] || [[ ! -d "phase-00-foundations" ]]; then
    echo "ERROR: run this from the swe-mastery-curriculum directory."
    exit 1
fi

# Init git if needed
if [[ ! -d ".git" ]]; then
    echo "[init] git init..."
    git init -b main
fi

# Stage everything
echo "[add] staging files..."
git add -A

# First commit (or amend if the user wants)
if git rev-parse --verify HEAD &>/dev/null; then
    echo "[skip] commit already exists; using current HEAD."
else
    echo "[commit] creating initial commit..."
    git commit -m "Initial commit: SWE Mastery Curriculum (phases 0–12)"
fi

# Set remote
REMOTE_URL="git@github.com:${GITHUB_USER}/${REPO_NAME}.git"
if git remote get-url origin &>/dev/null; then
    echo "[remote] origin already set; updating URL to $REMOTE_URL"
    git remote set-url origin "$REMOTE_URL"
else
    echo "[remote] adding origin → $REMOTE_URL"
    git remote add origin "$REMOTE_URL"
fi

# Push
echo "[push] pushing to GitHub..."
git push -u origin main

echo
echo "Done. Repo at: https://github.com/${GITHUB_USER}/${REPO_NAME}"
