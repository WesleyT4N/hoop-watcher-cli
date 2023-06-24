package main

import (
	"fmt"
	hoop_watcher "github.com/WesleyT4N/hoop-watcher-cli"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	teams        []hoop_watcher.NBATeam
	cursor       int
	selectedTeam hoop_watcher.NBATeam
}

func initialModel() model {
	allTeams := hoop_watcher.GetNBATeams(teamFilePath)
	return model{
		teams:        allTeams,
		cursor:       0,
		selectedTeam: allTeams[0],
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
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.teams)-1 {
				m.cursor++
			}
		case "enter", " ":
			m.selectedTeam = m.teams[m.cursor]
		}
	}
	return m, nil
}

func (m model) View() string {
	s := "Hoop Watcher CLI"
	for i, team := range m.teams {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		checked := " "
		if team == m.selectedTeam {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, team.Name)
	}

	s += "\nPress q to quit.\n"

	return s
}
