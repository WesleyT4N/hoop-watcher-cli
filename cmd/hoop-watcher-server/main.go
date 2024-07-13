package main

import (
	"log"
	"net/http"
	"os"
	"path"

	hoop_watcher "github.com/WesleyT4N/hoop-watcher-cli"
	"github.com/joho/godotenv"
)

var teamFilePath = path.Join(os.Getenv("HOME"), "bin", hoop_watcher.TeamFileName)

func main() {
	err := godotenv.Load(path.Join(os.Getenv("HOME"), ".env"))
	if err != nil {
		log.Fatal("Error occurred loading .env file")
	}

	router := http.NewServeMux()
	db, err := hoop_watcher.NewSqliteHoopWatcherDB("hoop-watcher-cli.db")
	if err != nil {
		log.Fatalf("Error occurred creating database connection: %v", err)
	}
	err = db.InitData(teamFilePath)
	if err != nil {
		log.Fatal("Error occurred initializing data in DB", err)
	}
	h := hoop_watcher.NewBaseHandler(db)

	router.HandleFunc("/", h.GetRoot)

	router.HandleFunc("GET /teams", h.GetTeams)
	router.HandleFunc("GET /teams/{abbrev}", h.GetTeam)
	router.HandleFunc("GET /teams/{abbrev}/highlights", h.GetTeamHighlights)

	log.Printf("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))

}
