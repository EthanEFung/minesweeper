package main

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type mainMenu struct {
	model  *model
	cursor int
}

func NewMainMenu(m *model) *mainMenu {
	return &mainMenu{m, 0}
}

func (m *mainMenu) update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m.model, tea.Quit
		case "j":
			if m.cursor >= 2 {
				break
			}
			m.cursor += 1
		case "k":
			if m.cursor <= 0 {
				break
			}
			m.cursor -= 1
		case "enter":
			switch m.cursor {
			case 0:
				m.model.current = m.model.playMenu
			case 1:
				m.model.current = m.model.instructions
			case 2:
				m.model.current = m.model.scores
			}
		}
	}
	return m.model, nil
}

func (m *mainMenu) view() string {
	builder := strings.Builder{}
	builder.WriteString("Welcome to Vim-Minesweeper\n")
	builder.WriteString("Press j to move the cursor down, and k to move the cursor up.\n\n")

	{ // option 1: Play menu
		builder.WriteString("[")
		if m.cursor == 0 {
			builder.WriteString(">")
		} else {
			builder.WriteString(" ")
		}
		builder.WriteString("]")
		builder.WriteString(" Play\n")
	}

	{ // option 2: How to play
		builder.WriteString("[")
		if m.cursor == 1 {
			builder.WriteString(">")
		} else {
			builder.WriteString(" ")
		}
		builder.WriteString("]")
		builder.WriteString(" How to play\n")
	}

	{ // option 3: scores
		builder.WriteString("[")
		if m.cursor == 2 {
			builder.WriteString(">")
		} else {
			builder.WriteString(" ")
		}
		builder.WriteString("]")
		builder.WriteString(" Scores\n\n")
	}

	builder.WriteString("Press 'enter' to select\n")
	return builder.String()
}
