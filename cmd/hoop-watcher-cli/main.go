package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"

	hoop_watcher "github.com/WesleyT4N/hoop-watcher-cli"
	"github.com/joho/godotenv"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

func main() {
	youtubeClient := newYoutubeClient()
	stdout := os.Stdout
	team := scanTeam()
	highlights := hoop_watcher.GetHighlights(team, stdout, youtubeClient)
	openHighlight(highlights)
}

func openHighlight(highlights []url.URL) {
	fmt.Print("Which one do you want to view? (enter the corresponding number)\n> ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if scanner.Err() != nil {
		log.Fatal("Error occurred parsing num")
	}
	num, err := strconv.Atoi(scanner.Text())
	if err != nil || !(0 <= num-1 && num-1 <= len(highlights)) {
		log.Fatal("Error occurred parsing num")
	}

	fmt.Println("Opening the highlight in your browser...")
	cmd := exec.Command("open", highlights[num-1].String())
	err = cmd.Run()
	if err != nil {
		fmt.Println("could not open highlight: ", err)
	}
}

func scanTeam() string {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter the NBA team you want to get highlights for:\n> ")
	scanner.Scan()
	if scanner.Err() != nil {
		log.Fatal("Error occurred parsing team")
	}

	team := scanner.Text()
	fmt.Println(team)
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
