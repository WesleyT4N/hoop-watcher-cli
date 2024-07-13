package hoop_watcher

import (
	"database/sql"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const initDB = `
CREATE TABLE IF NOT EXISTS teams(
    id INTEGER PRIMARY KEY NOT NULL,
    name TEXT NOT NULL,
	full_name TEXT NOT NULL,
    abbreviation TEXT NOT NULL,
	city TEXT NOT NULL,
	conference TEXT NOT NULL,
	division TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS games(
    id SERIAL PRIMARY KEY,
    home_team_id INTEGER NOT NULL REFERENCES teams(id),
    away_team_id INTEGER NOT NULL REFERENCES teams(id),
    date DATE NOT NULL,
	UNIQUE(home_team_id, away_team_id, date)
);

CREATE TABLE IF NOT EXISTS game_highlights(
    id SERIAL PRIMARY KEY,
    game_id INTEGER NOT NULL REFERENCES games(id),
    url VARCHAR(255) NOT NULL
);
`

type HoopWatcherDB interface {
	GetAllTeams() ([]NBATeam, error)
	GetTeamByAbbrev(abbrev string) (NBATeam, error)
	GetTeamHighlights(teamId int) ([]Highlight, error)
}

type SqliteHoopWatcherDB struct {
	db *sql.DB
}

func NewSqliteHoopWatcherDB(filePath string) (*SqliteHoopWatcherDB, error) {
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(initDB); err != nil {
		return nil, err
	}

	return &SqliteHoopWatcherDB{db: db}, nil
}

func (h *SqliteHoopWatcherDB) Close() error {
	return h.db.Close()
}

func (h *SqliteHoopWatcherDB) addAllTeams(filePath string) error {
	teams := GetNBATeamsFromJSON(filePath)
	stmt, err := h.db.Prepare("INSERT OR IGNORE INTO teams(id, name, full_name, abbreviation, city, conference, division) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	for _, team := range teams {
		_, err = stmt.Exec(team.Id, team.Name, team.FullName, team.Abbreviation, team.City, team.Conference, team.Division)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *SqliteHoopWatcherDB) InitData(teamFilePath string) error {
	if err := h.addAllTeams(teamFilePath); err != nil {
		return err
	}
	return nil
}

func (h *SqliteHoopWatcherDB) GetAllTeams() ([]NBATeam, error) {
	rows, err := h.db.Query("SELECT * FROM teams")
	if err != nil {
		return []NBATeam{}, err
	}
	defer rows.Close()

	teams := []NBATeam{}
	for rows.Next() {
		team := NBATeam{}
		if err := rows.Scan(
			&team.Id,
			&team.Name,
			&team.FullName,
			&team.Abbreviation,
			&team.City,
			&team.Conference,
			&team.Division,
		); err != nil {
			return []NBATeam{}, err
		}
		teams = append(teams, team)
	}
	return teams, nil
}

func (h *SqliteHoopWatcherDB) GetTeamByAbbrev(abbrev string) (NBATeam, error) {
	var team NBATeam
	row := h.db.QueryRow("SELECT * FROM teams WHERE abbrev = ?", strings.ToUpper(abbrev))
	if err := row.Scan(&team.Id, &team.Name, &team.Abbreviation); err != nil {
		return NBATeam{}, err
	}
	return team, nil
}

func (h *SqliteHoopWatcherDB) GetTeamHighlights(id int) ([]Highlight, error) {
	return []Highlight{}, nil
}
