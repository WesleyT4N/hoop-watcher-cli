package main

import (
	"log"
	"net/http"
	"os"
	"path"

	hoop_watcher "github.com/WesleyT4N/hoop-watcher-cli"
)

var teamFilePath = path.Join(os.Getenv("HOME"), "bin", hoop_watcher.TeamFileName)

func main() {
	router := http.NewServeMux()
	db, err := hoop_watcher.NewHoopWatcherDB("hoop-watcher-cli.db")
	if db.InitData(teamFilePath) != nil {
		log.Fatal("Error occurred initializing data in DB")
	}
	if err != nil {
		log.Fatal("Error occurred creating database connection")
	}
	h := hoop_watcher.NewBaseHandler(db)

	router.HandleFunc("/", h.GetRoot)
	router.HandleFunc("GET /team", h.GetTeams)
	router.HandleFunc("GET /team/{abbrev}", h.GetTeam)

	log.Printf("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))

}
