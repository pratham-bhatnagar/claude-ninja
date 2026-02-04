package workflow

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const stateMachineFileName = "state_machine.json"

type ProjectState string

const (
	ProjectStateIntake     ProjectState = "intake"
	ProjectStatePlanning   ProjectState = "planning"
	ProjectStateDecomposed ProjectState = "decomposed"
	ProjectStateAssigned   ProjectState = "assigned"
	ProjectStateExecuting  ProjectState = "executing"
	ProjectStateReview     ProjectState = "review"
	ProjectStatePRFlow     ProjectState = "pr_flow"
	ProjectStateMerged     ProjectState = "merged"
	ProjectStateVerified   ProjectState = "verified"
)

type TaskState string

const (
	TaskStateReady        TaskState = "ready"
	TaskStateWorking      TaskState = "working"
	TaskStateWaitingInput TaskState = "waiting_input"
	TaskStateReviewer     TaskState = "reviewer_check"
	TaskStateDone         TaskState = "done"
)

type PRState string

const (
	PRStateBranchReady PRState = "branch_ready"
	PRStateOpen        PRState = "pr_open"
	PRStateReviewed    PRState = "pr_reviewed"
	PRStateMerged      PRState = "merged"
)

type ProjectTransition struct {
	From  ProjectState `json:"from"`
	To    ProjectState `json:"to"`
	Actor string       `json:"actor,omitempty"`
	Note  string       `json:"note,omitempty"`
	At    time.Time    `json:"at"`
}

type TaskTransition struct {
	TaskID string    `json:"task_id"`
	From   TaskState `json:"from"`
	To     TaskState `json:"to"`
	Actor  string    `json:"actor,omitempty"`
	Note   string    `json:"note,omitempty"`
	At     time.Time `json:"at"`
}

type PRTransition struct {
	TaskID string    `json:"task_id"`
	From   PRState   `json:"from"`
	To     PRState   `json:"to"`
	Actor  string    `json:"actor,omitempty"`
	Note   string    `json:"note,omitempty"`
	At     time.Time `json:"at"`
}

type TaskRuntime struct {
	ID        string    `json:"id"`
	Title     string    `json:"title,omitempty"`
	State     TaskState `json:"state"`
	Skill     string    `json:"skill,omitempty"`
	Branch    string    `json:"branch,omitempty"`
	Worktree  string    `json:"worktree,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PRRuntime struct {
	TaskID    string    `json:"task_id"`
	State     PRState   `json:"state"`
	Number    int       `json:"number,omitempty"`
	URL       string    `json:"url,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Machine struct {
	Version            int                     `json:"version"`
	Current            ProjectState            `json:"current"`
	CreatedAt          time.Time               `json:"created_at"`
	UpdatedAt          time.Time               `json:"updated_at"`
	ProjectTransitions []ProjectTransition     `json:"project_transitions"`
	TaskTransitions    []TaskTransition        `json:"task_transitions"`
	PRTransitions      []PRTransition          `json:"pr_transitions"`
	Tasks              map[string]*TaskRuntime `json:"tasks"`
	PRs                map[string]*PRRuntime   `json:"prs"`
}

type Progress struct {
	TasksTotal  int
	TasksDone   int
	PRsTotal    int
	PRsMerged   int
	Current     ProjectState
	LastUpdated time.Time
}

var (
	ErrInvalidProjectTransition = errors.New("invalid project state transition")
	ErrInvalidTaskTransition    = errors.New("invalid task state transition")
	ErrInvalidPRTransition      = errors.New("invalid pr state transition")
	ErrTaskIDRequired           = errors.New("task id is required")
)

var projectTransitions = map[ProjectState]map[ProjectState]bool{
	ProjectStateIntake:     {ProjectStatePlanning: true},
	ProjectStatePlanning:   {ProjectStateDecomposed: true},
	ProjectStateDecomposed: {ProjectStateAssigned: true},
	ProjectStateAssigned:   {ProjectStateExecuting: true},
	ProjectStateExecuting:  {ProjectStateReview: true},
	ProjectStateReview:     {ProjectStateExecuting: true, ProjectStatePRFlow: true},
	ProjectStatePRFlow:     {ProjectStateMerged: true, ProjectStateExecuting: true},
	ProjectStateMerged:     {ProjectStateVerified: true},
	ProjectStateVerified:   {},
}

var taskTransitions = map[TaskState]map[TaskState]bool{
	TaskStateReady:        {TaskStateWorking: true},
	TaskStateWorking:      {TaskStateWaitingInput: true, TaskStateReviewer: true},
	TaskStateWaitingInput: {TaskStateWorking: true},
	TaskStateReviewer:     {TaskStateWorking: true, TaskStateDone: true},
	TaskStateDone:         {},
}

var prTransitions = map[PRState]map[PRState]bool{
	PRStateBranchReady: {PRStateOpen: true},
	PRStateOpen:        {PRStateReviewed: true},
	PRStateReviewed:    {PRStateBranchReady: true, PRStateMerged: true},
	PRStateMerged:      {},
}

func NewMachine() *Machine {
	now := time.Now().UTC()
	return &Machine{
		Version:            1,
		Current:            ProjectStateIntake,
		CreatedAt:          now,
		UpdatedAt:          now,
		ProjectTransitions: []ProjectTransition{},
		TaskTransitions:    []TaskTransition{},
		PRTransitions:      []PRTransition{},
		Tasks:              map[string]*TaskRuntime{},
		PRs:                map[string]*PRRuntime{},
	}
}

