package main

import (
	"errors"
	"log"
	"math/rand"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type gameState int64

const (
	playable gameState = iota
	won
	lost
)

func (gs gameState) String() string {
	switch gs {
	case playable:
		return "playable"
	case won:
		return "won"
	case lost:
		return "lost"
	default:
		return "unknown"
	}
}


type state int64

const (
	hidden state = iota
	revealed
	flagged
	focused
)

var baseStyle = lipgloss.NewStyle().Width(3).Height(1).Align(lipgloss.Center)
var hiddenStyle = baseStyle.Copy().Background(lipgloss.Color("#2e2e2e"))
var flaggedStyle = baseStyle.Copy()

var revealedStyles = []lipgloss.Style{
	// first will be explosive style
  baseStyle.Copy().Background(lipgloss.Color("#F00")).Foreground(lipgloss.Color("#000")),
	// next will be 0,
	baseStyle.Copy().Foreground(lipgloss.Color("#000")),
	//blue
	baseStyle.Copy().Foreground(lipgloss.Color("#00F")),
	//green
	baseStyle.Copy().Foreground(lipgloss.Color("#0F0")),
	// red
	baseStyle.Copy().Foreground(lipgloss.Color("#F00")),
	// purple
	baseStyle.Copy().Foreground(lipgloss.Color("#800080")),
	// black // since this is meant for dark mode, render white instead
	baseStyle.Copy(),
	// gray
	baseStyle.Copy().Foreground(lipgloss.Color("#808080")),
	// maroon
	baseStyle.Copy().Foreground(lipgloss.Color("#80000")),
	// turquoise
	baseStyle.Copy().Foreground(lipgloss.Color("#0FF")),
}

func createFocusedStyle(style lipgloss.Style) lipgloss.Style {
	return style.Copy().Background(lipgloss.Color("#696969"))
}

func (s state) view(val int, focused bool) string {
	switch s {
	case hidden:
		style := hiddenStyle
		if focused {
			style = createFocusedStyle(style)
		}
		return style.Render(" ")
	case flagged:
		style := flaggedStyle
		if focused {
			style = createFocusedStyle(style)
		}
		return style.Render("🚩")
	}
	var style lipgloss.Style = revealedStyles[val+1]
	if val == -1 {
		return style.Render("*")
	}
	if focused {
		style = createFocusedStyle(style)
	} 
	return style.Render(strconv.Itoa(val))
}

type coord struct {
	x, y int
}

func (c coord) unwrap() (x, y int) {
	return c.x, c.y
}

type model struct {
	// in each cell if it is a mine the int will be -1
	// otherwise each will be the number of mines next to it
	grid      [][]int
	states    [][]state
	cursor    coord
	gameState gameState
}

func (m model) Init() tea.Cmd {
	return nil
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	x, y := m.cursor.unwrap()
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "h":
			if x <= 0 {
				break
			}
			m.cursor = coord{x - 1, y}
		case "j":
			if y >= len(m.grid)-1 {
				break
			}
			m.cursor = coord{x, y + 1}
		case "k":
			if y <= 0 {
				break
			}
			m.cursor = coord{x, y - 1}
		case "l":
			if x >= len(m.grid[0])-1 {
				break
			}
			m.cursor = coord{x + 1, y}
		case " ":
			if m.states[y][x] == revealed {
				break
			}
			show(m.grid, m.states, x, y)
			m.gameState = evaluate(m.grid, m.states)
		case "m":
			if m.states[y][x] == revealed {
				break
			}
			if m.states[y][x] == flagged {
				m.states[y][x] = hidden
				break
			}
			m.states[y][x] = flagged
		case "r":
			return initializeModel(9, 9), nil
		}
	}
	return m, nil
}

func (m model) View() string {
	var s string
	s += "game status: " + m.gameState.String() + "\n"
	for y := range m.states {
		for x := range m.states[y] {
			state := m.states[y][x]
			coordinate := coord{x, y}
			s += state.view(m.grid[y][x], coordinate == m.cursor)
		}
		s += "\n\n"
	}	
	return s
}
	
func show(grid [][]int, states [][]state, initX, initY int) {
	// we'll go breadth first
	queue := []coord{{initX, initY}}
	nx, ny := len(grid[0]), len(grid)

	for len(queue) > 0 {
		// lets go level by level
		currentLen := len(queue)

		for i := 0; i < currentLen; i++ {
			x, y := queue[0].unwrap()
			queue = queue[1:]
			// here we see whether or not the value is 0
			// if so this means that we can show all of
			// the items around it

			states[y][x] = revealed
			
			if grid[y][x] != 0 {
				// if it isn't 0 then move on to the next in the queue
				continue
			}

			adjacent := []coord{
				{x - 1, y - 1}, {x - 1, y}, {x - 1, y + 1},
				{x, y - 1}, {x, y + 1},
				{x + 1, y - 1}, {x + 1, y}, {x + 1, y + 1},
			}

			for _, mod := range adjacent {
				mx, my := mod.unwrap()
				if my < 0 || my > ny-1 {
					continue
				}
				if mx < 0 || mx > nx-1 {
					continue
				}
				if states[my][mx] == hidden {
					queue = append(queue, coord{mx, my})
				}
			}
		}
	}
}

func evaluate(grid[][]int, states[][]state) gameState {
	var hasHidden bool
	for y := range states {
		for x := range states[y] {
			state, val := states[y][x], grid[y][x]
			if val == -1 && state == revealed {
				return lost
			}
			if val != -1 && state == hidden {
				hasHidden = true
			}
		}
	}
	if hasHidden {
		return playable
	}
	return won
}

func placeMines(grid [][]int, n int) error {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	if len(grid) <= 0 {
		return errors.New("cannot place mines on a non existant grid")
	}
	nRows, nCols := len(grid), len(grid[0])
	coords := []coord{}
	for y := range grid {
		for x := range grid[y] {
			coords = append(coords, coord{y, x})
		}
	}

	if n > len(coords) {
		return errors.New("cannot contain more mines than cells in grid")
	}

	r.Shuffle(len(coords), func(i, j int) {
		coords[i], coords[j] = coords[j], coords[i]
	})

	for i := 0; i < n; i++ {
		x, y := coords[i].unwrap()

		grid[y][x] = -1

		adjacent := []coord{
			{x - 1, y - 1}, {x - 1, y}, {x - 1, y + 1},
			{x, y - 1}, {x, y + 1},
			{x + 1, y - 1}, {x + 1, y}, {x + 1, y + 1},
		}

		for _, mod := range adjacent {
			mx, my := mod.unwrap()
			if mx < 0 || mx > nCols-1 {
				continue
			}
			if my < 0 || my > nRows-1 {
				continue
			}
			if grid[my][mx] == -1 {
				continue
			}
			grid[my][mx] += 1
		}
	}

	return nil
}

func initializeModel(width, height int) model {
	grid := make([][]int, height)

	for y := 0; y < height; y++ {
		grid[y] = make([]int, width)
		for x := 0; x < width; x++ {
			grid[y][x] = 0
		}
	}

	placeMines(grid, 10)

	states := make([][]state, height)

	for y := 0; y < height; y++ {
		states[y] = make([]state, width)
		for x := 0; x < width; x++ {
			states[y][x] = hidden
		}
	}

	return model{
		grid:      grid,
		states:    states,
		cursor:    coord{},
		gameState: playable,
	}
}

func main() {
	program := tea.NewProgram(initializeModel(9, 9))
	if err := program.Start(); err != nil {
		log.Fatalf("Booting Error: %v\n", err.Error())
	}
}
