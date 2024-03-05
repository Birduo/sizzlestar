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

	// creating inactive/active tabs
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	highlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Copy().Border(activeTabBorder, true)
	windowStyle       = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(2, 0).Align(lipgloss.Center).Border(lipgloss.NormalBorder()).UnsetBorderTop()
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

// stealing tab display from https://github.com/charmbracelet/bubbletea/blob/master/examples/tabs
// create rectangular tab border with custom edges
func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

func renderTabRow(m model) string {
	renderedTabs := []string{}
	for i, tab := range m.tabs {
		var style lipgloss.Style

		isFirst, isLast, isActive := i == 0, i == len(m.tabs)-1, i == m.activeTab

		if isActive {
			style = activeTabStyle.Copy()
		} else {
			style = inactiveTabStyle.Copy()
		}

		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "│"
		} else if isLast && !isActive {
			border.BottomRight = "┤"
		}

		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(tab.Name))
	}

	tabRow := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)

	return tabRow
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
