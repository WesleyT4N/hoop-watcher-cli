package hoop_watcher

import (
	"fmt"
	"io"
	"log"
	"net/url"
	"strings"
	"time"

	"google.golang.org/api/youtube/v3"
)

const DAILY_DATE_FORMAT = "2006-01-02"
const HUMAN_DATE_FORMAT = "January 2, 2006"

func TeamHighlightQueryString(teamNames []string, reqTime time.Time) string {
	dateStr := reqTime.Format(HUMAN_DATE_FORMAT)
	return fmt.Sprintf("'%s NBA Full Game Highlights %s'", strings.Join(teamNames, " vs "), dateStr)
}

// searchListByQ searches for videos based on a keyword query
func searchListByQ(service *youtube.Service, keywordQuery string) ([]*youtube.SearchResult, error) {
	call := service.Search.List([]string{"id", "snippet"}).
		Q(keywordQuery).
		Type("video").MaxResults(3)

	response, err := call.Do()
	if err != nil {
		return nil, err
	}

	return response.Items, nil
}

type Highlight struct {
	Title string
	URL   url.URL
}

func GetHighlightsForTUI(team NBATeam, time time.Time, youtubeClient *youtube.Service) (highlights []Highlight) {
	teamNames := []string{}
	teamNames = append(teamNames, team.Name)
	youtubeQueryString := TeamHighlightQueryString(teamNames, time)
	videos, err := searchListByQ(youtubeClient, youtubeQueryString)
	if err != nil {
		log.Fatalf("Error occurred fething youtube video urls")
	}

	for _, video := range videos {
		videoURL := fmt.Sprintf("https://www.youtube.com/watch?v=%v", video.Id.VideoId)
		parsedUrl, err := url.Parse(videoURL)
		if err != nil {
			log.Fatalf("Error parsing video URL")
		}
		highlights = append(highlights, Highlight{Title: video.Snippet.Title, URL: *parsedUrl})
	}
	return highlights
}

func GetHighlights(teams []NBATeam, out io.Writer, youtubeClient *youtube.Service, time time.Time) []url.URL {
	teamNames := []string{}
	for _, t := range teams {
		teamNames = append(teamNames, t.Name)
	}
	fmt.Fprintf(out, "Getting highlights for the %v\n\n", strings.Join(teamNames, " vs "))
	youtubeQueryString := TeamHighlightQueryString(teamNames, time)
	fmt.Println(youtubeQueryString)

	videos, err := searchListByQ(youtubeClient, youtubeQueryString)
	if err != nil {
		log.Fatalf("Error occurred fething youtube video urls")
	}

	fmt.Fprintln(out, "Found these matching highlights:")
	var highlightUrls []url.URL
	for i, video := range videos {
		videoURL := fmt.Sprintf("https://www.youtube.com/watch?v=%v", video.Id.VideoId)
		fmt.Fprintf(out, "[%d] %s : %s\n", i+1, videoURL, video.Snippet.Title)
		parsedUrl, err := url.Parse(videoURL)
		if err != nil {
			log.Fatalf("Error parsing video URL")
		}
		highlightUrls = append(highlightUrls, *parsedUrl)
	}
	return highlightUrls
}
