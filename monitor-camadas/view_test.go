package main

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func TestVendorPositions(t *testing.T) {
	names := []string{"Base", "Funções", "Navegação", "Mídia", "Ponteiro", "Números", "Símbolos"}
	if len(layerMaps) != len(names) {
		t.Fatalf("catálogo com %d camadas, esperado %d", len(layerMaps), len(names))
	}
	for id, name := range names {
		layer := &layerMaps[id]
		if layer.name != name || len(layer.keys) != 36 {
			t.Fatalf("camada %d: %s, %d posições", id, layer.name, len(layer.keys))
		}
		for i, key := range layer.keys {
			if key == "" || lipgloss.Width(key) > 5 {
				t.Errorf("%s: posição %d ausente ou maior que a tecla compacta: %q", name, i, key)
			}
		}
	}
	for _, tc := range []struct {
		layer, start int
		want         []string
	}{
		{0, 0, []string{"Q·`", "W·~", "E·'", "R", "T", "Y", "U·\"", "I·^", "O", "P"}},
		{0, 10, []string{"A·Alt", "S·Ctl", "D·Sh", "F·Gui", "G", "H", "J·Gui", "K·Sh", "L·Ctl", "'·Alt"}},
		{0, 20, []string{"Z·4", "X·Gr", "C", "V", "B", "N", "M", ",", ".·Gr", "/·4"}},
		{0, 30, []string{"Esc·3", "Tab·5", "Esp·2", "Ent·6", "Bsp·1", "Mudo"}},
		{1, 5, []string{"Tela", "F7", "F8", "F9", "F12"}},
		{1, 33, []string{"×", "▽", "×"}},
		{2, 15, []string{"Caps", "←", "↓", "↑", "→"}},
		{2, 25, []string{"Ins", "Iníc", "Pg↓", "Pg↑", "Fim"}},
		{3, 30, []string{"▽", "Tocar", "Parar", "Parar", "Tocar", "Mudo"}},
		{4, 20, []string{"▽", "Rol", "Mira", "B3", "×", "×", "B3", "Mira", "Rol", "▽"}},
		{4, 30, []string{"B3", "B2", "B1", "B1", "B2", "B3"}},
		{5, 25, []string{"NumLk", "1", "2", "3", "="}},
		{5, 30, []string{"×", "▽", "×", "EntN", "×", "0"}},
		{6, 0, []string{"!", "@", "#", "$", "%", "&", "*", "(", ")", "-·_"}},
		{6, 15, []string{"×", ";·:", "[·{", "]·}", "=·+"}},
		{6, 25, []string{"×", "<", ">", "\\·|", "?"}},
		{6, 30, []string{"×", "×", "Esp", "▽", "×", "×"}},
	} {
		got := layerMaps[tc.layer].keys[tc.start : tc.start+len(tc.want)]
		if !slices.Equal(got, tc.want) {
			t.Errorf("camada %d, posição %d: %v; esperado %v", tc.layer, tc.start, got, tc.want)
		}
	}
}

func TestMapAndLegendFitAfterResize(t *testing.T) {
	for id := range layerMaps {
		t.Run(layerMaps[id].name, func(t *testing.T) {
			m := newModel()
			m.Update(layerEvent{Layer: id})
			for _, size := range []tea.WindowSizeMsg{
				{Width: 85, Height: 24}, {Width: 128, Height: 36},
				{Width: 92, Height: 32}, {Width: 140, Height: 40}, {Width: 85, Height: 24},
			} {
				m.Update(size)
				content := m.View().Content
				w, h := lipgloss.Size(content)
				if w > size.Width || h > size.Height {
					t.Fatalf("tela %dx%d excede %dx%d", w, h, size.Width, size.Height)
				}
				plain := sgr.ReplaceAllString(content, "")
				lines := strings.Split(plain, "\n")
				if len(lines) <= size.Height/2 ||
					!strings.Contains(strings.Join(lines[:size.Height/2], "\n"), fmt.Sprintf("Camada %d · %s", id, layerMaps[id].name)) {
					t.Fatalf("cabeçalho ou área inferior ausente em %dx%d", size.Width, size.Height)
				}
				checkKeyboard(t, lines[:size.Height/2], &layerMaps[id])
				legend := strings.Join(lines[size.Height/2:], "\n")
				for _, text := range append([]string{legendHead, pairs, inherit}, layerMaps[id].legend...) {
					if !strings.Contains(legend, text) {
						t.Errorf("legenda ausente em %dx%d: %q", size.Width, size.Height, text)
					}
				}
			}
		})
	}
}

var sgr = regexp.MustCompile(`\x1b\[[0-9;:]*m`)

func checkKeyboard(t *testing.T, lines []string, layer *layerMap) {
	t.Helper()
	row, pos := 0, 0
	for _, line := range lines {
		if !strings.Contains(line, "│") {
			continue
		}
		parts := strings.Split(strings.Trim(strings.TrimSpace(line), "│"), "│     │")
		if len(parts) != 2 || row >= 4 {
			t.Fatalf("fileira física inválida: %s", line)
		}
		for _, hand := range parts {
			keys := strings.Fields(hand)
			if row < 3 && len(keys) != 5 || row == 3 && len(keys) != 3 {
				t.Fatalf("fileira %d: %d teclas na mão", row, len(keys))
			}
			for _, key := range keys {
				if strings.TrimSpace(key) != layer.keys[pos] {
					t.Errorf("posição %d: %q; esperada %q", pos, key, layer.keys[pos])
				}
				pos++
			}
		}
		row++
	}
	if pos != 36 {
		t.Fatalf("%d posições visíveis de 36", pos)
	}
}

func TestNoLayerUntilConfirmedAfterConnection(t *testing.T) {
	m := newModel()
	m.Update(tea.WindowSizeMsg{Width: 85, Height: 24})
	noMap := func(status string) {
		t.Helper()
		content := m.View().Content
		if !strings.Contains(content, "Sem camada confirmada") || !strings.Contains(content, status) || strings.Contains(content, "01  ESQUERDA") {
			t.Fatalf("esperado diagnóstico sem mapa; recebido %q", content)
		}
	}
	noMap(waiting)
	m.Update(layerEvent{Layer: 4})
	if !strings.Contains(m.View().Content, "Camada 4 · Ponteiro") {
		t.Fatal("evento válido não exibiu Ponteiro")
	}
	m.Update(layerEvent{Layer: -1, Status: "Sem acesso: ajuste a regra udev."})
	noMap("ajuste a regra udev")
	m.Update(layerEvent{Layer: -1, Status: "Aguardando mensagem do teclado."})
	m.Update(layerEvent{Layer: 7})
	noMap("Aguardando mensagem")
	m.Update(layerEvent{Layer: 0})
	if !strings.Contains(m.View().Content, "Camada 0 · Base") {
		t.Fatal("reconexão com mensagem da Base não restaurou o mapa")
	}
}

func TestBelowMinimumAndRecovery(t *testing.T) {
	m := newModel()
	m.Update(layerEvent{Layer: 6})
	for _, size := range []tea.WindowSizeMsg{{Width: 84, Height: 24}, {Width: 85, Height: 23}} {
		m.Update(size)
		text := m.View().Content
		if !strings.Contains(text, "85×24") || strings.Contains(text, "│") {
			t.Errorf("esperado aviso sem mapa cortado, recebido %q", text)
		}
	}
	m.Update(tea.WindowSizeMsg{Width: 85, Height: 24})
	if !strings.Contains(m.View().Content, "Camada 6 · Símbolos") {
		t.Fatal("resize perdeu a camada confirmada")
	}
}
