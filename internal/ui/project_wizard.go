package ui

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	ProjectModeNew  = 0
	ProjectModeOpen = 1
)

type ProjectWizard struct {
	mode          int
	nameInput     textinput.Model
	pathInput     textinput.Model
	focusIndex    int
	width         int
	height        int
	visible       bool
	validationErr string
}

func NewProjectWizard() *ProjectWizard {
	nameInput := textinput.New()
	nameInput.Placeholder = "project-name (optional)"
	nameInput.CharLimit = 80
	nameInput.Width = 40

	pathInput := textinput.New()
	pathInput.Placeholder = "~/projects/my-app"
	pathInput.CharLimit = 256
	pathInput.Width = 40

	cwd, err := os.Getwd()
	if err == nil {
		pathInput.SetValue(cwd)
	}

	return &ProjectWizard{
		mode:       ProjectModeNew,
		nameInput:  nameInput,
		pathInput:  pathInput,
		focusIndex: 0,
		visible:    false,
	}
}

func (w *ProjectWizard) SetSize(width, height int) {
	w.width = width
	w.height = height
}

func (w *ProjectWizard) Show() {
	w.visible = true
	w.validationErr = ""
	w.focusIndex = 0
	w.nameInput.SetValue("")
	w.nameInput.Focus()
	w.pathInput.Blur()
	cwd, err := os.Getwd()
	if err == nil {
		w.pathInput.SetValue(cwd)
	}
}

func (w *ProjectWizard) Hide() {
	w.visible = false
}

func (w *ProjectWizard) IsVisible() bool {
	return w.visible
}

func (w *ProjectWizard) Update(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "tab":
		w.focusIndex = (w.focusIndex + 1) % 2
	case "shift+tab":
		w.focusIndex--
		if w.focusIndex < 0 {
			w.focusIndex = 1
		}
	case "up", "k":
		w.mode = ProjectModeNew
	case "down", "j":
		w.mode = ProjectModeOpen
	}

	if w.focusIndex == 0 {
		w.nameInput.Focus()
		w.pathInput.Blur()
	} else {
		w.pathInput.Focus()
		w.nameInput.Blur()
	}

	var cmd tea.Cmd
	w.nameInput, cmd = w.nameInput.Update(msg)
	w.pathInput, _ = w.pathInput.Update(msg)
	return cmd
}

func (w *ProjectWizard) SetError(err string) {
	w.validationErr = err
}

func (w *ProjectWizard) GetValues() (mode int, name, path string) {
	mode = w.mode
	name = strings.TrimSpace(w.nameInput.Value())
	path = strings.Trim(strings.TrimSpace(w.pathInput.Value()), "'\"")

	if idx := strings.Index(path, "~/"); idx > 0 {
		path = path[idx:]
	}
	if strings.HasPrefix(path, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			path = filepath.Join(home, path[2:])
		}
	} else if path == "~" {
		if home, err := os.UserHomeDir(); err == nil {
			path = home
		}
	}
	path = filepath.Clean(path)

	return mode, name, path
}

func (w *ProjectWizard) View() string {
	if !w.visible {
		return ""
	}

	title := "Create or Open Project"
	boxWidth := min(70, max(40, w.width-8))
	boxHeight := min(18, max(12, w.height-8))

	modeStyle := lipgloss.NewStyle().Bold(true).Foreground(ColorAccent)
	labelStyle := lipgloss.NewStyle().Foreground(ColorTextDim)
	inputStyle := lipgloss.NewStyle().Foreground(ColorText)

	modeNew := "○ New project"
	modeOpen := "○ Open existing"
	if w.mode == ProjectModeNew {
		modeNew = "● New project"
	}
	if w.mode == ProjectModeOpen {
		modeOpen = "● Open existing"
	}

	var b strings.Builder
	b.WriteString(modeStyle.Render(title))
	b.WriteString("\n\n")
	b.WriteString(labelStyle.Render("Mode:"))
	b.WriteString("\n")
	b.WriteString("  " + inputStyle.Render(modeNew))
	b.WriteString("\n")
	b.WriteString("  " + inputStyle.Render(modeOpen))
	b.WriteString("\n\n")
	b.WriteString(labelStyle.Render("Project name (optional):"))
	b.WriteString("\n")
	b.WriteString(w.nameInput.View())
	b.WriteString("\n\n")
	b.WriteString(labelStyle.Render("Project path:"))
	b.WriteString("\n")
	b.WriteString(w.pathInput.View())
	b.WriteString("\n")

	hintStyle := lipgloss.NewStyle().Foreground(ColorComment)
	b.WriteString("\n")
	b.WriteString(hintStyle.Render("Enter = create/open  •  Tab = next  •  ↑/↓ = mode  •  Esc = cancel"))

	if w.validationErr != "" {
		errStyle := lipgloss.NewStyle().Foreground(ColorRed)
		b.WriteString("\n")
		b.WriteString(errStyle.Render(w.validationErr))
	}

	box := lipgloss.NewStyle().
		Width(boxWidth).
		Height(boxHeight).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorder).
		Render(b.String())

	return lipgloss.Place(w.width, w.height, lipgloss.Center, lipgloss.Center, box)
}
