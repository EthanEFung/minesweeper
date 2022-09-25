package main

import (
	"errors"
	"math/rand"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type gameState int64

const (
	pendingGame gameState = iota
	playableGame
	wonGame
	lostGame
)

func (gs gameState) String() string {
	switch gs {
	case playableGame:
		return "playable"
	case wonGame:
		return "won"
	case lostGame:
		return "lost"
	default:
		return "unknown"
	}
}

type gameMode int64

const (
	beginner gameMode = iota
	intermediate
	expert
	custom
)

func (gm gameMode) String() string {
	switch gm {
	case beginner:
		return "beginner"
	case intermediate:
		return "intermediate"
	case expert:
		return "expert"
	case custom:
		return "custom"
	default:
		return "unknown"
	}
}

type cellState int64

const (
	hidden cellState = iota
	revealed
	flagged
	focused
)

type coord struct {
	x, y int
}

func (c coord) unwrap() (int, int) { return c.x, c.y }

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
func (s cellState) view(val int, focused bool) string {
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

func evaluate(grid [][]int, states [][]cellState) gameState {
	var hasHidden bool
	for y := range states {
		for x := range states[y] {
			state, val := states[y][x], grid[y][x]
			if val == -1 && state == revealed {
				return lostGame
			}
			if val != -1 && state == hidden {
				hasHidden = true
			}
		}
	}
	if hasHidden {
		return playableGame
	}
	return wonGame
}

func placeMines(grid [][]int, n int) error {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	if len(grid) <= 0 {
		return errors.New("cannot place mines on a non-existent grid")
	}
	nRows, nCols := len(grid), len(grid[0])
	coords := []coord{}
	for y := range grid {
		for x := range grid[y] {
			coords = append(coords, coord{x, y})
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

func show(grid [][]int, states [][]cellState, initX, initY int) {
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

type game struct {
	model *model
	// in each cell if it is a mine the int will be -1
	// otherwise each will be the number of mines next to it
	grid       [][]int
	cellStates [][]cellState
	cursor     coord
	gameState  gameState
	mode       gameMode
	stopwatch  stopwatch.Model
	flags int
}

func NewGame(model *model) *game {
	return &game{
		model:      model,
		grid:       make([][]int, 0),
		cellStates: make([][]cellState, 0),
		cursor:     coord{},
		gameState:  pendingGame,
		stopwatch:  stopwatch.New(),
		flags: 0,
	}
}

func (g *game) update(msg tea.Msg) (tea.Model, tea.Cmd) {
	x, y := g.cursor.unwrap()
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return g.model, tea.Quit
		case "h":
			if x <= 0 {
				break
			}
			g.cursor = coord{x - 1, y}
		case "j":
			if y >= len(g.grid)-1 {
				break
			}
			g.cursor = coord{x, y + 1}
		case "k":
			if y <= 0 {
				break
			}
			g.cursor = coord{x, y - 1}
		case "l":
			if x >= len(g.grid[0])-1 {
				break
			}
			g.cursor = coord{x + 1, y}
		case "x":
			if g.cellStates[y][x] == revealed {
				break
			}
			show(g.grid, g.cellStates, x, y)
			g.gameState = evaluate(g.grid, g.cellStates)
			if g.gameState == wonGame || g.gameState == lostGame {
				return g.model, g.stopwatch.Stop()
			}
		case "f":
			if g.cellStates[y][x] == revealed {
				break
			}
			if g.cellStates[y][x] == flagged {
				g.cellStates[y][x] = hidden
				g.flags += 1
				break
			}
			g.cellStates[y][x] = flagged
			g.flags -= 1
		case "r":
			if g.mode == beginner {
				g.setBeginner()
			} else if g.mode == intermediate {
				g.setIntermediate()
			} else if g.mode == expert {
				g.setExpert()
			}
			if g.stopwatch.Running() {
				return g.model, g.stopwatch.Reset()
			}
			return g.model, tea.Batch(g.stopwatch.Reset(), g.stopwatch.Start())
		}
	}
	var cmd tea.Cmd
	g.stopwatch, cmd = g.stopwatch.Update(msg)
	return g.model, cmd
}

func (g *game) view() string {
	var s string
	s += "game status: " + g.gameState.String() + "\n"
	s += "time: " + g.stopwatch.View() + "\n"
	s += "mines: " + strconv.Itoa(g.flags) + "\n"
	for y := range g.cellStates {
		for x := range g.cellStates[y] {
			state := g.cellStates[y][x]
			coordinate := coord{x, y}
			s += state.view(g.grid[y][x], coordinate == g.cursor)
		}
		s += "\n\n"
	}
	return s
}

func (g *game) setGrid(width, height, mines int) {
	g.grid = make([][]int, height)

	for y := 0; y < height; y++ {
		g.grid[y] = make([]int, width)
		for x := 0; x < width; x++ {
			g.grid[y][x] = 0
		}
	}

	placeMines(g.grid, mines)

	states := make([][]cellState, height)

	for y := 0; y < height; y++ {
		states[y] = make([]cellState, width)
		for x := 0; x < width; x++ {
			states[y][x] = hidden
		}
	}
	g.cellStates = states
	g.gameState = playableGame
	g.flags = mines // the same number of flags as mines
}

func (g *game) setMode(mode gameMode) {
	g.mode = mode
	if mode == beginner {
		g.setBeginner()
	} else if mode == intermediate {
		g.setIntermediate()
	} else if mode == expert {
		g.setExpert()
	}
}

func (g *game) setBeginner() {
	g.setGrid(9, 9, 10)
}

func (g *game) setIntermediate() {
	g.setGrid(16, 16, 40)
}

func (g *game) setExpert() {
	g.setGrid(30, 16, 99)
}

// func (g *game) setCustom() {}
