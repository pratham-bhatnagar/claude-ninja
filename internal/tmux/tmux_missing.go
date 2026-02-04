package tmux

import (
	"fmt"
	"os"
	"os/exec"
)

// ═══════════════════════════════════════════════════════════════════════════
// Missing Functions - Added for compilation
// ═══════════════════════════════════════════════════════════════════════════

// IsInsideTmux returns true if the current process is running inside a tmux session
func IsInsideTmux() bool {
	// Check TMUX environment variable
	if os.Getenv("TMUX") != "" {
		return true
	}
	// Also check for TMUX_PANE
	if os.Getenv("TMUX_PANE") != "" {
		return true
	}
	return false
}

// DisplayPopup opens a tmux popup window with the specified command
// title is the popup title, width and height are in cells
// commandArgs are the command and arguments to run in the popup
func DisplayPopup(title string, width, height int, commandArgs ...string) error {
	args := []string{
		"display-popup",
		"-w", fmt.Sprintf("%d", width),
		"-h", fmt.Sprintf("%d", height),
		"-E", // Close popup when command exits
	}
	
	// Add the command to execute
	args = append(args, commandArgs...)
	
	cmd := exec.Command("tmux", args...)
	return cmd.Run()
}
