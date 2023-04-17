package main

import (
	"bufio"
	"fmt"
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
	team := scanTeam()
	hoop_watcher.GetHighlights(team, stdout, youtubeClient)
}

func scanTeam() string {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter the NBA team you want to get highlights for:\n> ")

	scanner.Scan()
	if scanner.Err() != nil {
		log.Fatal("Error occurred parsing team")
	}

	team := scanner.Text()
	return team
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
