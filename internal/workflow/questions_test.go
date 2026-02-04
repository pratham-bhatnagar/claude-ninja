package workflow

import (
	"path/filepath"
	"testing"
)

func TestAddAndLoadQuestions(t *testing.T) {
	projectPath := filepath.Join(t.TempDir(), "project")

	if err := AddQuestion(projectPath, "t-auth-1", "sess-1", "Need API contract confirmation"); err != nil {
		t.Fatalf("add question: %v", err)
	}
	if err := AddQuestion(projectPath, "t-auth-2", "sess-2", ""); err != nil {
		t.Fatalf("empty add should be no-op: %v", err)
	}

	items, err := LoadOpenQuestions(projectPath)
	if err != nil {
		t.Fatalf("load open questions: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("questions len=%d want 1", len(items))
	}
	if items[0].TaskID != "t-auth-1" {
		t.Fatalf("task id=%q want t-auth-1", items[0].TaskID)
	}
}
