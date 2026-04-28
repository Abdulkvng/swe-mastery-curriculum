# How to push this curriculum to GitHub

I've prepared everything as a Git repo with an initial commit already made. You just need to push it from your laptop. Three steps.

## Step 1 — Create the empty repo on GitHub

Go to https://github.com/new and:

- **Repository name:** `gap_knowledge`
- **Visibility:** Public or Private — your call. (Public looks great on your resume; Private if you want to polish first.)
- **Do NOT** check "Add a README" / "Add .gitignore" / "Choose a license" — this repo already has them. GitHub initializing them creates a conflict.
- Click **Create repository**.

## Step 2 — Download this curriculum to your laptop

I delivered the curriculum as files in this conversation. Download the zip from the Claude UI (or move the folder via `present_files`), then unzip it. You'll have a folder called `swe-mastery-curriculum/`.

## Step 3 — Push it

From your laptop terminal:

```bash
cd swe-mastery-curriculum

# Option A: use the included script (recommended)
./push-to-github.sh

# Option B: do it manually
git remote add origin git@github.com:Akvng/gap_knowledge.git
git branch -M main
git push -u origin main
```

> Note: the script defaults `GITHUB_USER=Akvng` based on your `kvng.dev` brand. If your GitHub username is different, run:
>
> ```bash
> GITHUB_USER=your-actual-username ./push-to-github.sh
> ```

## What I already did for you

- `git init -b main` ✅
- Initial commit with all 75 files ✅
- `.gitignore` covering env files, node_modules, target/, certs, etc. ✅
- `LICENSE` (MIT) ✅
- `README.md` with the full curriculum overview ✅

When you push, the repo's main branch will land with one clean initial commit.

## Verifying the push worked

After pushing:

```bash
# Should show "origin" pointing at github.com
git remote -v

# Open the repo in your browser
open https://github.com/Akvng/gap_knowledge
```

You should see all 13 phase folders + README + LICENSE rendered.

## Troubleshooting

**"Permission denied (publickey)"** → your SSH key isn't on GitHub. Phase 0's `setup.sh` covers this; the short version:
```bash
cat ~/.ssh/id_ed25519.pub | pbcopy
# paste into https://github.com/settings/keys
ssh -T git@github.com   # should say "Hi <username>!"
```

**"Repository not found"** → either the repo doesn't exist on GitHub yet (do step 1), or your username is wrong (re-run with `GITHUB_USER=...`).

**"Updates were rejected"** → the GitHub repo wasn't empty (you accidentally let GitHub create a README). Easiest fix: delete the repo on GitHub, recreate empty, push again. Or `git pull --rebase origin main` then push.

---

Once it's up, your repo is your portfolio. Pin it on your GitHub profile.
