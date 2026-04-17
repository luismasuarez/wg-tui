package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/luismasuarez/wg-tui/internal/nmcli"
)

type viewState int

const (
	viewList viewState = iota
	viewDetail
	viewQR
	viewAddForm
	viewConfirmDelete
	viewHelp
)

type tickMsg time.Time

type connectResultMsg struct{ err error }
type disconnectResultMsg struct{ err error }
type addResultMsg struct{ err error }
type deleteResultMsg struct{ err error }
type refreshMsg struct {
	conns []nmcli.Connection
	err   error
}
type detailMsg struct {
	conn nmcli.Connection
	err  error
}

var spinner = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// Model is the root Bubble Tea model.
type Model struct {
	conns       []nmcli.Connection
	cursor      int
	view        viewState
	detail      nmcli.Connection
	status      string
	statusOK    bool
	busy        bool
	spinnerIdx  int
	width       int
	height      int
	form        FormModel
	showHelp    bool
	qrText      string
}

var (
	styleActive   = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
	styleInactive = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	styleCursor   = lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true)
	styleStatus   = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	styleError    = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	styleTitle    = lipgloss.NewStyle().Bold(true).Underline(true)
	styleDim      = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)

func New(conns []nmcli.Connection) Model {
	return Model{conns: conns, statusOK: true}
}

func (m Model) Init() tea.Cmd {
	return tick()
}

func tick() tea.Cmd {
	return tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func spinTick() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		if m.width < 80 || m.height < 24 {
			m.status = fmt.Sprintf("Terminal muy pequeño (%dx%d), mínimo 80x24", m.width, m.height)
			m.statusOK = false
		}
		return m, nil

	case tickMsg:
		if m.busy {
			m.spinnerIdx = (m.spinnerIdx + 1) % len(spinner)
			return m, spinTick()
		}
		return m, tea.Batch(tick(), doRefresh())

	case refreshMsg:
		if msg.err == nil {
			m.conns = msg.conns
			if m.cursor >= len(m.conns) && m.cursor > 0 {
				m.cursor = len(m.conns) - 1
			}
		}
		return m, nil

	case connectResultMsg:
		m.busy = false
		if msg.err != nil {
			m.status = "Error: " + msg.err.Error()
			m.statusOK = false
		} else {
			m.status = "Conectado"
			m.statusOK = true
		}
		return m, doRefresh()

	case disconnectResultMsg:
		m.busy = false
		if msg.err != nil {
			m.status = "Error: " + msg.err.Error()
			m.statusOK = false
		} else {
			m.status = "Desconectado"
			m.statusOK = true
		}
		return m, doRefresh()

	case detailMsg:
		if msg.err != nil {
			m.status = "Error: " + msg.err.Error()
			m.statusOK = false
			m.view = viewList
		} else {
			m.detail = msg.conn
			m.view = viewDetail
		}
		return m, nil

	case addResultMsg:
		m.view = viewList
		if msg.err != nil {
			m.status = "Error al agregar: " + msg.err.Error()
			m.statusOK = false
		} else {
			m.status = "Peer agregado"
			m.statusOK = true
		}
		return m, doRefresh()

	case deleteResultMsg:
		m.view = viewList
		if msg.err != nil {
			m.status = "Error al eliminar: " + msg.err.Error()
			m.statusOK = false
		} else {
			m.status = "Conexión eliminada"
			m.statusOK = true
		}
		return m, doRefresh()

	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	// Global
	if key == "ctrl+c" {
		return m, tea.Quit
	}

	switch m.view {
	case viewList:
		return m.handleListKey(key)
	case viewDetail:
		if key == "esc" || key == "q" {
			m.view = viewList
		} else if key == "q" {
			m.view = viewQR
			m.qrText = buildQRText(m.detail)
		}
		return m, nil
	case viewQR:
		if key == "esc" {
			m.view = viewList
		}
		return m, nil
	case viewAddForm:
		return m.handleFormKey(key)
	case viewConfirmDelete:
		return m.handleConfirmKey(key)
	case viewHelp:
		m.view = viewList
		return m, nil
	}
	return m, nil
}

