package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/luismasuarez/wg-tui/internal/nmcli"
	"github.com/luismasuarez/wg-tui/ui"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Println("wg-tui", version)
		return
	}

	if err := nmcli.CheckPrerequisites(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	conns, err := nmcli.ListConnections()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error listando conexiones:", err)
		os.Exit(1)
	}

	p := tea.NewProgram(ui.New(conns), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
