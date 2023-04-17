package hoop_watcher

import (
	"fmt"
	"io"
)

func GetHighlights(team string, out io.Writer) {
	fmt.Fprintf(out, "Getting highlights for team %s", team)
}
