package main

import (
	hoop_watcher "github.com/WesleyT4N/hoop-watcher-cli"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type model struct {
	list         list.Model
	teams        []hoop_watcher.NBATeam
	cursor       int
	selectedTeam hoop_watcher.NBATeam
}

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func initialModel() model {
	allTeams := hoop_watcher.GetNBATeams(teamFilePath)
	items := []list.Item{}
	for _, team := range allTeams {
		items = append(items, item{
			title: team.Abbreviation,
			desc:  team.Name,
		})
	}
	list := list.New(items, list.NewDefaultDelegate(), 0, 0)
	list.Title = "Hoop Watcher CLI"
	return model{
		list: list,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}
