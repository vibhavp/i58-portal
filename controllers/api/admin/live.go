package admin

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/vibhavp/i58-portal/models"
)

var errFmt = `{"error": "%v"}`

func setLive(w http.ResponseWriter, r *http.Request) {
	if !IsLoggedIn(r) {
		http.Error(w, fmt.Sprintf(errFmt, "you are not logged in!"), http.StatusUnauthorized)
		return
	}

	err := r.ParseForm()
	if err != nil {
		return
	}
	query := r.Form

	matchID, err := strconv.Atoi(query.Get("match_id"))
	if err != nil {
		http.Error(w, fmt.Sprintf(errFmt, err), http.StatusBadRequest)
		return
	}

	if err := models.SetMatchLive(matchID); err != nil {
		http.Error(w, fmt.Sprintf(errFmt, err), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func unsetLive(w http.ResponseWriter, r *http.Request) {
	if !IsLoggedIn(r) {
		http.Error(w, fmt.Sprintf(errFmt, "you are not logged in!"), http.StatusUnauthorized)
		return
	}

	err := r.ParseForm()
	if err != nil {
		return
	}
	query := r.Form

	matchID, err := strconv.Atoi(query.Get("match_id"))
	if err != nil {
		http.Error(w, fmt.Sprintf(errFmt, err), http.StatusBadRequest)
		return
	}

	if err := models.UnsetMatchLive(matchID); err != nil {
		http.Error(w, fmt.Sprintf(errFmt, err), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
