package ui

import "github.com/charmbracelet/bubbletea"

type keyMsg struct{ key string }

func keyHandler(msg tea.KeyMsg) string {
	return msg.String()
}

// Keybinding help entries shown in the footer.
var helpKeys = []struct{ key, desc string }{
	{"↑/k", "arriba"},
	{"↓/j", "abajo"},
	{"enter", "conectar"},
	{"d", "desconectar"},
	{"i", "detalles"},
	{"q", "QR"},
	{"a", "agregar"},
	{"x", "eliminar"},
	{"r", "refrescar"},
	{"?", "ayuda"},
	{"ctrl+c", "salir"},
}
