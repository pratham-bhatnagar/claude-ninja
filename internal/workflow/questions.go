package workflow

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const questionsFileName = "questions.json"

type HumanQuestion struct {
	ID        string    `json:"id"`
	TaskID    string    `json:"task_id"`
	SessionID string    `json:"session_id"`
	Question  string    `json:"question"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func questionsPath(projectPath string) string {
	return filepath.Join(projectPath, ".planning", questionsFileName)
}

func LoadQuestions(projectPath string) ([]HumanQuestion, error) {
	data, err := os.ReadFile(questionsPath(projectPath))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var q []HumanQuestion
	if err := json.Unmarshal(data, &q); err != nil {
		return nil, err
	}
	return q, nil
}

func AddQuestion(projectPath, taskID, sessionID, question string) error {
	if question == "" {
		return nil
	}
	items, err := LoadQuestions(projectPath)
	if err != nil {
		return err
	}
	id := fmt.Sprintf("q-%d", time.Now().UTC().UnixNano())
	items = append(items, HumanQuestion{
		ID:        id,
		TaskID:    taskID,
		SessionID: sessionID,
		Question:  question,
		Status:    "open",
		CreatedAt: time.Now().UTC(),
	})
	if err := os.MkdirAll(filepath.Join(projectPath, ".planning"), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(questionsPath(projectPath), b, 0o644)
}

func LoadOpenQuestions(projectPath string) ([]HumanQuestion, error) {
	items, err := LoadQuestions(projectPath)
	if err != nil {
		return nil, err
	}
	var out []HumanQuestion
	for _, q := range items {
		if q.Status == "" || q.Status == "open" {
			out = append(out, q)
		}
	}
	return out, nil
}
