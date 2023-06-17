package hoop_watcher_test

import (
	"testing"
	"time"

	hoop_watcher "github.com/WesleyT4N/hoop-watcher-cli"
)

func TestHighlightQueryString(t *testing.T) {
	t.Run("it generates Youtube Query string for team and date", func(t *testing.T) {
		loc, _ := time.LoadLocation("Local")
		date := time.Date(2023, time.January, 1, 0, 0, 0, 0, loc)
		got := hoop_watcher.TeamHighlightQueryString([]string{"Knicks"}, date)

		want := "'Knicks NBA Full Game Highlights 2023-01-01'"
		if got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})
}
