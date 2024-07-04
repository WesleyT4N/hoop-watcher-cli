package hoop_watcher

import (
	"fmt"
	"net/http"
)

func GetRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}
