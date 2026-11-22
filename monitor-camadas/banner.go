package main

import (
	_ "embed"
	"strings"
)

//go:embed assets/dilemma.txt
var bannerArt string

var neonBanner = renderBanner(bannerArt)

func renderBanner(art string) string {
	palette := [...]string{"#5AF5FF", "#78BCFF", "#B388FF", "#D48CFF", "#FF79C6"}
	lines := strings.Split(strings.TrimRight(art, "\n"), "\n")
	width := 0
	for _, line := range lines {
		width = max(width, len(line)) // O asset contém apenas ASCII.
	}
	for i, line := range lines {
		lines[i] = style(palette[min(i, len(palette)-1)]).Render(line) + strings.Repeat(" ", width-len(line))
	}
	return strings.Join(lines, "\n")
}