func (m Model) handleListKey(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.conns)-1 {
			m.cursor++
		}
	case "enter":
		if len(m.conns) > 0 && !m.conns[m.cursor].Active {
			m.status = "Conectando…"
			m.statusOK = true
			m.busy = true
			name := m.conns[m.cursor].Name
			return m, tea.Batch(spinTick(), func() tea.Msg {
				return connectResultMsg{nmcli.Connect(name)}
			})
		}
	case "d":
		if len(m.conns) > 0 && m.conns[m.cursor].Active {
			m.status = "Desconectando…"
			m.statusOK = true
			m.busy = true
			name := m.conns[m.cursor].Name
			return m, tea.Batch(spinTick(), func() tea.Msg {
				return disconnectResultMsg{nmcli.Disconnect(name)}
			})
		}
	case "r":
		return m, doRefresh()
	case "i":
		if len(m.conns) > 0 {
			name := m.conns[m.cursor].Name
			return m, func() tea.Msg {
				c, err := nmcli.GetDetails(name)
				return detailMsg{c, err}
			}
		}
	case "q":
		if len(m.conns) > 0 {
			c := m.conns[m.cursor]
			m.qrText = buildQRText(c)
			m.view = viewQR
		}
	case "a":
		m.form = newForm()
		m.view = viewAddForm
	case "x":
		if len(m.conns) > 0 {
			m.view = viewConfirmDelete
		}
	case "?":
		m.showHelp = !m.showHelp
	}
	return m, nil
}

func (m Model) handleFormKey(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "esc":
		m.view = viewList
	case "tab", "down":
		m.form.next()
	case "shift+tab", "up":
		m.form.prev()
	case "enter":
		if m.form.focused == len(m.form.fields)-1 {
			if err := m.form.validate(); err != nil {
				m.status = err.Error()
				m.statusOK = false
				return m, nil
			}
			f := m.form
			return m, func() tea.Msg {
				return doAdd(f)
			}
		}
		m.form.next()
	default:
		m.form.input(key)
	}
	return m, nil
}

func (m Model) handleConfirmKey(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "y":
		name := m.conns[m.cursor].Name
		active := m.conns[m.cursor].Active
		return m, func() tea.Msg {
			if active {
				_ = nmcli.Disconnect(name)
			}
			return deleteResultMsg{nmcli.DeleteConnection(name)}
		}
	case "n", "esc":
		m.view = viewList
	}
	return m, nil
}

func (m Model) View() string {
	if m.width < 80 || m.height < 24 {
		return styleError.Render(fmt.Sprintf("Terminal muy pequeño (%dx%d), mínimo 80x24", m.width, m.height))
	}
	switch m.view {
	case viewDetail:
		return renderDetail(m)
	case viewQR:
		return renderQR(m)
	case viewAddForm:
		return renderForm(m)
	case viewConfirmDelete:
		return renderConfirm(m)
	default:
		return renderList(m)
	}
}

func doRefresh() tea.Cmd {
	return func() tea.Msg {
		conns, err := nmcli.ListConnections()
		return refreshMsg{conns, err}
	}
}

func buildQRText(c nmcli.Connection) string {
	var sb strings.Builder
	sb.WriteString("[Interface]\n")
	if c.AssignedIP != "" {
		sb.WriteString("Address = " + c.AssignedIP + "\n")
	}
	if c.PublicKey != "" {
		sb.WriteString("# PublicKey = " + c.PublicKey + "\n")
	}
	sb.WriteString("\n[Peer]\n")
	if c.PeerPublicKey != "" {
		sb.WriteString("PublicKey = " + c.PeerPublicKey + "\n")
	}
	if c.Endpoint != "" {
		sb.WriteString("Endpoint = " + c.Endpoint + "\n")
	}
	sb.WriteString("AllowedIPs = 0.0.0.0/0\n")
	return sb.String()
}
