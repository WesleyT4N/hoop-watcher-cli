package main

import (
	"log"
	"net/http"

	hoop_watcher "github.com/WesleyT4N/hoop-watcher-cli"
)

func main() {
	log.Printf("Starting server on port 8080")
	http.HandleFunc("/", hoop_watcher.GetRoot)

	log.Fatal(http.ListenAndServe(":8080", nil))

}
