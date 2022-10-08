package main

import (
	"encoding/csv"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

/*
scores ...
*/
type scores struct {
	model *model
	table table.Model
}

func NewScores(m *model) *scores {
	return &scores{m, NewTable()}
}

func (s *scores) update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return s.model, tea.Quit
		case "b":
			s.model.current = s.model.mainMenu
		case "j":
			s.table.MoveDown(1)
		case "k":
			s.table.MoveUp(1)
		}
	}
	return s.model, nil
}

func (s *scores) view() string {
	b := strings.Builder{}
	b.WriteString(s.table.View() + "\n")
	b.WriteString("Press 'b' to exit to the main menu.")
	return b.String()
}

/*
	reevaluate will read the csv file again and set the records to the new values
*/
func (s *scores) reevaluate() {
	s.table = NewTable()
}

func NewTable() table.Model {
	columns := []table.Column{
		{Title: "Rank", Width: 5},
		{Title: "Mode", Width: 12},
		{Title: "Time", Width: 10},
		{Title: "Player", Width: 6},
	}
	records, err := readCSV()

	if err != nil {
		log.Fatal(err)
	}

	rows := deriveRows(records)

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(10),
	)

	return t
}

func readCSV() ([][]string, error) {
	file, err := os.OpenFile("scores.csv", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	return reader.ReadAll()
}

func deriveRows(records [][]string) []table.Row {
	rows := []table.Row{}
	for _, record := range records {
		rows = append(rows, table.Row{"", record[3], record[1], record[0]})
	}
	return rows
}
