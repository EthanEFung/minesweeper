package main

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	// "github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

type keymap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Back   key.Binding
	Quit   key.Binding
}

func (k keymap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Select, k.Back, k.Quit}
}

func (k keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.Select, k.Back},
		{k.Quit},
	}
}

var keys = keymap{
	Up: key.NewBinding(
		key.WithKeys("k"),
		key.WithHelp("k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("j"),
		key.WithHelp("j", "move down"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select option"),
	),
	Back: key.NewBinding(
		key.WithKeys("b"),
		key.WithHelp("back", "previous menu"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	),
}

type playMenu struct {
	model  *model
	cursor int
	modes  []gameMode
	keys   keymap
}

func NewPlayMenu(model *model) *playMenu {
	modes := []gameMode{beginner, intermediate, expert}
	return &playMenu{model, 0, modes, keys}
}

func (m *playMenu) update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m.model, tea.Quit
		case key.Matches(msg, m.keys.Down):
			if m.cursor >= len(m.modes)-1 {
				break
			}
			m.cursor += 1
		case key.Matches(msg, m.keys.Up):
			if m.cursor <= 0 {
				break
			}
			m.cursor -= 1
		case key.Matches(msg, m.keys.Back):
			m.model.current = m.model.mainMenu
		case key.Matches(msg, m.keys.Select):
			m.model.game.setMode(m.modes[m.cursor])
			m.model.current = m.model.game
			return m.model, m.model.game.stopwatch.Start()
		}
	}
	return m.model, nil
}

func (m *playMenu) view() string {
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
