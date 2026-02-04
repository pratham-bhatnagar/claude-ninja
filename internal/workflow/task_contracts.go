package workflow

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type TaskContract struct {
	ID            string
	Title         string
	Objective     string
	Outputs       string
	DoneWhen      string
	AssignedTo    string
	SessionID     string
	Branch        string
	Worktree      string
	State         TaskState
	ReviewStatus  string
	TestStatus    string
	PRStatus      PRState
	UpdatedAtText string
	Path          string
}

func TaskIDForSession(sessionID, title string) string {
	base := strings.TrimSpace(sessionID)
	if len(base) > 8 {
		base = base[:8]
	}
	base = strings.ToLower(base)
	if base == "" {
		base = "task"
	}
	if title != "" {
		slug := strings.ToLower(title)
		slug = strings.NewReplacer(" ", "-", "_", "-", "/", "-").Replace(slug)
		var kept []rune
		for _, r := range slug {
			if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
				kept = append(kept, r)
			}
		}
		slug = strings.Trim(strings.Join(strings.Fields(strings.ReplaceAll(string(kept), "-", " ")), "-"), "-")
		if len(slug) > 18 {
			slug = slug[:18]
		}
		if slug != "" {
			return fmt.Sprintf("t-%s-%s", slug, base)
		}
	}
	return "t-" + base
}

func contractsDir(projectPath string) string {
	return filepath.Join(projectPath, ".planning", "tasks")
}

func contractPath(projectPath, taskID string) string {
	return filepath.Join(contractsDir(projectPath), taskID+".md")
}

func EnsureTaskContract(projectPath string, c TaskContract) (TaskContract, error) {
	if strings.TrimSpace(projectPath) == "" {
		return c, fmt.Errorf("project path is required")
	}
	if strings.TrimSpace(c.ID) == "" {
		return c, fmt.Errorf("task id is required")
	}
	if c.State == "" {
		c.State = TaskStateReady
	}
	if c.PRStatus == "" {
		c.PRStatus = PRStateBranchReady
	}
	if c.ReviewStatus == "" {
		c.ReviewStatus = "pending"
	}
	if c.TestStatus == "" {
		c.TestStatus = "pending"
	}
	if c.UpdatedAtText == "" {
		c.UpdatedAtText = time.Now().UTC().Format(time.RFC3339)
	}
	if c.Title == "" {
		c.Title = c.ID
	}
	if c.Objective == "" {
		c.Objective = "Define objective."
	}
	if c.Outputs == "" {
		c.Outputs = "Deliver code + summary."
	}
	if c.DoneWhen == "" {
		c.DoneWhen = "Code reviewed and tests passing."
	}

	if err := os.MkdirAll(contractsDir(projectPath), 0o755); err != nil {
		return c, err
	}
	c.Path = contractPath(projectPath, c.ID)

	if _, err := os.Stat(c.Path); err == nil {
		return LoadTaskContract(c.Path)
	} else if !os.IsNotExist(err) {
		return c, err
	}
	if err := writeContract(c.Path, c); err != nil {
		return c, err
	}
	return c, nil
}

func UpsertTaskContract(projectPath string, c TaskContract) (TaskContract, error) {
	if strings.TrimSpace(projectPath) == "" {
		return c, fmt.Errorf("project path is required")
	}
	if strings.TrimSpace(c.ID) == "" {
		return c, fmt.Errorf("task id is required")
	}
	path := contractPath(projectPath, c.ID)
	if existing, err := LoadTaskContract(path); err == nil {
		if c.Title == "" {
			c.Title = existing.Title
		}
		if c.Objective == "" {
			c.Objective = existing.Objective
		}
		if c.Outputs == "" {
			c.Outputs = existing.Outputs
		}
		if c.DoneWhen == "" {
			c.DoneWhen = existing.DoneWhen
		}
		if c.AssignedTo == "" {
			c.AssignedTo = existing.AssignedTo
		}
		if c.SessionID == "" {
			c.SessionID = existing.SessionID
		}
		if c.Branch == "" {
			c.Branch = existing.Branch
		}
		if c.Worktree == "" {
			c.Worktree = existing.Worktree
		}
		if c.State == "" {
			c.State = existing.State
		}
		if c.ReviewStatus == "" {
			c.ReviewStatus = existing.ReviewStatus
		}
		if c.TestStatus == "" {
			c.TestStatus = existing.TestStatus
		}
		if c.PRStatus == "" {
			c.PRStatus = existing.PRStatus
		}
	}
	if c.State == "" {
		c.State = TaskStateReady
	}
	if c.PRStatus == "" {
		c.PRStatus = PRStateBranchReady
	}
	if c.ReviewStatus == "" {
		c.ReviewStatus = "pending"
	}
	if c.TestStatus == "" {
		c.TestStatus = "pending"
	}
	c.UpdatedAtText = time.Now().UTC().Format(time.RFC3339)
	if c.Path == "" {
		c.Path = path
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return c, err
	}
	if err := writeContract(path, c); err != nil {
		return c, err
	}
	return c, nil
}

