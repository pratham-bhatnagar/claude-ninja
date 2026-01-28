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

Sharp, fast, and focused mission control for your AI coding agents.

[![GitHub Stars](https://img.shields.io/github/stars/pratham-bhatnagar/claude-ninja?style=for-the-badge&logo=github&color=yellow&labelColor=1a1b26)](https://github.com/pratham-bhatnagar/claude-ninja/stargazers)
[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=for-the-badge&logo=go&labelColor=1a1b26)](https://go.dev)
[![License](https://img.shields.io/badge/License-MIT-9ece6a?style=for-the-badge&labelColor=1a1b26)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-macOS%20%7C%20Linux%20%7C%20WSL-7aa2f7?style=for-the-badge&labelColor=1a1b26)](https://github.com/pratham-bhatnagar/claude-ninja)
[![Latest Release](https://img.shields.io/github/v/release/pratham-bhatnagar/claude-ninja?style=for-the-badge&color=e0af68&labelColor=1a1b26)](https://github.com/pratham-bhatnagar/claude-ninja/releases)

[Features](#features) . [Quick Start](#quick-start) . [Roadmap](#roadmap) . [License](#license)

</div>

## Overview

Claude Ninja is a nimble TUI that keeps all your AI coding sessions in one clean, lightning-fast workspace. See everything at a glance, switch instantly, and keep projects organized without tab chaos.

- See status for every agent at once (running, waiting, idle)
- Fork conversations to explore different approaches without losing context
- Attach and manage MCP servers per project or globally
- Fuzzy-search across sessions and jump in milliseconds
- Use git worktrees to isolate branches per agent

## Features

- Fork sessions with full context inheritance
- Toggle MCP servers with a built-in manager
- Socket pooling for efficient MCP usage
- Tmux-integrated notification bar for waiting sessions
- Smart status detection and global search

## Quick Start

This repository is in an active rebrand. The installed binary and some file paths may still use the previous name while the codebase migrates.

```bash
make build    # Build
make test     # Test
make lint     # Lint

# Run the TUI (current binary name may still be `agent-deck`)
agent-deck
```

## Roadmap

- Rename binary and config paths to `claude-ninja`
- Update installers and Homebrew formula
- Refresh docs and examples under `docs/`

## License

MIT License — see [LICENSE](LICENSE)

---

Questions or ideas? Open an issue: https://github.com/pratham-bhatnagar/claude-ninja/issues
