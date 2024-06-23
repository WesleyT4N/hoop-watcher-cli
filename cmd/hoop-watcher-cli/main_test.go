package main

import (
	"reflect"
	"testing"
	"time"

	hoop_watcher "github.com/WesleyT4N/hoop-watcher-cli"
)

func TestParseDate(t *testing.T) {
	t.Run("it returns time.Time object given a string of the form YYYY-MM-DD", func(t *testing.T) {
		got, err := parseDate("2020-01-01")
		if err != nil {
			t.Fatalf("Found err: %v", err)
		}
		want := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		if got != want {
			t.Fatalf("got %v, want %v", got, want)
		}
	})

	t.Run("it returns time.Time object given a string of the form YYYY-MM", func(t *testing.T) {
		got, err := parseDate("2020-01")
		if err != nil {
			t.Fatalf("Found err: %v", err)
		}
		want := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		if got != want {
			t.Fatalf("got %v, want %v", got, want)
		}
	})

	t.Run("raises error if invalid date", func(t *testing.T) {
		_, err := parseDate("2020-13")
		if err == nil {
			t.Fatal("Expected error but err was nil")
		}
		got := err.Error()
		want := "Invalid date"
		if got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
	})
}

func TestParseTeams(t *testing.T) {
	availableTeams := hoop_watcher.GetNBATeamsFromJSON(teamFilePath)

	t.Run("parses two team string into NBATeam", func(t *testing.T) {
		got, err := parseTeams("knicks,grizzlies", availableTeams)
		if err != nil {
			t.Fatalf("Found err: %v", err)
		}

		want := []hoop_watcher.NBATeam{
			{
				Name:         "New York Knicks",
				Abbreviation: "NYK",
			}, {
				Name:         "Memphis Grizzlies",
				Abbreviation: "MEM",
			},
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("parses singular team string into NBATeam", func(t *testing.T) {
		got, err := parseTeams("knicks", availableTeams)
		if err != nil {
			t.Fatalf("Found err: %v", err)
		}

		want := []hoop_watcher.NBATeam{
			{
				Name:         "New York Knicks",
				Abbreviation: "NYK",
			},
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("returns error on more than 2 teams", func(t *testing.T) {
		_, err := parseTeams("knicks,knicks,knicks", availableTeams)
		if err == nil {
			t.Fatal("Expected error but found none")
		}
		got := err.Error()
		want := "Invalid number of teams given"
		if got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
	})
	t.Run("returns nil on 0 teams", func(t *testing.T) {
		got, err := parseTeams("", availableTeams)
		if err != nil {
			t.Fatalf("Found err: %v", err)
		}
		if got != nil {
			t.Fatalf("Expected nil but found %v", got)
		}
	})
}
