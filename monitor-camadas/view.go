package main

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const (
	minColumns = 85
	minRows    = 24
	handGap    = "     "
	waiting    = "Aguardando o teclado Dilemma V3."
	legendHead = "LEGENDA  /  keymap vendor"
	pairs      = "a·b: toque = a; segurar = b. Tecla isolada: ação direta."
	inherit    = "▽ transparente: herda de uma camada inferior.  × sem ação."
)

var (
	backdrop = lipgloss.Color("#0D1023")
	ink      = lipgloss.Color("#EFF3FF")
	muted    = style("#9CA9D5")
	stroke   = style("#675995")
	heading  = style("#E8D2FF").Bold(true)
	live     = style("#5FFFD1").Bold(true)
	idle     = style("#FFD18A").Bold(true)
	left     = style("#8CEFFF")
	right    = style("#FFA6EC")
	legend   = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#8F72CD")).Padding(0, 2).Foreground(ink)
)

func style(hex string) lipgloss.Style { return lipgloss.NewStyle().Foreground(lipgloss.Color(hex)) }

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
	v.WindowTitle = "Dilemma V3 · Monitor de camadas"
	v.BackgroundColor, v.ForegroundColor = backdrop, ink
	return v
}

func (m *model) render() string {
	if m.width < minColumns || m.height < minRows {
		size := fmt.Sprintf("%d×%d", minColumns, minRows)
		if m.width < 18 {
			return size
		}
		return "Amplie para " + size + "."
	}

	keyWidth := min((m.width-9)/10-1, 11)
	handWidth := 5*(keyWidth+1) + 1
	panelWidth := 2*handWidth + len(handGap)
	card := legend.Width(panelWidth)
	center := func(text string) string {
		return lipgloss.PlaceHorizontal(panelWidth, lipgloss.Center, text)
	}
	title := center(neonBanner) + "\n" + center(muted.Render("MONITOR DE CAMADAS  //  PROCYON"))
	var upper, lower string
	if m.layer < 0 {
		upper = title + "\n\n" + center(idle.Render("○  Sem camada confirmada")) +
			"\n" + center(lipgloss.NewStyle().Width(panelWidth).Align(lipgloss.Center).Render(m.status))
		lower = card.Render(
			heading.Render("TELEMETRIA USB") + "\nO mapa aparece quando o firmware informa a camada ativa.\n" +
				muted.Render("q: sair"))
	} else {
		layer := &layerMaps[m.layer]
		name := heading.Render(fmt.Sprintf("Camada %d · %s", m.layer, layer.name))
		upper = title + "\n" + center(live.Render("● USB")+"   "+name+"   "+muted.Render("q: sair")) + "\n" +
			lipgloss.PlaceHorizontal(handWidth, lipgloss.Center, left.Render("01  ESQUERDA")) + handGap +
			lipgloss.PlaceHorizontal(handWidth, lipgloss.Center, right.Render("02  DIREITA")) +
			keyboard(layer, keyWidth, m.height >= 32)

		var guide strings.Builder
		guide.WriteString(heading.Render(legendHead))
		guide.WriteByte('\n')
		guide.WriteString(muted.Render(pairs + "\n" + inherit))
		for _, text := range layer.legend {
			guide.WriteByte('\n')
			guide.WriteString(stroke.Render("•") + " " + text)
		}
		lower = card.Render(guide.String())
	}
	upper = lipgloss.PlaceVertical(m.height/2, lipgloss.Center, upper)
	return lipgloss.PlaceHorizontal(m.width, lipgloss.Center, upper+"\n"+lower)
}
