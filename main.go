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
	mainMenu     *mainMenu
	playMenu     *playMenu
	instructions *instructions
	game         *game
	saveMenu     *saveMenu
	current      current
}

func NewModel() model {
	m := new(model)
	m.game = NewGame(m)
	m.playMenu = NewPlayMenu(m)
	m.mainMenu = NewMainMenu(m)
	m.instructions = NewInstructions(m)
	m.saveMenu = NewSaveMenu(m)
	m.current = m.mainMenu
	return *m
}

func (m model) Init() tea.Cmd {
	return nil
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
