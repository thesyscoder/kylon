package app

import tea "github.com/charmbracelet/bubbletea"

// Model is the main application state.
type Model struct{}

// InitialModel returns a new, initialized Model.
func InitialModel() Model {
	return Model{}
}

// Init is a Bubble Tea lifecycle method that initializes the application.
// For now, it doesn't need to do anything.
func (m Model) Init() tea.Cmd {
	return nil
}
