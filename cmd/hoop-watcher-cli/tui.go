package main

import (
	"log"
	"os"
	"os/exec"
	"time"

	hoop_watcher "github.com/WesleyT4N/hoop-watcher-cli"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"google.golang.org/api/youtube/v3"
)

var (
	docStyle  = lipgloss.NewStyle().Margin(1, 2)
	tableSyle = table.DefaultStyles()
)

type model struct {
	list            list.Model
	table           table.Model
	hasSelectedTeam bool
	highlights      map[Team][]hoop_watcher.Highlight
	yt              *youtube.Service
}

type Team struct {
	team hoop_watcher.NBATeam
}

func (i Team) FilterValue() string { return i.team.Name + i.team.Abbreviation }
func (i Team) Title() string       { return i.team.Abbreviation }
func (i Team) Description() string { return i.team.Name }

func initList() list.Model {
	allTeams := hoop_watcher.GetNBATeamsFromJSON(teamFilePath)
	var items []list.Item
	for _, team := range allTeams {
		items = append(items, Team{
			team: team,
		})
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Hoop Watcher CLI"
	l.SetShowStatusBar(true)
	l.DisableQuitKeybindings()
	return l
}

func initTable() table.Model {
	columns := []table.Column{
		{Title: "Video", Width: 100},
		{Title: "URL", Width: 100},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(false),
		table.WithHeight(10),
	)
	return t
}

func initialModel() model {
	return model{
		list:            initList(),
		table:           initTable(),
		hasSelectedTeam: false,
		highlights:      map[Team][]hoop_watcher.Highlight{},
		yt:              newYoutubeClient(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

type highlightLookupMsg struct {
	highlights []hoop_watcher.Highlight
}

func lookupHighlight(team hoop_watcher.NBATeam, yt *youtube.Service) tea.Cmd {
	return func() tea.Msg {
		return highlightLookupMsg{
			highlights: hoop_watcher.GetHighlightsForTUI(team, time.Now(), yt),
		}
	}
}

func (m model) Update(msg tea.Msg) (n tea.Model, cmd tea.Cmd) {
	log.Printf("Msg: %T, %v\n", msg, msg)
	log.Printf("Selected Team: %v\n", m.list.SelectedItem())
	log.Printf("Has Selected: %v\n", m.hasSelectedTeam)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			selectedItem := m.list.SelectedItem()
			if selectedItem != nil && !m.hasSelectedTeam && !m.list.SettingFilter() {
				selectedTeam := selectedItem.(Team)
				m.hasSelectedTeam = true
				return m, lookupHighlight(selectedTeam.team, m.yt)
			} else if m.table.Focused() {
				cmd := exec.Command("open", m.table.SelectedRow()[1])
				if cmd.Run() != nil {
					os.Exit(1)
				}
			}
		case "esc":
			m.list.ResetFilter()
			if m.list.SelectedItem() != nil && m.hasSelectedTeam {
				m.list.ResetSelected()
				m.hasSelectedTeam = false
			}
			if m.table.Focused() {
				m.table.Blur()
			}
		}
	case highlightLookupMsg:
		highlights := msg.highlights
		selectedTeam := m.list.SelectedItem().(Team)
		m.highlights[selectedTeam] = highlights
		var rows []table.Row
		for _, h := range highlights {
			rows = append(rows, table.Row{h.Title, h.URL.String()})
		}
		m.table.SetRows(rows)
		m.table.Focus()
		return m, nil
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	if m.table.Focused() {
		m.table, cmd = m.table.Update(msg)
	} else {
		m.list, cmd = m.list.Update(msg)
	}
	return m, cmd
}

func (m model) View() string {
	if m.hasSelectedTeam {
		selectedTeam := m.list.SelectedItem().(Team)
		if m.highlights[selectedTeam] != nil {
			{
				return docStyle.Render(m.table.View())
			}
		}
	}
	return docStyle.Render(m.list.View())
}
