package hoop_watcher

import (
	"fmt"
	"io"
	"time"
)

const DAILY_DATE_FORMAT = "2006-01-02"

func TeamHighlightQueryString(team string, reqTime time.Time) string {
	dateStr := reqTime.Format(DAILY_DATE_FORMAT)
	return fmt.Sprintf("%s NBA Highlights %s", team, dateStr)
}

func GetHighlights(team string, out io.Writer) {
	fmt.Fprintf(out, "Getting highlights for team %s", team)
}
