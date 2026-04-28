# Beginner Map: How This Curriculum Connects

Start here if terms like `zsh`, `bash`, `PATH`, `terminal`, `shell`, `Postgres`, or `Redis` feel random.

## The big picture

Software engineering is a stack of layers:

```text
Your laptop
  -> Operating system: macOS or Linux
  -> Terminal: the window you type in
  -> Shell: zsh or bash, the program that reads your commands
  -> Commands: cd, ls, git, python, node, docker
  -> Code: Go, Python, TypeScript, Rust
  -> Apps/APIs: servers that receive requests and return responses
  -> Databases: Postgres, Redis, SQLite, etc. store data
  -> Infrastructure: Docker, Kubernetes, CI/CD, cloud, monitoring
```

When you type:

```bash
git push origin main
```

several layers work together:

1. Your terminal displays the text window.
2. Your shell reads the command.
3. The shell finds `git` using your `PATH`.
4. Git reads your project files.
5. Git connects to GitHub over the network.
6. GitHub stores your code remotely.

That is why Phase 0 matters. It is not random command memorization. It is the control panel for the rest of software engineering.

## Core beginner definitions

### Terminal

A terminal is the app or window where you type commands.

Examples: macOS Terminal, iTerm2, VS Code integrated terminal.

The terminal mostly handles input and output. It does not truly understand commands by itself.

### Shell

A shell is the program inside the terminal that reads your commands and runs programs.

Examples: `zsh`, `bash`, `fish`.

### zsh

`zsh` means Z shell. It is the default shell on many modern Macs.

When you type:

```bash
cd projects
ls
```

`zsh` is usually the thing reading those commands.

Your zsh config file is usually:

```bash
~/.zshrc
```

That is where aliases, PATH edits, and shell settings usually go.

### bash

`bash` means Bourne Again SHell. It is another shell and is very common on Linux servers and in shell scripts.

Many scripts start with:

```bash
#!/usr/bin/env bash
```

That line tells the operating system to run the file using bash.

### zsh vs bash

They do the same core job: read commands and run programs.

| Term | What it is | Where you see it |
|---|---|---|
| `zsh` | Interactive shell on macOS | daily terminal use |
| `bash` | Common scripting shell | `.sh` scripts, Linux servers, CI jobs |

At first, do not overthink the difference. Learn the commands they share: `cd`, `ls`, `pwd`, `mkdir`, `rm`, `chmod`, variables, pipes, redirects, and scripts.

### PATH

`PATH` is a list of folders where your shell searches for programs.

When you type:

```bash
python
```

the shell checks every folder in `PATH` until it finds a program named `python`.

Check it with:

```bash
echo $PATH
```

### File path

A file path is the address of a file or folder.

Examples:

```bash
/Users/abdul/projects/app/main.go
./push-to-github.sh
../README.md
```

| Type | Example | Meaning |
|---|---|---|
| Absolute path | `/Users/abdul/projects` | full address from root `/` |
| Relative path | `./script.sh` | address from your current folder |
| Home path | `~/projects` | shortcut for your home folder |

### chmod

`chmod` means change mode. It changes file permissions.

Example:

```bash
chmod +x push-to-github.sh
```

This means: make the script executable so you can run it.

### Git vs GitHub

Git tracks versions of your code on your laptop.

GitHub hosts Git repositories online.

Git is the tool. GitHub is the website/service.

### Repository

A repository, or repo, is a project folder tracked by Git.

If a folder has a hidden `.git` folder, it is a Git repo.

Check with:

```bash
ls -la
```

## How the phases relate

- Phase 0 teaches the terminal, shell, Git, GitHub, and CI/CD.
- Phase 1 teaches how computers communicate over networks.
- Phase 2 teaches Go and backend design patterns.
- Phase 3 teaches data structures and performance thinking.
- Phase 4 teaches databases: Postgres, Redis, indexes, ACID, transactions, replication, sharding, and time-series databases.
- Phase 5 teaches APIs, which are how apps expose and modify database-backed data.
- Phase 6 teaches processes, threads, locks, and concurrency.
- Phase 7 teaches distributed systems, which explains scaling, replication, consensus, and consistency.
- Phase 8 connects the fundamentals to Datadog-style data infrastructure.
- Phase 9 combines everything into capstones.

## Database mastery check

Yes, this repo has database content.

The main database phase is:

```text
phase-04-databases/README.md
```

There is also a B-tree project:

```text
phase-04-databases/projects/btree-from-scratch/README.md
```

That project teaches the data structure behind database indexes. It is useful, but it is not enough by itself to master databases.

To truly master databases, build these in order:

1. SQL playground with users, posts, comments, likes, and audit logs.
2. Postgres-backed Task API.
3. Redis cache and rate limiter.
4. Mini database connection pool or mini pgbouncer.
5. B-tree from scratch.
6. Query tuning lab with `EXPLAIN ANALYZE`.
7. Mini time-series database for metrics.

The mental model:

```text
Postgres = durable source of truth
Redis = fast temporary helper
Indexes = speed up reads but slow down writes
Transactions = group multiple changes safely
Connection pools = reuse expensive database connections
Query planner = database brain that chooses how to run SQL
```

## Recommended order if you are shaky

1. Read this file.
2. Read `phase-00-foundations/README.md`.
3. Build `phase-00-foundations/projects/dev-bootstrap`.
4. Read `phase-04-databases/README.md` slowly.
5. Build a SQL playground before jumping into B-trees.
6. Then build the B-tree project.
7. Then build a Postgres-backed API.

Do not rush the vocabulary. Being able to explain simple terms clearly is what makes advanced systems click later.
