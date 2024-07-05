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

type HoopWatcherDB interface {
	GetAllTeams() ([]NBATeam, error)
	GetTeamByAbbrev(abbrev string) (NBATeam, error)
	SetTeamFavorite(teamId int, favorite bool) error
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

func (h *SqliteHoopWatcherDB) addTeam(name, abbrev string) error {
	// prepared statement
	stmt, err := h.db.Prepare("INSERT OR IGNORE INTO teams(name, abbrev) VALUES(?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(name, abbrev)
	return err
}

func (h *SqliteHoopWatcherDB) addAllTeams(filePath string) error {
	teams := GetNBATeamsFromJSON(filePath)
	for _, team := range teams {
		if err := h.addTeam(team.Name, team.Abbreviation); err != nil {
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

func (h *SqliteHoopWatcherDB) GetTeamByAbbrev(abbrev string) (NBATeam, error) {
	var team NBATeam
	row := h.db.QueryRow("SELECT * FROM teams WHERE abbrev = ?", abbrev)
	if err := row.Scan(&team.Id, &team.Name, &team.Abbreviation, &team.IsFavorited); err != nil {
		return NBATeam{}, err
	}
	return team, nil
}

func (h *SqliteHoopWatcherDB) SetTeamFavorite(teamId int, favorite bool) error {
	stmt, err := h.db.Prepare("UPDATE teams SET is_favorited = ? WHERE id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(favorite, teamId)
	if err != nil {
		return err
	}
	return nil
}
