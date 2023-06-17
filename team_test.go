package hoop_watcher_test

import (
	"testing"

	hoop_watcher "github.com/WesleyT4N/hoop-watcher-cli"
)

const NUMBER_OF_NBA_TEAMS = 30

func TestGetTeams(t *testing.T) {
	t.Run("it loads teams JSON file", func(t *testing.T) {
		teams := hoop_watcher.GetNBATeams("./nba_teams.json")
		if len(teams) != NUMBER_OF_NBA_TEAMS {
			t.Errorf("expected %d number of teams but loaded %d", len(teams), NUMBER_OF_NBA_TEAMS)
		}

		for _, team := range teams {
			if team.Name == "" {
				t.Errorf("invalid team Name %s", team.Name)
			}
			if team.Name == "" {
				t.Errorf("invalid team Abbreviation %s", team.Abbreviation)
			}
		}
	})
}

type getTeamFromQueryTestCase struct {
	Query            string
	ExpectedTeamName string
}

func TestGetTeamFromQuery(t *testing.T) {
	t.Run("gets team from name", func(t *testing.T) {
		cases := []getTeamFromQueryTestCase{
			{
				"NEW YORK KNICKS",
				"New York Knicks",
			},
			{
				"KNICKS",
				"New York Knicks",
			},
			{
				"knicks",
				"New York Knicks",
			},
			{
				"nyk",
				"New York Knicks",
			},
			{
				"ny knicks",
				"New York Knicks",
			},
		}
		for _, c := range cases {
			got := hoop_watcher.GetTeamFromQuery(c.Query, hoop_watcher.GetNBATeams("./nba_teams.json"))
			want := c.ExpectedTeamName
			if got == nil {
				t.Errorf("got %v want %s", got, want)
			}
			if got.Name != want {
				t.Errorf("got %v want %s", got.Name, want)
			}
		}
	})
}
