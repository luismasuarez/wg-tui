package ui

import (
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/luismasuarez/wg-tui/internal/nmcli"
	"github.com/luismasuarez/wg-tui/internal/wg"
)

type field struct {
	label    string
	value    string
	required bool
}

// FormModel holds the state for the add-peer form.
type FormModel struct {
	fields  []field
	focused int
	pubKey  string // generated public key shown after creation
}

func newForm() FormModel {
	return FormModel{
		fields: []field{
			{label: "Nombre", required: true},
			{label: "Endpoint (host:port)", required: true},
			{label: "Clave pública del servidor", required: true},
			{label: "IP asignada (CIDR)", required: true},
			{label: "Preshared key (opcional)", required: false},
			{label: "DNS (opcional)", required: false},
		},
	}
}

func (f *FormModel) next() {
	if f.focused < len(f.fields)-1 {
		f.focused++
	}
}

func (f *FormModel) prev() {
	if f.focused > 0 {
		f.focused--
	}
}

func (f *FormModel) input(key string) {
	switch key {
	case "backspace":
		v := f.fields[f.focused].value
		if len(v) > 0 {
			f.fields[f.focused].value = v[:len(v)-1]
		}
	default:
		if len(key) == 1 {
			f.fields[f.focused].value += key
		}
	}
}

func (f *FormModel) validate() error {
	for _, field := range f.fields {
		if field.required && strings.TrimSpace(field.value) == "" {
			return fmt.Errorf("campo requerido: %s", field.label)
		}
	}
	return nil
}

func (f FormModel) render() string {
	var sb strings.Builder
	for i, field := range f.fields {
		prefix := "  "
		label := field.label
		if i == f.focused {
			prefix = "▶ "
			label = styleCursor.Render(label)
		}
		value := field.value
		if i == f.focused {
			value += "█"
		}
		if value == "" && i != f.focused {
			value = strings.Repeat("─", 20)
		}
		sb.WriteString(fmt.Sprintf("%s%-30s %s\n", prefix, label+":", value))
	}
	if f.pubKey != "" {
		sb.WriteString("\n  " + styleActive.Render("Tu clave pública: "+f.pubKey) + "\n")
		sb.WriteString("  " + styleDim.Render("(compártela con el servidor)") + "\n")
	}
	return sb.String()
}

func doAdd(f FormModel) tea.Msg {
	privKey, pubKey, err := wg.GenerateKeypair()
	if err != nil {
		return addResultMsg{errors.New("error generando keypair: " + err.Error())}
	}
	_ = pubKey // shown in status after creation

	name := strings.TrimSpace(f.fields[0].value)
	endpoint := strings.TrimSpace(f.fields[1].value)
	serverPubKey := strings.TrimSpace(f.fields[2].value)
	assignedIP := strings.TrimSpace(f.fields[3].value)
	presharedKey := strings.TrimSpace(f.fields[4].value)
	dns := strings.TrimSpace(f.fields[5].value)

	err = nmcli.AddConnection(name, endpoint, serverPubKey, assignedIP, privKey, presharedKey, dns)
	return addResultMsg{err}
}
