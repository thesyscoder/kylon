package app

// View renders the TUI.
func (m Model) View() string {
	s := "Welcome to Kylon, the Kubernetes disaster recovery app!\n\n"
	s += "Press 'q' or 'Ctrl+C' to quit.\n"
	return s
}
