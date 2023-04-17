package hoop_watcher

import (
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/api/youtube/v3"
)

const DAILY_DATE_FORMAT = "2006-01-02"

func TeamHighlightQueryString(team string, reqTime time.Time) string {
	dateStr := reqTime.Format(DAILY_DATE_FORMAT)
	return fmt.Sprintf("'%s NBA Full Game Highlights %s'", team, dateStr)
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

func GetHighlights(team string, out io.Writer, youtubeClient *youtube.Service) {

	fmt.Fprintf(out, "Getting highlights for team: %s\n\n", team)
	youtubeQueryString := TeamHighlightQueryString(team, time.Now())

	videos, err := searchListByQ(youtubeClient, youtubeQueryString)
	if err != nil {
		log.Fatalf("Error occurred fething youtube video urls")
	}

	fmt.Fprintln(out, "Found these matching highlights:")
	for _, video := range videos {
		videoURL := fmt.Sprintf("https://www.youtube.com/watch?v=%v", video.Id.VideoId)
		fmt.Fprintf(out, videoURL+" : "+video.Snippet.Title+"\n")
	}
}
