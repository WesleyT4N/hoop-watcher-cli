package hoop_watcher

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func TestGetTeams(t *testing.T) {
	t.Run("returns all teams", func(t *testing.T) {
		db := newMockDB()
		db.getAllTeams = func() ([]NBATeam, error) {
			return []NBATeam{
				{
					Id:           1,
					Name:         "Atlanta Hawks",
					Abbreviation: "ATL",
				},
				{
					Id:           2,
					Name:         "Boston Celtics",
					Abbreviation: "BOS",
				},
			}, nil
		}
		req, _ := http.NewRequest("GET", "/teams", nil)
		rr := httptest.NewRecorder()
		h := NewBaseHandler(db)
		h.GetTeams(rr, req)

		var got []NBATeam
		json.Unmarshal(rr.Body.Bytes(), &got)
		want := []NBATeam{
			{
				Id:           1,
				Name:         "Atlanta Hawks",
				Abbreviation: "ATL",
			},
			{
				Id:           2,
				Name:         "Boston Celtics",
				Abbreviation: "BOS",
			},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
	t.Run("404 if no teams found", func(t *testing.T) {
		db := newMockDB()
		db.getAllTeams = func() ([]NBATeam, error) {
			return []NBATeam{}, sql.ErrNoRows
		}
		req, _ := http.NewRequest("GET", "/teams", nil)
		rr := httptest.NewRecorder()
		h := NewBaseHandler(db)
		h.GetTeams(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("got %d, want %d", rr.Code, http.StatusNotFound)
		}

		if rr.Body.String() != "No results found\n" {
			t.Errorf("got %s, want %s", rr.Body.String(), "No results found\n")
		}
	})

	t.Run("500 if error occurs", func(t *testing.T) {
		log.SetOutput(io.Discard)
		defer func() {
			log.SetOutput(os.Stderr)
		}()

		db := newMockDB()
		db.getAllTeams = func() ([]NBATeam, error) {
			return []NBATeam{}, errors.New("unknown error")
		}
		req, _ := http.NewRequest("GET", "/teams", nil)
		rr := httptest.NewRecorder()
		h := NewBaseHandler(db)
		h.GetTeams(rr, req)

		if rr.Code != http.StatusInternalServerError {
			t.Errorf("got %d, want %d", rr.Code, http.StatusInternalServerError)
		}

		if rr.Body.String() != "Internal Server Error\n" {
			t.Errorf("got %s, want %s", rr.Body.String(), "Internal Server Error\n")
		}
	})
}

type mockHoopWatcherDB struct {
	getAllTeams       func() ([]NBATeam, error)
	getTeamByAbbrev   func(abbrev string) (NBATeam, error)
	setTeamFavorite   func(id int, fav bool) error
	getTeamHighlights func(id int) ([]Highlight, error)
}

func (m *mockHoopWatcherDB) GetAllTeams() ([]NBATeam, error) {
	return m.getAllTeams()
}

func (m *mockHoopWatcherDB) GetTeamByAbbrev(abbrev string) (NBATeam, error) {
	return m.getTeamByAbbrev(abbrev)
}

func (m *mockHoopWatcherDB) SetTeamFavorite(id int, fav bool) error {
	return m.setTeamFavorite(id, fav)
}

func (m *mockHoopWatcherDB) GetTeamHighlights(id int) ([]Highlight, error) {
	return m.getTeamHighlights(id)
}

func newMockDB() *mockHoopWatcherDB {
	return &mockHoopWatcherDB{
		getAllTeams: func() ([]NBATeam, error) {
			return []NBATeam{}, nil
		},
		getTeamByAbbrev: func(abbrev string) (NBATeam, error) {
			return NBATeam{}, nil
		},
		getTeamHighlights: func(id int) ([]Highlight, error) {
			return []Highlight{}, nil
		},
	}
}
