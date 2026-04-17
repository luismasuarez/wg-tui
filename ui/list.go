package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func renderList(m Model) string {
	var sb strings.Builder

	sb.WriteString(styleTitle.Render("wg-tui — WireGuard connections") + "\n\n")

	if len(m.conns) == 0 {
		sb.WriteString(styleDim.Render("No WireGuard connections found in NetworkManager.\n"))
		sb.WriteString(styleDim.Render("Press 'a' to add one.\n"))
	} else {
		for i, c := range m.conns {
			cursor := "  "
			if i == m.cursor {
				cursor = styleCursor.Render("▶ ")
			}
			var name string
			if c.Active {
				name = styleActive.Render(fmt.Sprintf("● %s", c.Name))
			} else {
				name = styleInactive.Render(fmt.Sprintf("○ %s", c.Name))
			}
			sb.WriteString(cursor + name + "\n")
		}
	}

	sb.WriteString("\n")
	if m.status != "" {
		st := styleStatus
		if !m.statusOK {
			st = styleError
		}
		prefix := ""
		if m.busy {
			prefix = spinner[m.spinnerIdx] + " "
		}
		sb.WriteString(st.Render(prefix+m.status) + "\n")
	}

	if m.showHelp {
		sb.WriteString("\n" + renderHelp())
	} else {
		sb.WriteString(styleDim.Render("? help  ctrl+c quit") + "\n")
	}

	return lipgloss.NewStyle().Padding(1, 2).Render(sb.String())
}

func renderHelp() string {
	var sb strings.Builder
	for _, k := range helpKeys {
		sb.WriteString(fmt.Sprintf("  %-12s %s\n", k.key, k.desc))
	}
	return styleDim.Render(sb.String())
}
