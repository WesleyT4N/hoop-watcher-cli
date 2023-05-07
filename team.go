package hoop_watcher

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

type NBATeam struct {
	Name         string
	Abbreviation string
}

const teamFileLoadErrorMessage = "Error occurred loading NBA teams"

func GetNBATeams() []NBATeam {
	f, err := os.Open("nba_teams.json")
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

func queryMatchesTeam(query string, team NBATeam) bool {
	lowerCaseQuery := strings.ToLower(query)
	lowerCaseTeamName := strings.ToLower(team.Name)
	lowerCaseTeamAbbreviation := strings.ToLower(team.Abbreviation)
	if strings.Contains(lowerCaseQuery, lowerCaseTeamName) || strings.Contains(lowerCaseQuery, lowerCaseTeamAbbreviation) {
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

func GetTeamFromQuery(query string, nbaTeams []NBATeam) *NBATeam {
	for _, team := range nbaTeams {
		if queryMatchesTeam(query, team) {
			return &team
		}
	}
	return nil
}
