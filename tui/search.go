package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

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
	input := textinput.New()
	input.Placeholder = "Start typing to search cells.."
	input.PlaceholderStyle = getBarStyle()
	input.TextStyle = getBarStyle()
	input.Prompt = ""
	input.Focus()
	input.Blink()

	return &searchModel{input: input}
}

func (s *searchModel) Init() tea.Cmd {
	return textinput.Blink
}

// A searchMessage contains the contents of the search bar, and is
// sent to other models when the user stops typing briefly.
type searchMessage string

func searchCommand(val string) func() tea.Msg {
	return func() tea.Msg {
		return searchMessage(val)
	}
}

func (s *searchModel) View() string {
	doc := strings.Builder{}
	doc.WriteString(statusStyle.Render("Query"))
	doc.WriteString(s.input.View())

	return barStyle.
		Width(80).
		Render(doc.String())
}

func (s *searchModel) Update(msg tea.Msg) (*searchModel, tea.Cmd) {
	cmd := s.updateSubModels(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyRunes, tea.KeyBackspace:
			v := s.input.Value()
			if v != "" {
				cmd = searchCommand(s.input.Value())
			}
		}
	}

	return s, cmd
}

func (s *searchModel) updateSubModels(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	s.input, cmd = s.input.Update(msg)
	return cmd
}

var (
	statusStyle = getStatusStyle()
	barStyle    = getBarStyle()
)
