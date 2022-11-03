package main

import (
	"encoding/csv"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

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

type sortable [][]string

func (records sortable) Len() int      { return len(records) }
func (records sortable) Swap(i, j int) { records[i], records[j] = records[j], records[i] }
func (records sortable) Less(i, j int) bool {
	a, b := records[i], records[j]

	durationA, err := time.ParseDuration(a[1])
	if err != nil {
		log.Fatal(err)
	}
	durationB, err := time.ParseDuration(b[1])
	if err != nil {
		log.Fatal(err)
	}
	if durationA != durationB {
		return durationA < durationB
	}

	playedA, err := time.Parse(time.RFC3339, a[2])
	if err != nil {
		log.Fatal(err)
	}
	playedB, err := time.Parse(time.RFC3339, b[2])
	if err != nil {
		log.Fatal(err)
	}
	return playedA.Before(playedB)

}

func NewScores(m *model) *scores {
	records, err := readCSV()
	if err != nil {
		log.Fatal(err)
	}
	return &scores{m, NewTable(records)}
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
	records, err := readCSV()
	if err != nil {
		log.Fatal(err)
	}
	s.table = NewTable(records)
}

func NewTable(records sortable) table.Model {
	columns := []table.Column{
		{Title: "Rank", Width: 5},
		{Title: "Mode", Width: 12},
		{Title: "Time", Width: 10},
		{Title: "Player", Width: 6},
	}

	sort.Sort(records)

	rows := deriveRows(records)

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(10),
	)

	cursor := latestIndex(records)
	t.SetCursor(cursor)

	return t
}

func readCSV() (sortable, error) {
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
	for i, record := range records {
		rows = append(rows, table.Row{strconv.Itoa(i + 1), record[3], record[1], record[0]})
	}
	return rows
}

func latestIndex(records [][]string) int {
	if len(records) == 0 {
		return 0
	}
	var l int
	for i, record := range records {
		t, err := time.Parse(time.RFC3339, record[2])
		if err != nil {
			log.Fatal(err)
		}
		latest, err := time.Parse(time.RFC3339, records[l][2])
		if err != nil {
			log.Fatal(err)
		}
		if latest.Before(t) {
			l = i
		}
	}
	if l == 0 {
		log.Fatal(records[0][2])
	}
	return l
}
