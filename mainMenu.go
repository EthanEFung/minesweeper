package main

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type mainMenu struct {
	model   *model
	cursor  int
	list    list.Model
	message tea.Msg
}

var (
	mainMenuStyle     = lipgloss.NewStyle().Margin(1, 2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(4).Foreground(lipgloss.Color("170"))
)

type item string

func (i item) FilterValue() string { return string(i) }

type delegate struct{}

func (d delegate) Height() int                               { return 1 }
func (d delegate) Spacing() int                              { return 0 }
func (d delegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d delegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("  %s", i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return selectedItemStyle.Render("> " + s[2:])
		}
	}

	fmt.Fprint(w, fn(str))
}

func NewMainMenu(m *model) *mainMenu {
	items := []list.Item{
		item("Play"),
		item("How to play"),
		item("Scores"),
	}
	list := list.New(items, delegate{}, 20, 14)
	list.SetShowStatusBar(false)
	list.SetFilteringEnabled(false)
	list.Title = "Welcome to Vim-Minesweeper"
	return &mainMenu{m, 0, list, nil}
}

func (m *mainMenu) update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m.model, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m.model, tea.Quit
		case "enter":
			switch m.list.SelectedItem().FilterValue() {
			case "Play":
				m.model.current = m.model.playMenu
			case "How to play":
				m.model.current = m.model.instructions
			case "Scores":
				m.model.current = m.model.scores
			}
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m.model, cmd
}

func (m *mainMenu) view() string {
	return mainMenuStyle.Render(m.list.View())
}
