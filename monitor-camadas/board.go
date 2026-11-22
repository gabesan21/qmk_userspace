package main

import (
	"strings"

	"charm.land/lipgloss/v2"
)

func keyboard(layer *layerMap, width int, spacious bool) string {
	var board strings.Builder
	cap := lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Background(lipgloss.Color("#202747"))
	dimCap := cap.Foreground(lipgloss.Color("#9BA8CF")).Background(lipgloss.Color("#171C36"))
	holdCap := cap.Foreground(lipgloss.Color("#FFE4FA")).Background(lipgloss.Color("#41305F"))
	handCaps := [...]lipgloss.Style{cap.Foreground(lipgloss.Color("#BDF5FF")), cap.Foreground(lipgloss.Color("#FFD1F4"))}
	edges := [...]string{left.Render("│"), right.Render("│")}
	var regularRim, thumbRim string
	if spacious {
		capRim := "╰" + strings.Repeat("─", width-2) + "╯"
		regularRim = stroke.Render(" " + strings.Repeat(capRim+" ", 4) + capRim)
		thumbRim = stroke.Render(" " + strings.Repeat(capRim+" ", 2) + capRim)
	}
	// Polegares começam duas posições para dentro na mão esquerda.
	for row := range 4 {
		board.WriteByte('\n')
		start, count, indent, rim := row*10, 5, "", regularRim
		if row == 3 {
			count = 3
			indent = strings.Repeat(" ", 2*(width+1))
			rim = thumbRim
		}
		board.WriteString(indent)
		for hand, edge := range edges {
			if hand == 1 {
				board.WriteString(handGap)
			}
			board.WriteString(edge)
			for i, key := range layer.keys[start+hand*count : start+(hand+1)*count] {
				if i > 0 {
					board.WriteByte(' ')
				}
				appearance := handCaps[hand]
				switch {
				case key == transparent || key == inactive:
					appearance = dimCap
				case strings.Contains(key, "·"):
					appearance = holdCap
				}
				board.WriteString(appearance.Render(key))
			}
			board.WriteString(edge)
		}
		if spacious {
			board.WriteByte('\n')
			board.WriteString(indent)
			board.WriteString(rim)
			board.WriteString(" " + handGap)
			board.WriteString(rim)
		}
	}
	return board.String()
}
