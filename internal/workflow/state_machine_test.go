package workflow

import (
	"errors"
	"path/filepath"
	"testing"
)

func TestProjectTransitionsHappyPath(t *testing.T) {
	m := NewMachine()
	path := []ProjectState{
		ProjectStatePlanning,
		ProjectStateDecomposed,
		ProjectStateAssigned,
		ProjectStateExecuting,
		ProjectStateReview,
		ProjectStatePRFlow,
		ProjectStateMerged,
		ProjectStateVerified,
	}
	for _, to := range path {
		if err := m.TransitionProject(to, "manager", ""); err != nil {
			t.Fatalf("transition to %s failed: %v", to, err)
		}
	}
	if m.Current != ProjectStateVerified {
		t.Fatalf("current state = %s, want %s", m.Current, ProjectStateVerified)
	}
}

func TestProjectTransitionInvalid(t *testing.T) {
	m := NewMachine()
	err := m.TransitionProject(ProjectStateExecuting, "manager", "")
	if !errors.Is(err, ErrInvalidProjectTransition) {
		t.Fatalf("expected ErrInvalidProjectTransition, got %v", err)
	}
}

func TestTaskAndPRLoops(t *testing.T) {
	m := NewMachine()
	if _, err := m.EnsureTask("T-001", "frontend shell"); err != nil {
		t.Fatalf("ensure task: %v", err)
	}
	if err := m.TransitionTask("T-001", TaskStateWorking, "agent", ""); err != nil {
		t.Fatalf("task working: %v", err)
	}
	if err := m.TransitionTask("T-001", TaskStateReviewer, "agent", "ready for review"); err != nil {
		t.Fatalf("task reviewer: %v", err)
	}
	if err := m.TransitionTask("T-001", TaskStateWorking, "reviewer", "requested changes"); err != nil {
		t.Fatalf("task back to working: %v", err)
	}
	if err := m.TransitionTask("T-001", TaskStateReviewer, "agent", "resubmitted"); err != nil {
		t.Fatalf("task reviewer 2: %v", err)
	}
	if err := m.TransitionTask("T-001", TaskStateDone, "reviewer", "accepted"); err != nil {
		t.Fatalf("task done: %v", err)
	}

	if err := m.TransitionPR("T-001", PRStateOpen, "agent", "created PR"); err != nil {
		t.Fatalf("pr open: %v", err)
	}
	if err := m.TransitionPR("T-001", PRStateReviewed, "manager", "left comments"); err != nil {
		t.Fatalf("pr reviewed: %v", err)
	}
	if err := m.TransitionPR("T-001", PRStateBranchReady, "agent", "addressed comments"); err != nil {
		t.Fatalf("pr back to branch ready: %v", err)
	}
	if err := m.TransitionPR("T-001", PRStateOpen, "agent", "pushed updates"); err != nil {
		t.Fatalf("pr reopen: %v", err)
	}
	if err := m.TransitionPR("T-001", PRStateReviewed, "manager", "approved"); err != nil {
		t.Fatalf("pr reviewed 2: %v", err)
	}
	if err := m.TransitionPR("T-001", PRStateMerged, "manager", "merged"); err != nil {
		t.Fatalf("pr merged: %v", err)
	}

	progress := m.Progress()
	if progress.TasksTotal != 1 || progress.TasksDone != 1 {
		t.Fatalf("unexpected task progress: %+v", progress)
	}
	if progress.PRsTotal != 1 || progress.PRsMerged != 1 {
		t.Fatalf("unexpected PR progress: %+v", progress)
	}
}

func TestPersistenceRoundTrip(t *testing.T) {
	dir := t.TempDir()
	projectPath := filepath.Join(dir, "project")
	if err := EnsureProjectMachine(projectPath); err != nil {
		t.Fatalf("ensure machine: %v", err)
	}

	m, err := Load(projectPath)
	if err != nil {
		t.Fatalf("load machine: %v", err)
	}
	if m.Current != ProjectStateIntake {
		t.Fatalf("current state = %s, want %s", m.Current, ProjectStateIntake)
	}

	if err := m.TransitionProject(ProjectStatePlanning, "manager", "kickoff"); err != nil {
		t.Fatalf("transition project: %v", err)
	}
	if err := Save(projectPath, m); err != nil {
		t.Fatalf("save machine: %v", err)
	}

	reloaded, err := Load(projectPath)
	if err != nil {
		t.Fatalf("reload machine: %v", err)
	}
	if reloaded.Current != ProjectStatePlanning {
		t.Fatalf("reloaded state = %s, want %s", reloaded.Current, ProjectStatePlanning)
	}
	if len(reloaded.ProjectTransitions) != 1 {
		t.Fatalf("project transitions = %d, want 1", len(reloaded.ProjectTransitions))
	}
}
