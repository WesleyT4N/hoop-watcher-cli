package hoop_watcher

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const initDB = `
CREATE TABLE IF NOT EXISTS teams(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    abbrev TEXT NOT NULL,
    is_favorited BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS games(
    id SERIAL PRIMARY KEY,
    home_team_id INTEGER NOT NULL REFERENCES teams(id),
    away_team_id INTEGER NOT NULL REFERENCES teams(id),
    date DATE NOT NULL
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
