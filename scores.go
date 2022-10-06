package main

import tea "github.com/charmbracelet/bubbletea"

/*
scores ...
*/
type scores struct {
	model *model
}

func NewScores(m *model) *scores {
	return &scores{m}
}

func (s *scores) update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return s.model, tea.Quit
		case "b":
			s.model.current = s.model.mainMenu
		}
	}
	return s.model, nil
}

func (s *scores) view() string {
	return "score has been saved. Press 'b' to exit to the main menu."
}
