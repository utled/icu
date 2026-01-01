package ui

import (
	"database/sql"
	"fmt"
	"snafu/data"
	"snafu/search"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

type styles struct {
	BorderColor lipgloss.Color
	InputField  lipgloss.Style
}

type searchResult struct {
	Rows []data.SearchResult
	Err  error
}

type Model struct {
	width         int
	height        int
	inputField    textinput.Model
	viewPort      viewport.Model
	searchResults []data.SearchResult
	style         styles
	dbConnection  *sql.DB
	err           error
}

func defaultStyles() styles {
	style := new(styles)
	style.BorderColor = lipgloss.Color("36")
	style.InputField = lipgloss.NewStyle().BorderForeground(style.BorderColor).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(80)
	return *style
}

func NewModel() Model {
	inputStyles := defaultStyles()
	textInput := textinput.New()
	textInput.Placeholder = "Enter search term"
	textInput.Width = 30
	textInput.Focus()

	return Model{
		inputField:    textInput,
		style:         inputStyles,
		searchResults: []data.SearchResult{},
	}
}

func limitPathLength(s string, max int) string {
	if len(s) <= max {
		return s
	}
	if max < 3 {
		return s[:max]
	}
	return s[:max-3] + "..."
}

func (model Model) renderTable() string {
	if len(model.searchResults) == 0 {
		return ""
	}

	pathWidth := int(float64(model.width) * 0.3)
	nameWidth := int(float64(model.width) * 0.2)
	headers := []string{"Path", "Name", "Size", "Modified", "Last Used", "Metadata Changed"}
	var rows [][]string
	for _, entry := range model.searchResults {
		rows = append(rows, []string{
			limitPathLength(entry.Path, pathWidth),
			limitPathLength(entry.Name, nameWidth),
			strconv.Itoa(int(entry.Size)),
			time.Unix(0, entry.ModificationTime).Format("2006-01-02 15:04:05"),
			time.Unix(0, entry.AccessTime).Format("2006-01-02 15:04:05"),
			time.Unix(0, entry.MetaDataChangeTime).Format("2006-01-02 15:04:05"),
		})
	}

	table := table.New().
		Border(lipgloss.HiddenBorder()).
		//BorderStyle(lipgloss.NewStyle().Foreground(model.style.BorderColor)).
		Headers(headers...).
		Rows(rows...)

	return table.Render()
}

func (model Model) searchDB(searchString string) tea.Cmd {
	return func() tea.Msg {
		results, err := search.Index(model.dbConnection, searchString)
		if err != nil {
			fmt.Println(err)
		}
		return searchResult{Rows: results, Err: nil}
	}
}

func (model Model) Init() tea.Cmd {
	return textinput.Blink
}

func (model Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	oldInput := model.inputField.Value()

	var viewPortCmd tea.Cmd
	model.viewPort, viewPortCmd = model.viewPort.Update(msg)
	cmds = append(cmds, viewPortCmd)

	var inputCmd tea.Cmd
	model.inputField, inputCmd = model.inputField.Update(msg)
	cmds = append(cmds, inputCmd)

	newInput := model.inputField.Value()
	if newInput != oldInput && newInput != "" {
		cmds = append(cmds, model.searchDB(newInput))
	}
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		model.width = msg.Width
		model.height = msg.Height
		model.viewPort.Width = msg.Width
		model.viewPort.Height = msg.Height - 10
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return model, tea.Quit
		}
		switch msg.Type {
		case tea.KeyEnter:
			inputValue := model.inputField.Value()
			return model, model.searchDB(inputValue)
		}
	case searchResult:
		model.searchResults = msg.Rows
		model.err = msg.Err
		tableString := model.renderTable()
		model.viewPort.SetContent(tableString)
		model.viewPort.GotoTop()
		return model, nil
	}
	return model, tea.Batch(cmds...)
}

func (model Model) View() string {
	if model.width == 0 {
		return "loading..."
	}
	return lipgloss.Place(
		model.width,
		model.height,
		lipgloss.Center,
		lipgloss.Top,
		lipgloss.JoinVertical(
			lipgloss.Center,
			model.style.InputField.Render(model.inputField.View()),
			"\n"+model.viewPort.View(),
		),
	)
}
