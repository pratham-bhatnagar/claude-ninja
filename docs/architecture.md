# Architecture Overview (ASCII)

```
User
  |
  v
+----------------------------------+
| agent-deck (TUI/CLI mission ctrl)|
| sessions, groups, search, fork   |
+------------------+---------------+
                   |
                   v
+----------------------------------+
| Worktree/Branch Manager          |
| one agent = one worktree/branch  |
+------------------+---------------+
                   |
                   v
+----------------------------------+
| Manager Session (Claude Code)    |
| reads .planning/, assigns tasks  |
| nudges agents, resolves conflicts|
+------------------+---------------+
                   |
        +----------+----------+
        |                     |
        v                     v
+-------------------+  +-------------------+
| Agent Sessions     |  | GSD Subagents     |
| (Claude Code)      |  | planner/executor/ |
| work in branches   |  | verifier          |
+-------------------+  +-------------------+
                   |
                   v
           .planning/ state files
```

## What each part does

- **agent-deck TUI/CLI**: Central dashboard to start, monitor, and switch between agent sessions. It also handles session grouping, search, and fork workflows.
- **Worktree/Branch Manager**: Ensures every agent runs in its own Git worktree and branch, so work never overwrites other agents.
- **Manager Session**: A dedicated Claude Code session that reads `.planning/` state, assigns tasks to agents, and nudges them when blocked. It also decides how to resolve conflicts.
- **Agent Sessions**: Individual Claude Code sessions that implement tasks inside their own branches.
- **GSD Subagents**: Planner/executor/verifier agents that GSD spawns to keep execution structured and to verify outcomes against actual code.
- **.planning/**: Canonical project state (PROJECT, REQUIREMENTS, ROADMAP, STATE, VERIFICATION).

## How review works (simple)

- GSD runs a verifier subagent after execution. The verifier checks actual code changes and writes a verification report in `.planning/VERIFICATION.md`.
- If gaps are found, GSD generates fix plans. If user input is needed, the Manager asks you directly.

## When an agent needs you

- agent-deck marks the session as **Waiting**.
- The Manager session surfaces the request and routes your answer back to that agent so work can resume.

## If agents conflict

- Conflicts are isolated because each agent works in its own worktree branch.
- The Manager compares outcomes and picks a winning branch (or asks you to choose), then merges it.
