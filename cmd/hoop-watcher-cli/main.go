package main

import (
	"log"
	"net/http"
	"os"

	hoop_watcher "github.com/WesleyT4N/hoop-watcher-cli"
	"github.com/joho/godotenv"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

func main() {
	youtubeClient := newYoutubeClient()
	stdout := os.Stdout
	hoop_watcher.GetHighlights("Knicks", stdout, youtubeClient)
}

func newYoutubeClient() *youtube.Service {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error occurred loading .env file")
	}

	youtubeApiKey := os.Getenv("YOUTUBE_API_KEY")
	youtubeClient, err := youtube.New(
		&http.Client{
			Transport: &transport.APIKey{Key: youtubeApiKey},
		},
	)
	if err != nil {
		log.Fatal("Error occurred setting up Youtube Client")
	}

	return youtubeClient
}
