/*
	Purpose: Main game code for Sizzle Star
	- Run core Init, Update, and View for the game
*/

package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	tabs      []tab
	activeTab int
}

// waiting .25s then start animation
func (m model) Init() tea.Cmd {
	return tea.Sequence(wait(time.Second/4), animate())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "up":
			moveSelection(&m.tabs[m.activeTab], -1)
			return m, nil
		case "down":
			moveSelection(&m.tabs[m.activeTab], 1)
			return m, nil
		case "left":
			moveTab(&m, -1)
			return m, nil
		case "right":
			moveTab(&m, 1)
			return m, nil
		default:
			return m, nil
		}

	case frameMsg:
		return m, animate()
	default:
		return m, nil
	}
}

// view: this takes model state and renders it to the console
func (m model) View() string {
	var out strings.Builder

	// rendering tabs into a bar
	tabRow := renderTabRow(m)
	fmt.Fprint(&out, tabRow)
	fmt.Fprint(&out, "\n")

	curTab := m.tabs[m.activeTab]
	curSel := curTab.Upgrades[curTab.selection]
	fmt.Fprint(&out, windowStyle.Width((lipgloss.Width(tabRow) - windowStyle.GetHorizontalFrameSize())).Render(curSel.Description))

	// printing the color-enabled grid
	// fmt.Fprint(&out, borderStyle.Render(displayGrid(grid)))

	fmt.Fprint(&out, "\nPress q to quit\n")
	return out.String()
}

func main() {
	m := loadBaseModel()

	// grid init
	for i := range grid {
		grid[i] = make([]rune, settings.Cols)
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
