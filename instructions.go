package main

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type instructions struct {
	model *model
}

func NewInstructions(m *model) *instructions {
	return &instructions{m}
}

func (m *instructions) update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m.model, tea.Quit
		case "b":
			m.model.current = m.model.mainMenu
		}
	}
	return m.model, nil
}

func (i *instructions) view() string {
	b := strings.Builder{}
	b.WriteString("This is vim minesweeper, a simple terminal game to practice vim operations.\n")
	b.WriteString("The goal of minesweeper is to reveal all the cells in the grid that do not have mines.\n\n")

	b.WriteString("The first thing you'll notice in game is there are no arrow keys to navigate the cursor!\nThis is an intentional choice.\n")
	b.WriteString("Instead focus on moving the cursor using 'h','j','k', and 'l'.\n\n")

	b.WriteString("Press 'x' and mimic removing a character to select and reveal a cell.\n")
	b.WriteString("Press 'd' on a revealed number to select and reveal all non-flagged adjacent cells\n          (mimicking deleting a word).\n")
	b.WriteString("Press 'q' at any point (in game or not) to terminate the program.\n\n")

	b.WriteString("Here are some commands you can issue that does not mimic vim.\n\n")
	b.WriteString("You can toggle flags on unrevealed cells by pressing 'f'.\n")
	b.WriteString("At any point during play you can press 'r' to reset the game.\n\n")

	b.WriteString("Now...press 'b', to go to the main menu and get sweeping!\n")

	return b.String()
}
