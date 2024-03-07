/*
	Purpose: Host menu structure and functionality
*/

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	settings = loadConfig()

	// maxLWidth is the maximum width of descriptions within the window
	maxRWidth = 7
	maxLWidth = settings.Cols - maxRWidth
	// determined to permit " $1234" for numbers

	// creating inactive/active tabs
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	highlightColor    = lipgloss.AdaptiveColor{Light: "#194D33", Dark: "#2D8659"}
	inactiveTabStyle  = lipgloss.NewStyle().
				Border(inactiveTabBorder, true).
				BorderForeground(highlightColor).
				Padding(0, 1)
	activeTabStyle = inactiveTabStyle.Copy().Border(activeTabBorder, true)
	windowStyle    = lipgloss.NewStyle().
			BorderForeground(highlightColor).
			Width(settings.Cols).
			Height(settings.Rows).
			Border(lipgloss.NormalBorder()).
			UnsetBorderTop()
	helpStyle           = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1).Faint(true)
	selectedItemStyle   = lipgloss.NewStyle().Foreground(highlightColor)
	unselectedItemStyle = lipgloss.NewStyle().Faint(true)
)

// moveTab(model, n), moves the tab selection by n amount
func moveTab(m *model, n int) {
	// m.activeTab = min(max(0, m.activeTab+n), len(m.tabs)-1)
	m.activeTab = min(max(0, m.activeTab+n), len(m.tabs)-1)
}

// moveSelection(tab, n), moves the selection within a tab by n amount
func moveSelection(t *tab, n int) {
	t.selection = min(max(0, t.selection+n), len(t.Upgrades)-1)
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
		renderedTabs = append(renderedTabs, style.Render(tab.Icon))
	}

	tabRow := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)

	return tabRow
}

func renderUpgradeRow(left string, right string, style lipgloss.Style) string {
	leftSide := style.Copy().Width(maxLWidth).Render(left)
	rightSide := style.Copy().Align(lipgloss.Right).Width(maxRWidth).Render(right)
	return lipgloss.JoinHorizontal(lipgloss.Top, leftSide, rightSide)
}

func renderTabContent(m model) string {
	var out strings.Builder

	curTab := m.tabs[m.activeTab]
	// fmt.Fprint(&out, curSel, "\n")
	// to render per upgrade: description, cost, prod, owned
	for i := max(curTab.selection-1, 0); i < len(curTab.Upgrades) && 3*i < settings.Rows; i++ {
		var style lipgloss.Style
		if i == curTab.selection {
			style = selectedItemStyle.Copy()
		} else {
			style = unselectedItemStyle.Copy()
		}

		description := curTab.Upgrades[i].Description
		owned := fmt.Sprintf("x%d", curTab.Upgrades[i].owned)
		fmt.Fprint(&out, renderUpgradeRow(description, owned, style)+"\n")
		cost := fmt.Sprintf("Cost: %.2f", curTab.Upgrades[i].Cost)
		production := fmt.Sprintf("%.2f¥/s", curTab.Upgrades[i].Production)
		fmt.Fprint(&out, renderUpgradeRow(cost, production, style)+"\n")
	}

	return out.String()
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
