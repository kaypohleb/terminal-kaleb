# terminal-kaleb


SSH server ([Wish](https://github.com/charmbracelet/wish)) that serves a [Bubble Tea](https://github.com/charmbracelet/bubbletea) TUI per session.

## Prerequisite

Host key for the server (run from repo root):

```bash
mkdir -p .ssh
ssh-keygen -t ed25519 -f .ssh/id_ed25519 -N ""
```

## Run

```bash
go run .
```

Optional: `SSH_HOST` (default `localhost`), `SSH_PORT` (default `23234`).

## Connect

```bash
ssh -p 23234 localhost
```

Press `q` to disconnect.

For local dev, you can use [Wish’s `~/.ssh/config` tip](https://github.com/charmbracelet/wish#pro-tip) to avoid `known_hosts` churn.


Choose your own adventure
##page - option destination - action