package admin

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/vibhavp/i58-portal/models"
)

func setWins(w http.ResponseWriter, r *http.Request) {
	if !IsLoggedIn(r) {
		http.Error(w, "You are not logged in!", http.StatusUnauthorized)
		return
	}

	if err := r.ParseForm(); err != nil {
		return
	}

	wins, err := strconv.Atoi(r.Form.Get("wins"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(r.Form.Get("teamid"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := models.SetWins(uint(id), wins); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, "Done!")
}

func setLosses(w http.ResponseWriter, r *http.Request) {
	if !IsLoggedIn(r) {
		http.Error(w, "You are not logged in!", http.StatusUnauthorized)
		return
	}

	if err := r.ParseForm(); err != nil {
		return
	}

	losses, err := strconv.Atoi(r.Form.Get("losses"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(r.Form.Get("teamid"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := models.SetLosses(uint(id), losses); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, "Done!")
}
