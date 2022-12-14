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

	width, height int

	currModeIdx int
	modes       []searchMode
}

func newSearchModel() *searchModel {
	modes := []searchMode{
		{
			title:       "Keyword",
			placeholder: "Search via keyword matching..",
			color:       "#684EFF",
			mode:        search.Keyword,
		},
		{
			title:       "Phrase",
			placeholder: "Search via a phrase..",
			color:       "#9F2EEB",
			mode:        search.Phrase,
		},
		{
			title:       "Fuzzy",
			placeholder: "Search via fuzzy word matching..",
			color:       "#09ae70",
			mode:        search.Fuzzy,
		},
		{
			title:       "Wildcard",
			placeholder: "Search using '?' and '*' wildcards..",
			color:       "#F25D94",
			mode:        search.Wildcard,
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
		barStyle.Width(s.width-3).Render(title+input)+"\n",
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
			cmd = tea.Batch(searchCommand(s.modes[s.currModeIdx].mode, s.input.Value()))
		case tea.KeyRunes, tea.KeyBackspace:
			v := s.input.Value()
			if v == "" {
				return s, cmd
			}
			cmd = tea.Batch(searchCommand(s.modes[s.currModeIdx].mode, s.input.Value()))
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

func (s *searchModel) setDimensions(width, height int) {
	s.width, s.height = width, height
}

func (s *searchModel) updateSubModels(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	s.input, cmd = s.input.Update(msg)
	return cmd
}

func getBarStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#FFF", Dark: "#FFF"}).
		Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#262626"})
}

func getStatusStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Inherit(getBarStyle()).
		Foreground(lipgloss.Color("#FFF")).
		Padding(0, 1).
		MarginRight(1).
		Bold(true)
}
