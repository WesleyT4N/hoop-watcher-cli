package main

import (
	"os"

	hoop_watcher "github.com/WesleyT4N/hoop-watcher-cli"
)

func main() {
	stdout := os.Stdout
	hoop_watcher.GetHighlights("Knicks", stdout)
}
