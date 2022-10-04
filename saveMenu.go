package main

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type saveMenu struct {
	model    *model
	initials []rune
}

func NewSaveMenu(m *model) *saveMenu {
	return &saveMenu{m, []rune{}}
}

func (m *saveMenu) view() string {
	b := strings.Builder{}
	b.WriteString("\n")
	for _, char := range m.initials {
		b.WriteRune(char)
	}
	b.WriteString("\nSave? (y / n)\n")

	return b.String()
}

func (m *saveMenu) update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m.model, tea.Quit
		case "n", "b":
			m.model.game = NewGame(m.model)
			m.model.current = m.model.mainMenu
		case "y", "w":
			save(m.model.game, m.initials)
		}
	}
	return m.model, nil
}

func save(game *game, initials []rune) {
	// ...
}
