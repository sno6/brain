package tui

import (
	"github.com/sno6/brain/search"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	statusStyle = getStatusStyle()
	barStyle    = getBarStyle()
)

type searchMode struct {
	mode        search.Mode
	title       string
	placeholder string
	color       lipgloss.Color
}

type searchModel struct {
	input textinput.Model
	help  *helpModel

	currModeIdx int
	modes       []searchMode
}

func newSearchModel() *searchModel {
	modes := []searchMode{
		{
			title:       "Match",
			placeholder: "Search via keyword matching..",
			color:       "62",
			mode:        search.Match,
		},
		{
			title:       "Fuzzy",
			placeholder: "Search via fuzzy word matching..",
			color:       "22",
			mode:        search.Fuzzy,
		},
		{
			title:       "Regexp",
			placeholder: "Search via regular expressions..",
			color:       "55",
			mode:        search.Regexp,
		},
	}

	input := textinput.New()
	input.Placeholder = modes[0].placeholder
	input.PlaceholderStyle = getBarStyle()
	input.TextStyle = getBarStyle()
	input.Prompt = ""
	input.Focus()
	input.Blink()

	return &searchModel{
		input: input,
		help:  newHelpModel(PageSearch),
		modes: modes,
	}
}

func (s *searchModel) Init() tea.Cmd {
	return textinput.Blink
}

func (s *searchModel) View() string {
	mode := s.modes[s.currModeIdx]
	title := statusStyle.Background(mode.color).Render(mode.title)
	input := s.input.View()

	return lipgloss.JoinVertical(
		0,
		barStyle.Width(80).Render(title+input)+"\n",
		s.help.View(),
	)
}

func (s *searchModel) Update(msg tea.Msg) (*searchModel, tea.Cmd) {
	cmd := s.updateSubModels(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab:
			s.toggleMode()
		case tea.KeyRunes, tea.KeyBackspace:
			v := s.input.Value()
			if v == "" {
				return s, cmd
			}

			cmd = tea.Batch(
				searchCommand(s.modes[s.currModeIdx].mode, s.input.Value()))
		}
	}

	return s, cmd
}

func (s *searchModel) toggleMode() {
	idx := s.currModeIdx
	if idx == len(s.modes)-1 {
		s.currModeIdx = 0
	} else {
		s.currModeIdx = idx + 1
	}

	// A bit gross but we use the input's own placeholder model here
	// so we need to update its contents. In render is where everything
	// else will change.
	newMode := s.modes[s.currModeIdx]
	s.input.Placeholder = newMode.placeholder
}

func (s *searchModel) updateSubModels(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	s.input, cmd = s.input.Update(msg)
	return cmd
}

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
