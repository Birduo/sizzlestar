package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	fps  = 60
	rows = 10
	cols = 30
)

var (
	// particle style
	pStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("#575BD8"))
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

func purgeParticles(particles []particle) []particle {
	out := []particle{}
	for _, particle := range particles {
		x := int(math.Round(particle.x))
		y := int(math.Round(particle.y))

		if x < cols && y < rows && x >= 0 && y >= 0 {
			out = append(out, particle)
		}
	}

	return out
}

type particle struct {
	x, xVel float64
	y, yVel float64
	ch      rune
}

// defining the model
type model struct {
	particles []particle
}

// setting init
func (m model) Init() tea.Cmd {
	return tea.Sequence(wait(time.Second/2), animate())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "p":
			p := particle{
				math.Round(cols / 2),
				20 * (rand.Float64() - .5) / fps,
				math.Round(rows / 2),
				20 * (rand.Float64() - .5) / fps,
				'p',
			}

			m.particles = append(m.particles, p)
			return m, nil
		case "r":
			p := particle{
				math.Round(cols / 2),
				20 * (rand.Float64() - .5) / fps,
				math.Round(rows / 2),
				20 * (rand.Float64() - .5) / fps,
				rune(rand.Intn(16777216)), // 16777216 = 256^3
			}

			m.particles = append(m.particles, p)
			return m, nil
		default:
			return m, nil
		}

	case frameMsg:
		// incrementing pos by associated vel
		for index, particle := range m.particles {
			m.particles[index].x += particle.xVel
			m.particles[index].y += particle.yVel
		}

		// removing particles out of range
		m.particles = purgeParticles(m.particles)

		return m, animate()
	default:
		return m, nil
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
			case 'p': // particle to render according to pStyle
				fmt.Fprint(&out, pStyle.Render(" "))
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

// view: this takes model state and renders it to the console
func (m model) View() string {
	var out strings.Builder

	// grid init
	grid := make([][]rune, rows)
	for i := range grid {
		grid[i] = make([]rune, cols)
	}

	// particle render
	for _, particle := range m.particles {
		x := int(math.Round(particle.x))
		y := int(math.Round(particle.y))

		grid[y][x] = particle.ch
	}

	// printing the color-enabled grid
	fmt.Fprint(&out, displayGrid(grid))

	fmt.Fprint(&out, "Particles:", len(m.particles))

	return out.String()
}

func main() {
	m := model{}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
