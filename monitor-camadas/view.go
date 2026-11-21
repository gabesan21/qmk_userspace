package main

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const (
	minColumns = 64
	minRows    = 18
	waiting    = "Aguardando o teclado Dilemma V3."
	legendHead = "Legenda · mapa fixo do keymap vendor"
	pairs      = "a·b: toque = a; segurar = b. Tecla isolada: ação direta."
	inherit    = "▽ transparente: herda de uma camada inferior.  × sem ação."
)

var heading = lipgloss.NewStyle().Bold(true)

type model struct {
	width, height int
	layer         int
	status        string
	content       string
}

func newModel() *model { return &model{layer: -1, status: waiting} }

func (m *model) Init() tea.Cmd { return nil }

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if msg.Width == m.width && msg.Height == m.height {
			return m, nil
		}
		m.width, m.height = msg.Width, msg.Height
	case layerEvent:
		if msg.Layer < -1 || msg.Layer >= len(layerMaps) {
			return m, nil
		}
		if msg.Layer == -1 && msg.Status == "" {
			msg.Status = waiting
		}
		if msg.Layer == m.layer && msg.Status == m.status {
			return m, nil // O heartbeat não muda o mapa desenhado.
		}
		m.layer, m.status = msg.Layer, msg.Status
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
		return m, nil
	default:
		return m, nil
	}
	m.content = m.render()
	return m, nil
}

func (m *model) View() tea.View {
	v := tea.NewView(m.content)
	v.AltScreen = true
	return v
}

func (m *model) render() string {
	if m.width < minColumns || m.height < minRows {
		if m.width < 18 {
			return "64×18"
		}
		return "Amplie para 64×18."
	}
	if m.layer < 0 {
		return heading.Render("Dilemma V3 · Sem camada confirmada") + "\n\n" +
			lipgloss.NewStyle().Width(m.width).Render(m.status) + "\n\nq: sair"
	}

	layer := &layerMaps[m.layer]
	keyWidth := min((m.width-14)/10, 9) // Dez teclas, doze bordas e vão de duas células.
	halfWidth := 5*(keyWidth+1) + 1
	var screen strings.Builder
	fmt.Fprintf(&screen, "%s   q: sair\n", heading.Render(fmt.Sprintf("Dilemma V3 · Camada %d · %s", m.layer, layer.name)))
	screen.WriteString(lipgloss.PlaceHorizontal(halfWidth, lipgloss.Center, "Esquerda"))
	screen.WriteString("  ")
	screen.WriteString(lipgloss.PlaceHorizontal(halfWidth, lipgloss.Center, "Direita"))
	topRows := 2
	for row := range 4 {
		screen.WriteByte('\n')
		if row > 0 && m.height >= 26 {
			screen.WriteByte('\n')
			topRows++
		}
		start, count := row*10, 5
		if row == 3 {
			count = 3
			screen.WriteString(strings.Repeat(" ", 2*(keyWidth+1)))
		}
		writeHand(&screen, layer.keys[start:start+count], keyWidth)
		screen.WriteString("  ")
		writeHand(&screen, layer.keys[start+count:start+2*count], keyWidth)
		topRows++
	}
	screen.WriteString(strings.Repeat("\n", m.height/2-topRows+1))
	screen.WriteString(heading.Render(legendHead))
	for _, text := range []string{pairs, inherit} {
		screen.WriteByte('\n')
		screen.WriteString(text)
	}
	for _, text := range layer.legend {
		screen.WriteByte('\n')
		screen.WriteString(text)
	}
	return lipgloss.PlaceHorizontal(m.width, lipgloss.Center, screen.String())
}

func writeHand(screen *strings.Builder, keys []string, width int) {
	screen.WriteRune('│')
	for _, key := range keys {
		screen.WriteString(lipgloss.PlaceHorizontal(width, lipgloss.Center, key))
		screen.WriteRune('│')
	}
}
