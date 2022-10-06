package main

import (
	"encoding/csv"
	"log"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type saveMenu struct {
	model    *model
	initials []rune
	cursor   int
}

var baseFocusedStyle = createFocusedStyle(baseStyle)

func NewSaveMenu(m *model) *saveMenu {
	return &saveMenu{
		model:    m,
		initials: []rune{'A', 'A', 'A'},
		cursor:   0,
	}
}

func (m *saveMenu) view() string {
	b := strings.Builder{}
	b.WriteString("\n")
	b.WriteString("Change the initials using h, j, k, and l. Press y to save.\n")
	b.WriteString("Pressing n will take you to the menu\n\n")
	for i, char := range m.initials {
		str := string(char)
		if i == m.cursor {
			b.WriteString(baseFocusedStyle.Render(str))
			continue
		}
		b.WriteString(baseStyle.Render(str))
	}
	b.WriteString("\n\nSave? (y / n)\n")

	return b.String()
}

func (m *saveMenu) update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m.model, tea.Quit
		case "n":
			m.model.game = NewGame(m.model)
			m.model.current = m.model.mainMenu
		case "y":
			save(m.model.game, m.initials)
			m.model.current = m.model.scores
		case "h":
			if m.cursor > 0 {
				m.cursor--
			}
		case "l":
			if m.cursor < 2 {
				m.cursor++
			}
		case "j":
			if m.initials[m.cursor] < 'Z' {
				m.initials[m.cursor]++
			}
		case "k":
			if m.initials[m.cursor] > 'A' {
				m.initials[m.cursor]--
			}
		}
	}
	return m.model, nil
}

func save(game *game, initials []rune) {
	// ...
	file, err := os.OpenFile("scores.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	err = writer.Write([]string{string(initials), game.stopwatch.Elapsed().String(), time.Now().String()})

	if err != nil {
		log.Fatal(err)
	}
	writer.Flush()
}
