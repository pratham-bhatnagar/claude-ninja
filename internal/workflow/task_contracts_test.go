package workflow

import (
	"path/filepath"
	"testing"
)

func TestTaskContractCreateUpdateLoad(t *testing.T) {
	projectPath := filepath.Join(t.TempDir(), "proj")

	created, err := EnsureTaskContract(projectPath, TaskContract{
		ID:         "t-auth-001",
		Title:      "Auth API",
		Objective:  "Build login endpoint.",
		SessionID:  "sess-1234",
		AssignedTo: "backend-agent",
		Branch:     "feat/auth",
		Worktree:   "/tmp/wt-auth",
	})
	if err != nil {
		t.Fatalf("ensure contract: %v", err)
	}
	if created.State != TaskStateReady {
		t.Fatalf("state=%s want %s", created.State, TaskStateReady)
	}

	updated, err := UpsertTaskContract(projectPath, TaskContract{
		ID:           "t-auth-001",
		State:        TaskStateWorking,
		ReviewStatus: "in_review",
		TestStatus:   "running",
		PRStatus:     PRStateOpen,
	})
	if err != nil {
		t.Fatalf("upsert contract: %v", err)
	}
	if updated.State != TaskStateWorking {
		t.Fatalf("state=%s want %s", updated.State, TaskStateWorking)
	}
	if updated.Title != "Auth API" {
		t.Fatalf("title=%q want Auth API", updated.Title)
	}

	all, err := LoadTaskContracts(projectPath)
	if err != nil {
		t.Fatalf("load contracts: %v", err)
	}
	if len(all) != 1 {
		t.Fatalf("contracts len=%d want 1", len(all))
	}
	if all[0].PRStatus != PRStateOpen {
		t.Fatalf("pr=%s want %s", all[0].PRStatus, PRStateOpen)
	}
}

func TestTaskIDForSession(t *testing.T) {
	id := TaskIDForSession("1234567890abcdef", "Frontend Shell")
	if id == "" || id[:2] != "t-" {
		t.Fatalf("unexpected task id: %q", id)
	}
}
