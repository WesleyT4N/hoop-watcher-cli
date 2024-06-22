package hoop_watcher

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

type NBATeam struct {
	Id           int
	Name         string
	Abbreviation string
	IsFavorited  bool
}

const TeamFileName = "nba_teams.json"

const teamFileLoadErrorMessage = "Error occurred loading NBA teams"

func GetNBATeamsFromJSON(filePath string) []NBATeam {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(teamFileLoadErrorMessage)
	}

	var nbaTeams []NBATeam
	err = json.NewDecoder(f).Decode(&nbaTeams)
	f.Close()
	if err != nil {
		log.Fatal(teamFileLoadErrorMessage)
	}

	return nbaTeams
}

func getTeamMatchTokens(team NBATeam) (teamToMatchTokens []string) {
	lowerCaseTeamName := strings.ToLower(team.Name)
	lowerCaseTeamAbbreviation := strings.ToLower(team.Abbreviation)
	splitTeamName := strings.Split(lowerCaseTeamName, " ")
	teamCity := strings.Join(splitTeamName[:len(splitTeamName)-1], " ")
	shortenedTeamName := splitTeamName[len(splitTeamName)-1]
	return []string{
		lowerCaseTeamName,
		lowerCaseTeamAbbreviation,
		teamCity,
		shortenedTeamName,
	}
}

func queryMatchesTeam(query string, team NBATeam) bool {
	lowerCaseQuery := strings.ToLower(query)
	lowerCaseTeamName := strings.ToLower(team.Name)
	lowerCaseTeamAbbreviation := strings.ToLower(team.Abbreviation)
	if strings.Contains(lowerCaseQuery, lowerCaseTeamName) || lowerCaseQuery == lowerCaseTeamAbbreviation {
		return true
	}

	splitTeamName := strings.Split(lowerCaseTeamName, " ")
	teamCity := strings.Join(splitTeamName[:len(splitTeamName)-1], " ")
	shortenedTeamName := splitTeamName[len(splitTeamName)-1]
	if strings.Contains(lowerCaseQuery, teamCity) || strings.Contains(lowerCaseQuery, shortenedTeamName) {
		return true
	}
	return false
}

func FuzzyGetTeamFromQuery(query string, nbaTeams []NBATeam) *NBATeam {
	matchTokensByTeam := map[NBATeam][]string{}
	for _, team := range nbaTeams {
		matchTokensByTeam[team] = getTeamMatchTokens(team)
	}
	var closestTeam NBATeam
	maxMatches := 0
	for team, matchTokens := range matchTokensByTeam {
		matches := len(fuzzy.FindNormalizedFold(query, matchTokens))
		if matches > maxMatches {
			closestTeam = team
			maxMatches = matches
		}
	}
	return &closestTeam
}

func GetTeamFromQuery(query string, nbaTeams []NBATeam) *NBATeam {
	for _, team := range nbaTeams {
		if queryMatchesTeam(query, team) {
			return &team
		}
	}
	return nil
}
