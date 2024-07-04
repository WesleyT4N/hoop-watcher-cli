package hoop_watcher

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type BaseHandler struct {
	db *HoopWatcherDB
}

func NewBaseHandler(db *HoopWatcherDB) *BaseHandler {
	return &BaseHandler{db: db}
}

func handleDBError(w http.ResponseWriter, err error) {
	if err == sql.ErrNoRows {
		http.Error(w, "No results found", http.StatusNotFound)
		return
	}
	log.Printf("DB Error occurred: %v", err)
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (h *BaseHandler) GetRoot(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, map[string]string{"status": "ok"})
}

func (h *BaseHandler) GetTeams(w http.ResponseWriter, r *http.Request) {
	teams, err := h.db.GetAllTeams()
	if err != nil {
		handleDBError(w, err)
		return
	}
	writeJSON(w, teams)
}

func (h *BaseHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	abbrev := r.PathValue("abbrev")
	log.Printf("abbrev: %s", abbrev)
	team, err := h.db.GetTeamByAbbrev(abbrev)
	if err != nil {
		handleDBError(w, err)
		return
	}
	writeJSON(w, team)
}
