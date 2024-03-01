/*
	Purpose: Host menu structure and functionality
*/

package main

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	settings = loadConfig()

	// this style just sets a border around the given text
	borderStyle = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder())

	// this is the grid to be drawn by the main tab
	grid = make([][]rune, settings.Rows)
)

// moveSelection(tab, n), moves the inner-tab selection by n amount
func moveSelection(t *tab, n int) {
	t.selection = min(max(0, t.selection+n), len(t.Upgrades)-1)
}

// moveTab(model, n), moves the tab selection by n amount
func moveTab(m *model, n int) {
	// m.activeTab = min(max(0, m.activeTab+n), len(m.tabs)-1)
	m.activeTab = min(max(0, m.activeTab+n), len(m.tabs)-1)
}

// helper functions for selections taken from the bubbletea examples
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// grid animation / display functions
type frameMsg time.Time

func animate() tea.Cmd {
	return tea.Tick(time.Second/time.Duration(settings.Fps), func(t time.Time) tea.Msg {
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
