package ui

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/pratham-bhatnagar/claude-ninja/internal/session"
	"github.com/pratham-bhatnagar/claude-ninja/internal/workflow"
)

func TestEnsureOrchestratorPrefersClaude(t *testing.T) {
	home := &Home{
		instances:    []*session.Instance{},
		instanceByID: make(map[string]*session.Instance),
	}
	projectPath := "/tmp/project-alpha"

	codex := &session.Instance{ID: "c1", Title: "Codex", Tool: "codex", ProjectPath: projectPath, CreatedAt: time.Now()}
	claude := &session.Instance{ID: "c2", Title: "Claude", Tool: "claude", ProjectPath: projectPath, CreatedAt: time.Now().Add(-time.Minute)}
	home.instances = []*session.Instance{codex, claude}

	orchestrator := home.ensureOrchestratorForProject(projectPath)
	if orchestrator == nil {
		t.Fatal("expected orchestrator, got nil")
	}
	if orchestrator.ID != claude.ID {
		t.Fatalf("expected Claude orchestrator, got %s", orchestrator.ID)
	}
	if !claude.IsManager {
		t.Fatalf("expected Claude to be marked as orchestrator")
	}
}

func TestEnsureOrchestratorMostRecentActivity(t *testing.T) {
	home := &Home{
		instances:    []*session.Instance{},
		instanceByID: make(map[string]*session.Instance),
	}
	projectPath := "/tmp/project-beta"

	old := &session.Instance{ID: "o1", Title: "Old", Tool: "codex", ProjectPath: projectPath, CreatedAt: time.Now().Add(-2 * time.Hour)}
	newer := &session.Instance{ID: "o2", Title: "Newer", Tool: "codex", ProjectPath: projectPath, CreatedAt: time.Now().Add(-time.Minute)}
	home.instances = []*session.Instance{old, newer}

	orchestrator := home.ensureOrchestratorForProject(projectPath)
	if orchestrator == nil {
		t.Fatal("expected orchestrator, got nil")
	}
	if orchestrator.ID != newer.ID {
		t.Fatalf("expected most recent orchestrator, got %s", orchestrator.ID)
	}
}

func TestGetPhaseSummary(t *testing.T) {
	tempDir := t.TempDir()
	planningDir := filepath.Join(tempDir, ".planning")
	if err := os.MkdirAll(planningDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	roadmap := "## Phase 1: Foundation\n\n## Phase 2: Features\n"
	state := "Current Phase: 1\n"
	if err := os.WriteFile(filepath.Join(planningDir, "ROADMAP.md"), []byte(roadmap), 0o644); err != nil {
		t.Fatalf("write roadmap: %v", err)
	}
	if err := os.WriteFile(filepath.Join(planningDir, "STATE.md"), []byte(state), 0o644); err != nil {
		t.Fatalf("write state: %v", err)
	}

	home := &Home{}
	summary := home.getPhaseSummary(tempDir)
	if summary != "Phase 1 of 2" {
		t.Fatalf("unexpected summary: %q", summary)
	}
}

func TestGetPhaseSummaryFromStateMachine(t *testing.T) {
	projectDir := t.TempDir()
	if err := workflow.EnsureProjectMachine(projectDir); err != nil {
		t.Fatalf("ensure machine: %v", err)
	}
	m, err := workflow.Load(projectDir)
	if err != nil {
		t.Fatalf("load machine: %v", err)
	}
	if _, err := m.EnsureTask("T-001", "frontend"); err != nil {
		t.Fatalf("ensure task: %v", err)
	}
	if err := m.TransitionProject(workflow.ProjectStatePlanning, "manager", ""); err != nil {
		t.Fatalf("transition project: %v", err)
	}
	if err := m.TransitionTask("T-001", workflow.TaskStateWorking, "agent", ""); err != nil {
		t.Fatalf("transition task: %v", err)
	}
	if err := workflow.Save(projectDir, m); err != nil {
		t.Fatalf("save machine: %v", err)
	}

	home := &Home{}
	summary := home.getPhaseSummary(projectDir)
	if summary != "planning · 0/1 tasks done" {
		t.Fatalf("unexpected machine summary: %q", summary)
	}
}
