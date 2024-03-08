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
)

type model struct {
	state     gameState
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
			m.saveGameState()
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

	fmt.Fprint(&out, windowStyle.Render(renderTabContent(m)))
	// fmt.Fprint(&out, renderTabContent(m))

	fmt.Fprint(&out, helpStyle.Render("Press q to quit"))
	return out.String()
}

func main() {
	m := loadBaseModel()

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