func FilePath(projectPath string) string {
	return filepath.Join(projectPath, ".planning", stateMachineFileName)
}

func EnsureProjectMachine(projectPath string) error {
	if strings.TrimSpace(projectPath) == "" {
		return nil
	}
	_, err := os.Stat(FilePath(projectPath))
	if err == nil {
		return nil
	}
	if !os.IsNotExist(err) {
		return err
	}
	return Save(projectPath, NewMachine())
}

func Load(projectPath string) (*Machine, error) {
	data, err := os.ReadFile(FilePath(projectPath))
	if err != nil {
		return nil, err
	}
	var m Machine
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("decode state machine: %w", err)
	}
	if m.Tasks == nil {
		m.Tasks = map[string]*TaskRuntime{}
	}
	if m.PRs == nil {
		m.PRs = map[string]*PRRuntime{}
	}
	return &m, nil
}

func Save(projectPath string, m *Machine) error {
	if strings.TrimSpace(projectPath) == "" {
		return fmt.Errorf("project path is required")
	}
	if m == nil {
		return fmt.Errorf("machine is required")
	}
	planningDir := filepath.Join(projectPath, ".planning")
	if err := os.MkdirAll(planningDir, 0o755); err != nil {
		return err
	}
	m.UpdatedAt = time.Now().UTC()
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	tmp := filepath.Join(planningDir, stateMachineFileName+".tmp")
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, FilePath(projectPath))
}

func CanProjectTransition(from, to ProjectState) bool {
	return projectTransitions[from][to]
}

func CanTaskTransition(from, to TaskState) bool {
	return taskTransitions[from][to]
}

func CanPRTransition(from, to PRState) bool {
	return prTransitions[from][to]
}

func (m *Machine) TransitionProject(to ProjectState, actor, note string) error {
	if !CanProjectTransition(m.Current, to) {
		return fmt.Errorf("%w: %s -> %s", ErrInvalidProjectTransition, m.Current, to)
	}
	prev := m.Current
	m.Current = to
	m.ProjectTransitions = append(m.ProjectTransitions, ProjectTransition{
		From:  prev,
		To:    to,
		Actor: actor,
		Note:  note,
		At:    time.Now().UTC(),
	})
	return nil
}

func (m *Machine) EnsureTask(taskID, title string) (*TaskRuntime, error) {
	if stringsTrim(taskID) == "" {
		return nil, ErrTaskIDRequired
	}
	if t, ok := m.Tasks[taskID]; ok {
		if stringsTrim(title) != "" {
			t.Title = title
		}
		return t, nil
	}
	t := &TaskRuntime{
		ID:        taskID,
		Title:     title,
		State:     TaskStateReady,
		UpdatedAt: time.Now().UTC(),
	}
	m.Tasks[taskID] = t
	return t, nil
}

func (m *Machine) TransitionTask(taskID string, to TaskState, actor, note string) error {
	task, err := m.EnsureTask(taskID, "")
	if err != nil {
		return err
	}
	if !CanTaskTransition(task.State, to) {
		return fmt.Errorf("%w: %s -> %s (%s)", ErrInvalidTaskTransition, task.State, to, taskID)
	}
	prev := task.State
	task.State = to
	task.UpdatedAt = time.Now().UTC()
	m.TaskTransitions = append(m.TaskTransitions, TaskTransition{
		TaskID: taskID,
		From:   prev,
		To:     to,
		Actor:  actor,
		Note:   note,
		At:     task.UpdatedAt,
	})
	return nil
}

func (m *Machine) TransitionPR(taskID string, to PRState, actor, note string) error {
	if stringsTrim(taskID) == "" {
		return ErrTaskIDRequired
	}
	pr, ok := m.PRs[taskID]
	if !ok {
		pr = &PRRuntime{
			TaskID:    taskID,
			State:     PRStateBranchReady,
			UpdatedAt: time.Now().UTC(),
		}
		m.PRs[taskID] = pr
	}
	if !CanPRTransition(pr.State, to) {
		return fmt.Errorf("%w: %s -> %s (%s)", ErrInvalidPRTransition, pr.State, to, taskID)
	}
	prev := pr.State
	pr.State = to
	pr.UpdatedAt = time.Now().UTC()
	m.PRTransitions = append(m.PRTransitions, PRTransition{
		TaskID: taskID,
		From:   prev,
		To:     to,
		Actor:  actor,
		Note:   note,
		At:     pr.UpdatedAt,
	})
	return nil
}

func (m *Machine) Progress() Progress {
	p := Progress{
		Current:     m.Current,
		LastUpdated: m.UpdatedAt,
	}
	for _, t := range m.Tasks {
		p.TasksTotal++
		if t.State == TaskStateDone {
			p.TasksDone++
		}
	}
	for _, pr := range m.PRs {
		p.PRsTotal++
		if pr.State == PRStateMerged {
			p.PRsMerged++
		}
	}
	return p
}

func (m *Machine) SortedTasks() []*TaskRuntime {
	out := make([]*TaskRuntime, 0, len(m.Tasks))
	for _, t := range m.Tasks {
		out = append(out, t)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].ID < out[j].ID
	})
	return out
}

func stringsTrim(s string) string { return strings.TrimSpace(s) }
