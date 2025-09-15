package app

import tea "github.com/charmbracelet/bubbletea"

// Update is the core logic for the Bubble Tea program.
// It receives messages (e.g., keypresses) and updates the model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}
