package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func renderDetail(m Model) string {
	c := m.detail
	var sb strings.Builder

	sb.WriteString(styleTitle.Render("Detalles: "+c.Name) + "\n\n")

	field := func(label, value string) {
		if value == "" {
			value = styleDim.Render("—")
		}
		sb.WriteString(fmt.Sprintf("  %-20s %s\n", label, value))
	}

	state := styleInactive.Render("inactiva")
	if c.Active {
		state = styleActive.Render("activa")
	}
	field("Estado:", state)
	field("Interfaz:", c.Interface)
	field("IP asignada:", c.AssignedIP)
	field("Endpoint:", c.Endpoint)
	field("Clave pública:", c.PublicKey)
	field("Clave pública peer:", c.PeerPublicKey)

	if c.Active {
		field("RX:", formatBytes(c.RxBytes))
		field("TX:", formatBytes(c.TxBytes))
		if !c.LastHandshake.IsZero() {
			field("Último handshake:", c.LastHandshake.Format("2006-01-02 15:04:05"))
		}
	}

	sb.WriteString("\n" + styleDim.Render("esc volver  q QR"))

	return lipgloss.NewStyle().Padding(1, 2).Render(sb.String())
}

func formatBytes(b int64) string {
	if b == 0 {
		return "—"
	}
	switch {
	case b >= 1<<30:
		return fmt.Sprintf("%.2f GB", float64(b)/(1<<30))
	case b >= 1<<20:
		return fmt.Sprintf("%.2f MB", float64(b)/(1<<20))
	case b >= 1<<10:
		return fmt.Sprintf("%.2f KB", float64(b)/(1<<10))
	default:
		return fmt.Sprintf("%d B", b)
	}
}

func renderConfirm(m Model) string {
	name := ""
	if len(m.conns) > 0 {
		name = m.conns[m.cursor].Name
	}
	var sb strings.Builder
	sb.WriteString(styleTitle.Render("Eliminar conexión") + "\n\n")
	sb.WriteString(fmt.Sprintf("  ¿Eliminar %q?\n\n", name))
	sb.WriteString("  " + styleError.Render("y") + " confirmar   " + styleDim.Render("n/esc cancelar"))
	return lipgloss.NewStyle().Padding(1, 2).Render(sb.String())
}

func renderQR(m Model) string {
	var sb strings.Builder
	sb.WriteString(styleTitle.Render("QR — escanear con WireGuard") + "\n\n")
	qr := generateQR(m.qrText)
	sb.WriteString(qr + "\n")
	sb.WriteString(styleDim.Render("esc volver"))
	return lipgloss.NewStyle().Padding(1, 2).Render(sb.String())
}

func renderForm(m Model) string {
	var sb strings.Builder
	sb.WriteString(styleTitle.Render("Agregar peer WireGuard") + "\n\n")
	sb.WriteString(m.form.render())
	sb.WriteString("\n" + styleDim.Render("tab siguiente  esc cancelar  enter confirmar"))
	return lipgloss.NewStyle().Padding(1, 2).Render(sb.String())
}

// Stub — replaced by actual QR rendering in qr.go
var generateQR func(text string) string = func(text string) string {
	return strings.TrimSpace(text)
}