func LoadTaskContracts(projectPath string) ([]TaskContract, error) {
	dir := contractsDir(projectPath)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var out []TaskContract
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		c, err := LoadTaskContract(filepath.Join(dir, e.Name()))
		if err == nil && c.ID != "" {
			out = append(out, c)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}

func LoadTaskContract(path string) (TaskContract, error) {
	f, err := os.Open(path)
	if err != nil {
		return TaskContract{}, err
	}
	defer f.Close()

	var c TaskContract
	c.Path = path
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		switch {
		case strings.HasPrefix(line, "# Task "):
			header := strings.TrimPrefix(line, "# Task ")
			parts := strings.SplitN(header, ": ", 2)
			c.ID = strings.TrimSpace(parts[0])
			if len(parts) > 1 {
				c.Title = strings.TrimSpace(parts[1])
			}
		case strings.HasPrefix(line, "- State: "):
			c.State = TaskState(strings.TrimSpace(strings.TrimPrefix(line, "- State: ")))
		case strings.HasPrefix(line, "- Review: "):
			c.ReviewStatus = strings.TrimSpace(strings.TrimPrefix(line, "- Review: "))
		case strings.HasPrefix(line, "- Tests: "):
			c.TestStatus = strings.TrimSpace(strings.TrimPrefix(line, "- Tests: "))
		case strings.HasPrefix(line, "- PR: "):
			c.PRStatus = PRState(strings.TrimSpace(strings.TrimPrefix(line, "- PR: ")))
		case strings.HasPrefix(line, "- Assigned To: "):
			c.AssignedTo = strings.TrimSpace(strings.TrimPrefix(line, "- Assigned To: "))
		case strings.HasPrefix(line, "- Session ID: "):
			c.SessionID = strings.TrimSpace(strings.TrimPrefix(line, "- Session ID: "))
		case strings.HasPrefix(line, "- Branch: "):
			c.Branch = strings.TrimSpace(strings.TrimPrefix(line, "- Branch: "))
		case strings.HasPrefix(line, "- Worktree: "):
			c.Worktree = strings.TrimSpace(strings.TrimPrefix(line, "- Worktree: "))
		case strings.HasPrefix(line, "- Updated At: "):
			c.UpdatedAtText = strings.TrimSpace(strings.TrimPrefix(line, "- Updated At: "))
		}
	}
	if err := sc.Err(); err != nil {
		return TaskContract{}, err
	}
	if c.State == "" {
		c.State = TaskStateReady
	}
	if c.PRStatus == "" {
		c.PRStatus = PRStateBranchReady
	}
	if c.ReviewStatus == "" {
		c.ReviewStatus = "pending"
	}
	if c.TestStatus == "" {
		c.TestStatus = "pending"
	}
	return c, nil
}

func writeContract(path string, c TaskContract) error {
	var b strings.Builder
	fmt.Fprintf(&b, "# Task %s: %s\n\n", c.ID, c.Title)
	fmt.Fprintf(&b, "- State: %s\n", c.State)
	fmt.Fprintf(&b, "- Review: %s\n", c.ReviewStatus)
	fmt.Fprintf(&b, "- Tests: %s\n", c.TestStatus)
	fmt.Fprintf(&b, "- PR: %s\n", c.PRStatus)
	fmt.Fprintf(&b, "- Assigned To: %s\n", c.AssignedTo)
	fmt.Fprintf(&b, "- Session ID: %s\n", c.SessionID)
	fmt.Fprintf(&b, "- Branch: %s\n", c.Branch)
	fmt.Fprintf(&b, "- Worktree: %s\n", c.Worktree)
	fmt.Fprintf(&b, "- Updated At: %s\n\n", c.UpdatedAtText)
	fmt.Fprintf(&b, "## Objective\n%s\n\n", c.Objective)
	fmt.Fprintf(&b, "## Expected Output\n%s\n\n", c.Outputs)
	fmt.Fprintf(&b, "## Done When\n%s\n", c.DoneWhen)
	return os.WriteFile(path, []byte(b.String()), 0o644)
}
