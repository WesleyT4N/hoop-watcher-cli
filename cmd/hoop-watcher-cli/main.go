package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	hoop_watcher "github.com/WesleyT4N/hoop-watcher-cli"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var SupportedDateFormats = []string{
	"2006-01-02",
	"2006-01",
}

var teamFilePath = path.Join(os.Getenv("HOME"), "bin", hoop_watcher.TeamFileName)

func parseDate(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Now(), nil
	}
	for _, format := range SupportedDateFormats {
		gameDate, err := time.Parse(format, dateStr)
		if err != nil {
			continue
		}
		return gameDate, nil
	}
	return time.Now(), fmt.Errorf("Invalid date")
}

func parseTeams(teamStr string, availableTeams []hoop_watcher.NBATeam) (teams []hoop_watcher.NBATeam, err error) {
	trimmedTeamStr := strings.TrimSpace(teamStr)
	if trimmedTeamStr == "" {
		return nil, nil
	}
	rawTeams := strings.Split(trimmedTeamStr, ",")
	if len(rawTeams) > 2 {
		return nil, fmt.Errorf("Invalid number of teams given")
	}

	for _, team := range rawTeams {
		parsedTeam := hoop_watcher.GetTeamFromQuery(team, availableTeams)
		if parsedTeam != nil {
			teams = append(teams, *parsedTeam)
		}
	}
	return teams, nil
}

func parseFavoriteTeams(favoriteStr string) (favoriteTeams []int, err error) {
	trimmedFavoriteStr := strings.TrimSpace(favoriteStr)
	if trimmedFavoriteStr == "" {
		return nil, nil
	}
	rawFavoriteTeams := strings.Split(trimmedFavoriteStr, ",")
	for _, favoriteTeam := range rawFavoriteTeams {
		parsedFavoriteTeam, err := strconv.Atoi(favoriteTeam)
		if err != nil {
			return nil, fmt.Errorf("Invalid team id")
		}
		favoriteTeams = append(favoriteTeams, parsedFavoriteTeam)
	}
	return favoriteTeams, nil
}

func parseFlags(availableTeams []hoop_watcher.NBATeam) (useTui bool, date time.Time, teams []hoop_watcher.NBATeam, favTeams []int, err error) {
	tuiArg := flag.Bool("tui", false, "Use the TUI")
	dateArg := flag.String("d", "", "Date of the highlights to fetch in the format YYYY-MM-DD")
	teamsArg := flag.String("tm", "", "Which teams are playing (max 2) joined by ','")
	favoriteArg := flag.String("f", "", "Favorite the team inputted")
	flag.Parse()

	date, dateErr := parseDate(*dateArg)
	if dateErr != nil {
		return *tuiArg, date, teams, favTeams, dateErr
	}
	teams, teamErr := parseTeams(*teamsArg, availableTeams)
	if teamErr != nil {
		return *tuiArg, date, teams, favTeams, teamErr
	}
	favTeams, favTeamErr := parseFavoriteTeams(*favoriteArg)
	if favTeamErr != nil {
		return *tuiArg, date, teams, favTeams, favTeamErr
	}

	return *tuiArg, date, teams, favTeams, nil
}

func runCLI() {
	db, err := hoop_watcher.NewSqliteHoopWatcherDB("hoop-watcher-cli.db")
	allTeams := hoop_watcher.GetNBATeamsFromDB(db)
	allTeamsById := make(map[int]hoop_watcher.NBATeam)
	for _, team := range allTeams {
		allTeamsById[team.Id] = team
	}

	if err != nil {
		log.Fatal("Error occurred setting up DB")
	}
	if db.InitData(teamFilePath) != nil {
		log.Fatal("Error occurred initializing data in DB")
	}
	useTui, date, teams, favTeams, err := parseFlags(allTeams)
	if useTui {
		runTUI()
		return
	}

	youtubeClient := newYoutubeClient()
	stdout := os.Stdout

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if len(favTeams) > 0 {
		fmt.Println("Favorited the following teams:")
		for _, teamId := range favTeams {
			if err := db.SetTeamFavorite(teamId, true); err != nil {
				fmt.Println("Error occurred setting team as favorite")
				os.Exit(1)
			}
			fmt.Println(allTeamsById[teamId].Name)
		}
		return
	}

	if len(teams) == 0 {
		teams, err = scanTeam(allTeams)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}

	highlights := hoop_watcher.GetHighlights(teams, stdout, youtubeClient, date)
	err = openHighlight(highlights)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func runTUI() {
	err := godotenv.Load(path.Join(os.Getenv("HOME"), ".env"))
	if err != nil {
		log.Fatal("Error occurred loading .env file")
	}
	if os.Getenv("DEBUG") == "1" {
		f, err := tea.LogToFile("debug.log", "[DEBUG]")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func main() {
	runCLI()
}

func openHighlight(highlights []url.URL) error {
	fmt.Print("Which one do you want to view? (enter the corresponding number)\n> ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if scanner.Err() != nil {
		return errors.New("Error occurred parsing num")
	}
	num, err := strconv.Atoi(scanner.Text())
	if err != nil || !(0 <= num-1 && num-1 <= len(highlights)) {
		return errors.New("Error occurred parsing num")
	}

	fmt.Println("Opening the highlight in your browser...")
	cmd := exec.Command("open", highlights[num-1].String())
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("could not open highlight: %v", err)
	}
	return nil
}

func scanTeam(allTeams []hoop_watcher.NBATeam) ([]hoop_watcher.NBATeam, error) {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter the NBA team you want to get highlights for:\n> ")
	scanner.Scan()
	if scanner.Err() != nil {
		return nil, errors.New("Error occurred parsing team")
	}

	team := scanner.Text()

	parsedTeam := hoop_watcher.FuzzyGetTeamFromQuery(team, allTeams)
	if parsedTeam == nil {
		return nil, errors.New("Unknown team")
	}
	return []hoop_watcher.NBATeam{*parsedTeam}, nil
}

func newYoutubeClient() *youtube.Service {
	err := godotenv.Load(path.Join(os.Getenv("HOME"), ".env"))
	if err != nil {
		log.Fatal("Error occurred loading .env file")
	}

	youtubeApiKey := os.Getenv("YOUTUBE_API_KEY")
	ctx := context.Background()
	youtubeClient, err := youtube.NewService(ctx, option.WithAPIKey(youtubeApiKey))
	if err != nil {
		log.Fatal("Error occurred setting up Youtube Client")
	}

	return youtubeClient
}
