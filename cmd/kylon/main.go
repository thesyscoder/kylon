package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/thesyscoder/kylon/internal/app"
)

func main() {
	p := tea.NewProgram(app.InitialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Uh oh, an error occurred: %v\n", err)
		os.Exit(1)
	}
}
