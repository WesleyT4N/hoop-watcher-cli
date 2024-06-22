package hoop_watcher

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const initDB = `
CREATE TABLE IF NOT EXISTS teams(
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    abbrev TEXT NOT NULL,
    is_favorited BOOLEAN NOT NULL DEFAULT FALSE,
	UNIQUE(name, abbrev)
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

type HoopWatcherDB struct {
	db *sql.DB
}

func NewHoopWatcherDB(filePath string) (*HoopWatcherDB, error) {
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(initDB); err != nil {
		return nil, err
	}

	return &HoopWatcherDB{db: db}, nil
}

func (h *HoopWatcherDB) Close() error {
	return h.db.Close()
}

func (h *HoopWatcherDB) addTeam(name, abbrev string) error {
	// prepared statement
	stmt, err := h.db.Prepare("INSERT OR IGNORE INTO teams(name, abbrev) VALUES(?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(name, abbrev)
	return err
}

func (h *HoopWatcherDB) addAllTeams(filePath string) error {
	teams := GetNBATeamsFromJSON(filePath)
	for _, team := range teams {
		if err := h.addTeam(team.Name, team.Abbreviation); err != nil {
			return err
		}
	}
	return nil
}

func (h *HoopWatcherDB) InitData(teamFilePath string) error {
	if err := h.addAllTeams(teamFilePath); err != nil {
		return err
	}
	return nil
}

func (h *HoopWatcherDB) GetAllTeams() ([]NBATeam, error) {
	rows, err := h.db.Query("SELECT * FROM teams")
	if err != nil {
		return []NBATeam{}, err
	}
	defer rows.Close()

	teams := []NBATeam{}
	for rows.Next() {
		var id int
		var name, abbrev string
		var isFavorited bool
		if err := rows.Scan(&id, &name, &abbrev, &isFavorited); err != nil {
			return []NBATeam{}, err
		}
		teams = append(teams, NBATeam{Id: id, Name: name, Abbreviation: abbrev, IsFavorited: isFavorited})
	}
	return teams, nil
}
