# Claude Ninja - Inspired by this subreddit

Hey everyone, I wanted to share a project I built after learning from a lot of posts and discussions in this subreddit.

## What it is

Claude Ninja is a project-first, terminal-based workspace for orchestrating multiple AI coding agents.

- One project can run many agents in parallel
- Each agent works in its own git branch/worktree
- A manager/orchestrator coordinates planning, execution, and verification
- Waiting/stalled agents are surfaced so work keeps moving

## Why I built it

I kept running into the same issues while using agents at scale:

- Context gets fragmented across sessions
- Parallel work can overwrite or conflict
- Planning quality varies a lot
- Reviewer feedback loops are inconsistent

So I started building a manager-driven system that keeps planning and execution structured.

## Vision

The direction I am actively building toward:

- High-quality plan generation from a single prompt
- Automatic decomposition into sub-tasks
- One sub-agent per sub-task with explicit output contracts
- Reviewer agent loops before manager acceptance
- PR-driven integration with comment -> fix -> re-review cycles
- Optional skill lookup per sub-task (`skills.sh`) to pick the best workflow/tooling

## Obsidian direction

I am also exploring an Obsidian-style knowledge layer:

- Markdown-first task sheets
- Graph view for manager <-> agents <-> tasks <-> PRs
- Better visibility into progress and dependency flow

## Feedback request

Would love feedback on:

1. The manager/sub-agent/reviewer architecture
2. The best way to model progress and quality gates
3. Obsidian integration ideas that would make this genuinely useful daily

If this sounds interesting, I can also share architecture diagrams and the state-machine model I am using next.
