<div align="center">

```
       ✦   .    .         .        ✦
   .        .        .        .         .
         .        .      .          .
     ☁️        ✷      CLAUDE NINJA      ✷        ☁️
         .        .      .          .
   .        .        .        .         .
       ✦      .        sky above, code below     ✦
```

# Claude Ninja

Project-first, manager-driven workspace for running multiple AI coding agents safely and fast.

[![GitHub Stars](https://img.shields.io/github/stars/pratham-bhatnagar/claude-ninja?style=for-the-badge&logo=github&color=yellow&labelColor=1a1b26)](https://github.com/pratham-bhatnagar/claude-ninja/stargazers)
[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=for-the-badge&logo=go&labelColor=1a1b26)](https://go.dev)
[![License](https://img.shields.io/badge/License-MIT-9ece6a?style=for-the-badge&labelColor=1a1b26)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-macOS%20%7C%20Linux%20%7C%20WSL-7aa2f7?style=for-the-badge&labelColor=1a1b26)](#)

[Features](#features) · [Installation](#installation) · [Quick Start](#quick-start) · [Usage](#usage) · [Contributing](#contributing) · [Roadmap](#roadmap) · [License](#license)

</div>

## Overview

Claude Ninja keeps every agent isolated in its own branch/worktree, routes all sub-agent output through a single Manager, and provides a focused TUI to plan, execute, and verify work across projects.

- Project-first TUI: switch projects with Left/Right; browse sub-agents with Up/Down
- Manager-only interaction: reply once in the Manager; aggregate sub-agent output there
- Worktree isolation: each agent works on its own branch/worktree for safe concurrency
- Structured planning: GSD-style plan/execute/verify via lightweight `.planning/` state
- Nudging loop: surface waiting agents to the Manager to keep flow unblocked

## Features

- Multi-project navigation and global search
- Manager panel with status, waiting agents, and planning commands
- Session forking and context inheritance
- Git worktree helpers for isolated branches
- Optional MCP integrations and pooling (where available)

## Installation

Prerequisites: Go 1.24+, git, tmux.

From source:

```bash
make build            # Build to ./build
make install          # Install to /usr/local/bin (sudo)
# or
make install-user     # Install to ~/.local/bin (no sudo)
```

Note: The current binary name in Makefile is `agent-deck` during migration. After install, run `agent-deck` to launch. The CLI will be renamed to `claude-ninja` in an upcoming update.

## Quick Start

```bash
# Launch the TUI
agent-deck

# Create/import sessions per project, then mark your Manager with "!"
# Send a sub-agent’s output to the Manager with "x"
```

## Usage

Core workflow:

- Create or import sessions per project
- Mark a Manager session: select and press `!`
- Switch projects with Left/Right; navigate sessions with Up/Down
- Press `x` on a sub-agent to send its output to the Manager
- Use the Manager to plan, execute, verify, and merge

GSD-style planning commands (run from the Manager):

- `/gsd:new-project` initialize goals and context
- `/gsd:plan-phase 1` draft a structured plan
- `/gsd:execute-phase 1` run sub-agents for phase 1
- `/gsd:verify-work 1` verify outputs against code

Key shortcuts:

- Left/Right: switch projects
- Up/Down: navigate sessions
- `!`: mark/unmark Manager session
- `x`: send sub-agent output to Manager
- `v`: cycle right pane (Manager / Both / Output / Stats)

## Contributing

We’d love your help! Ways to contribute:

- Triage issues and propose improvements
- Implement features from the Roadmap
- Improve docs and examples in `docs/` and `demos/`

Development setup:

```bash
make test     # Run tests
make fmt      # Format code
make lint     # Lint
```

Please read CONTRIBUTING.md before opening a PR. If proposing larger changes, open a discussion/issue first for alignment.

Security: If you discover a vulnerability, please open a private security advisory on GitHub rather than filing a public issue.

## Roadmap

- Rename binary/paths to `claude-ninja`
- Update installers and package distribution
- Expand docs (configuration, MCP, worktrees)
- Starter templates and example projects

## License

MIT License — see LICENSE.

---

Questions or ideas? Open an issue: https://github.com/pratham-bhatnagar/claude-ninja/issues
