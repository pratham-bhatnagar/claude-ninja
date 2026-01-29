package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type OrchestratorInputDialog struct {
	input   textarea.Model
	width   int
	height  int
	visible bool
}

func NewOrchestratorInputDialog() *OrchestratorInputDialog {
	input := textarea.New()
	input.Placeholder = "Message to the Orchestrator…"
	input.CharLimit = 4000
	input.Focus()
	input.ShowLineNumbers = false
	input.Prompt = "> "
	return &OrchestratorInputDialog{
		input:   input,
		visible: false,
	}
}

func (d *OrchestratorInputDialog) SetSize(width, height int) {
	d.width = width
	d.height = height
}

func (d *OrchestratorInputDialog) Show() {
	d.visible = true
	d.input.SetValue("")
	d.input.Focus()
}

func (d *OrchestratorInputDialog) Hide() {
	d.visible = false
}

func (d *OrchestratorInputDialog) IsVisible() bool {
	return d.visible
}

func (d *OrchestratorInputDialog) Value() string {
	return strings.TrimSpace(d.input.Value())
}

func (d *OrchestratorInputDialog) Update(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	d.input, cmd = d.input.Update(msg)
	return cmd
}

func (d *OrchestratorInputDialog) View() string {
	if !d.visible {
		return ""
	}

	boxWidth := min(80, max(50, d.width-8))
	boxHeight := min(18, max(10, d.height-8))
	inputWidth := boxWidth - 6
	inputHeight := boxHeight - 8
	if inputHeight < 3 {
		inputHeight = 3
	}

	d.input.SetWidth(inputWidth)
	d.input.SetHeight(inputHeight)

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(ColorAccent)
	labelStyle := lipgloss.NewStyle().Foreground(ColorTextDim)
	hintStyle := lipgloss.NewStyle().Foreground(ColorComment)

	var b strings.Builder
	b.WriteString(titleStyle.Render("Orchestrator Chat"))
	b.WriteString("\n")
	b.WriteString(labelStyle.Render("Send a message or queue it for later:"))
	b.WriteString("\n\n")
	b.WriteString(d.input.View())
	b.WriteString("\n")
	b.WriteString(hintStyle.Render("Enter = send  •  Ctrl+S = queue  •  Esc = cancel"))

	box := lipgloss.NewStyle().
		Width(boxWidth).
		Height(boxHeight).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorder).
		Render(b.String())

	return lipgloss.Place(d.width, d.height, lipgloss.Center, lipgloss.Center, box)
}
