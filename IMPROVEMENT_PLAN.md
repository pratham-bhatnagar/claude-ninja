# Claude Ninja Improvement Plan

## Executive Summary

After testing Claude Ninja's autonomous agent system, I've identified critical issues preventing the codebase from compiling and running. This document outlines the problems found and the iterative improvement plan.

---

## Phase 1: Critical Fixes (Current)

### Issue 1: Module Name Mismatch 🚨 CRITICAL
**Problem:** The codebase has inconsistent module naming:
- `go.mod` declares: `github.com/asheshgoplani/agent-deck`
- Code imports: `github.com/pratham-bhatnagar/claude-ninja`

**Impact:** Project cannot compile, go mod tidy fails, external dependencies get pulled incorrectly.

**Fix:** 
1. Update `go.mod` module name to `github.com/pratham-bhatnagar/claude-ninja`
2. Replace all import paths in Go files
3. Remove dependency on external `github.com/asheshgoplani/agent-deck`

### Issue 2: Missing Struct Fields and Methods 🚨 CRITICAL
**Problem:** Code references non-existent fields/methods:

| Location | Missing Reference | Context |
|----------|------------------|---------|
| `internal/ui/home.go:4735` | `settings.AttachMode` | InstanceSettings struct |
| `internal/ui/home.go:4737` | `tmux.IsInsideTmux()` | tmux package function |
| `internal/ui/home.go:4742` | `tmux.IsInsideTmux()` | tmux package function |
| `internal/ui/home.go:4760` | `tmux.IsInsideTmux()` | tmux package function |
| `internal/ui/home.go:4776` | `tmux.DisplayPopup()` | tmux package function |
| `internal/ui/home.go:4776` | `settings.PopupWidth` | InstanceSettings struct |
| `internal/ui/home.go:4776` | `settings.PopupHeight` | InstanceSettings struct |
| `internal/ui/home.go:5464` | `agents[i].IsManager` | Instance struct field |
| `internal/ui/home.go:5467` | `agents[i].IsManager` | Instance struct field |

**Impact:** Build fails with 10+ compilation errors.

**Fix Options:**
1. **Option A:** Add missing fields to structs
2. **Option B:** Comment out unused code paths
3. **Option C:** Implement stub methods that return safe defaults

**Recommendation:** Implement Option A for core functionality, stub out experimental features.

---

## Phase 2: Architecture Improvements

### Improvement 1: State Machine Robustness
**Current:** State machine JSON file can get corrupted during concurrent access.

**Improvement:**
- Add file locking for state_machine.json
- Implement atomic writes (write to temp, then rename)
- Add state machine versioning and migration

### Improvement 2: Error Handling
**Current:** Many functions return generic errors without context.

**Improvement:**
- Wrap errors with `fmt.Errorf("context: %w", err)`
- Add structured logging with zap or slog
- Implement error categorization (user vs system errors)

### Improvement 3: Test Coverage
**Current:** Tests exist but coverage is inconsistent.

**Improvement:**
- Add unit tests for state_machine.go transitions
- Add integration tests for full agent lifecycle
- Add benchmarks for concurrent agent operations

---

## Phase 3: Autonomous Agent Loop

### Current State Machine Flow

```
Project Lifecycle:
Intake → Planning → Decomposed → Assigned → Executing → Review → PRFlow → Merged → Verified

Task Lifecycle:
Ready → Working → (WaitingInput) → ReviewerCheck → Done

PR Lifecycle:
BranchReady → Open → Reviewed → Merged
```

### Proposed Enhancements

1. **Add Retry Logic:** Tasks that fail should auto-retry with exponential backoff
2. **Add Timeout Handling:** Tasks stuck in "Working" for >30min should be flagged
3. **Add Dependency Graph:** Visualize task dependencies before execution
4. **Add Metrics:** Track cycle time, success rate, agent utilization

---

## Phase 4: Gas Town Integration

### Current Integration
- Beads database initialized ✓
- Convoy tracking working ✓
- Agent beads created ✓

### Improvements Needed

1. **Hook up Discord notifications:** Already configured in cron jobs
2. **Add automated PR creation:** Use `gt` commands to open PRs
3. **Add self-healing:** Detect stuck agents and respawn
4. **Add work balancing:** Distribute tasks across polecats evenly

---

## Implementation Priority

### Week 1: Critical Fixes
- [ ] Fix module name mismatch
- [ ] Add missing struct fields
- [ ] Implement missing tmux functions
- [ ] Get project to compile

### Week 2: Testing
- [ ] Add unit tests for fixed code
- [ ] Run full integration test
- [ ] Test autonomous agent loop
- [ ] Document test results

### Week 3: Enhancements
- [ ] Add retry logic
- [ ] Add timeout handling
- [ ] Add metrics collection
- [ ] Create monitoring dashboard

### Week 4: Integration
- [ ] Hook up Discord channel
- [ ] Test end-to-end flow
- [ ] Create demo video/screenshots
- [ ] Write final documentation

---

## Testing Checklist

### Build Tests
- [ ] `make build` completes without errors
- [ ] `make test` passes all tests
- [ ] `make lint` shows no issues
- [ ] Binary runs: `./build/agent-deck --help`

### Unit Tests
- [ ] State machine transitions work correctly
- [ ] Task contracts are validated
- [ ] PR state changes propagate

### Integration Tests
- [ ] Create project with `/ninja:new-project`
- [ ] Plan phase with `/ninja:plan-phase 1`
- [ ] Execute with `/ninja:execute-phase 1`
- [ ] Verify with `/ninja:verify-work 1`

### Autonomous Agent Tests
- [ ] Spawn 10 polecats
- [ ] Assign tasks to all polecats
- [ ] Monitor completion
- [ ] Verify PR creation

---

## Success Metrics

| Metric | Current | Target |
|--------|---------|--------|
| Build Status | ❌ Fails | ✅ Passes |
| Test Coverage | ~60% | >80% |
| Compilation Errors | 10+ | 0 |
| Autonomous Task Completion | N/A | 90%+ |
| Discord Updates | Manual | Automated |

---

## Screenshots Required for PR

1. **Build Success:** Terminal showing `make build` completing
2. **Test Results:** Terminal showing `make test` passing
3. **TUI Screenshot:** Claude Ninja interface running
4. **State Machine:** Diagram or JSON showing state transitions
5. **Agent Swarm:** Multiple agents working concurrently
6. **Discord Integration:** Bot posting updates to channel

---

## Conclusion

Claude Ninja has a solid architecture but needs critical fixes before it's production-ready. The autonomous agent system design is sound, but implementation gaps (missing fields, module mismatches) prevent it from running. With the fixes outlined above, it can become a powerful tool for orchestrating AI coding agents.

**Next Steps:**
1. Implement Phase 1 fixes (this PR)
2. Test on local Mac system
3. Create PR with screenshots
4. Iterate based on review comments
5. Hook up Discord for real-time updates
