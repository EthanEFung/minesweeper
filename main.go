package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

type current interface {
	view() string
	update(tea.Msg) (tea.Model, tea.Cmd)
}

type model struct {
	menu    *menu
	game    *game
	current current
}

func NewModel() model {
	m := new(model)
	m.game = NewGame(m)
	m.menu = NewMenu(m)
	m.current = m.menu
	// hmm this seems suspicious
	return *m
}

func (m model) Init() tea.Cmd {
	return m.game.stopwatch.Init()
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.current.update(msg)
}
func (m model) View() string {
	return m.current.view()
}

func main() {
	program := tea.NewProgram(NewModel())
	if err := program.Start(); err != nil {
		log.Fatalf("Booting Error: %v\n", err.Error())
	}
}
