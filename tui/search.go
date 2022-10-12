package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type searchMessage string

func getBarStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
		Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})
}

func getStatusStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Inherit(getBarStyle()).
		Background(lipgloss.Color("62")).
		Foreground(lipgloss.Color("230")).
		Padding(0, 1).
		MarginRight(1)
}

type searchModel struct {
	input textinput.Model
}

func newSearchModel() *searchModel {
	return &searchModel{input: initInput()}
}

func initInput() textinput.Model {
	query := textinput.New()
	query.Placeholder = "Start typing to search cells.."
	query.PlaceholderStyle = getBarStyle()
	query.TextStyle = getBarStyle()
	query.Prompt = ""
	query.Focus()
	query.Blink()
	return query
}

func (s *searchModel) Init() tea.Cmd {
	return textinput.Blink
}

func searchCommand(val string) func() tea.Msg {
	return func() tea.Msg {
		return searchMessage(val)
	}
}

func (s *searchModel) Update(msg tea.Msg) (*searchModel, tea.Cmd) {
	var cmds []tea.Cmd

	// On any keyboard action..
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyRunes:
			cmds = append(cmds, searchCommand(s.input.Value()))
		}
	}

	var cmd tea.Cmd
	s.input, cmd = s.input.Update(msg)
	cmds = append(cmds, cmd)

	return s, tea.Batch(cmds...)
}

var (
	statusStyle = getStatusStyle()
	barStyle    = getBarStyle()
)

func (s *searchModel) View() string {
	doc := strings.Builder{}
	doc.WriteString(statusStyle.Render("Query"))
	doc.WriteString(s.input.View())
	return barStyle.Width(100).Render(doc.String())
}
