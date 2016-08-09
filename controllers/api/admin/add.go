package admin

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/vibhavp/i58-portal/models"
)

var reValidURL = regexp.MustCompile(`logs.tf/(\d+)`)

func addMatch(w http.ResponseWriter, r *http.Request) {
	if !IsLoggedIn(r) {
		http.Error(w, "You are not logged in!", http.StatusUnauthorized)
		return
	}
	query := r.URL.Query()

	url := query.Get("url")
	if url == "" || !reValidURL.MatchString(url) {
		http.Error(w, "Missing/Invalid logs URL", http.StatusBadRequest)
		return
	}

	team1 := query.Get("team1")
	if team1 == "" {
		http.Error(w, "Missing team1", http.StatusBadRequest)
		return
	}

	team2 := query.Get("team2")
	if team2 == "" {
		http.Error(w, "Missing team2", http.StatusBadRequest)
		return
	}

	stage := query.Get("stage")
	if stage == "" {
		http.Error(w, "Missing stage", http.StatusBadRequest)
		return
	}

	page := query.Get("page")
	if page == "" {
		http.Error(w, "Missing page", http.StatusBadRequest)
		return
	}

	m := reValidURL.FindStringSubmatch(url)
	id, err := strconv.Atoi(m[1])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = models.AddMatch(id, team1, team2, stage, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Added!")
}

func addHighlight(w http.ResponseWriter, r *http.Request) {
	if !IsLoggedIn(r) {
		http.Error(w, "You are not logged in!", http.StatusUnauthorized)
		return
	}

	url := r.URL.Query().Get("url")
	if url == "" || !reValidURL.MatchString(url) {
		http.Error(w, "Missing/Invalid logs URL", http.StatusBadRequest)
		return
	}

	m := reValidURL.FindStringSubmatch(url)
	id, err := strconv.Atoi(m[1])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	highlight := r.URL.Query().Get("highlight")
	if highlight == "" {
		http.Error(w, "Missing page", http.StatusBadRequest)
		return
	}

	title := r.URL.Query().Get("title")
	if title == "" {
		http.Error(w, "Missing title", http.StatusBadRequest)
		return
	}

	err = models.AddHighlight(id, highlight, title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
