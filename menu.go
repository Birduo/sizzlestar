/*
	Purpose: Host menu functionality and structure
*/

package main

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// constants to be set in a config file later
const (
	fps  = 60
	rows = 10
	cols = 30
)

var (
	borderStyle = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder())
)

type frameMsg time.Time

func animate() tea.Cmd {
	return tea.Tick(time.Second/fps, func(t time.Time) tea.Msg {
		return frameMsg(t)
	})
}

func wait(d time.Duration) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(d)
		return nil
	}
}

// displayGrid: this takes a rune grid and prints a color-enabled grid
func displayGrid(grid [][]rune) string {
	var out strings.Builder

	for _, row := range grid {
		for _, char := range row {
			switch char {
			case 0: // empty cell: render as nothing
				fmt.Fprint(&out, " ")
			default: // render as the given rgb of the last slice of the i32
				hexColor := fmt.Sprintf("#%06X", char)

				// cool debug hex tracker
				// fmt.Fprint(&out, hexColor)

				randStyle := lipgloss.NewStyle().
					Background(lipgloss.Color(hexColor))

				fmt.Fprint(&out, randStyle.Render(" "))
			}
		}
		fmt.Fprint(&out, "\n")
	}

	return out.String()
}
