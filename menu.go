package main

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type menu struct {
	model   *model
	cursor  int
	modes   []gameMode
}

func NewMenu(model *model) *menu {
	modes := []gameMode{beginner, intermediate, expert}
	return &menu{model, 0, modes}
}

func (m *menu) update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m.model, tea.Quit
		case "j":
			if m.cursor >= len(m.modes)-1 {
				break
			}
			m.cursor += 1
		case "k":
			if m.cursor <= 0 {
				break
			}
			m.cursor -= 1
		case "enter":
			m.model.game.setMode(m.modes[m.cursor])
			m.model.current = m.model.game
		}
	}
	return m.model, nil
}

func (m *menu) view() string {
	b := strings.Builder{}
	for i, mode := range m.modes {
		if i == m.cursor {
			b.WriteString("[>] ")
		} else {
			b.WriteString("[ ] ")
		}
		b.WriteString(mode.String())
		b.WriteRune('\n')
	}
	return b.String()
}
