package admin

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/vibhavp/i58-portal/models"
)

var reValidURL = regexp.MustCompile(`logs.tf/(\d+)`)

func addMatch(w http.ResponseWriter, r *http.Request) {
	if !IsLoggedIn(r) {
		http.Error(w, "You are not logged in!", http.StatusUnauthorized)
		return
	}

	err := r.ParseForm()
	if err != nil {
		return
	}
	query := r.Form

	url := strings.TrimSpace(query.Get("url"))
	if url == "" || !reValidURL.MatchString(url) {
		http.Error(w, "Missing/Invalid logs URL", http.StatusBadRequest)
		return
	}

	team1 := strings.TrimSpace(query.Get("team1"))
	if team1 == "" {
		http.Error(w, "Missing team1", http.StatusBadRequest)
		return
	}

	team1Score, err := strconv.Atoi(query.Get("team1_score"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	team2 := strings.TrimSpace(query.Get("team2"))
	if team2 == "" {
		http.Error(w, "Missing team2", http.StatusBadRequest)
		return
	}

	team2Score, err := strconv.Atoi(query.Get("team2_score"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stage := strings.TrimSpace(query.Get("stage"))
	if stage == "" {
		http.Error(w, "Missing stage", http.StatusBadRequest)
		return
	}

	page := strings.TrimSpace(query.Get("page"))
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

	err = models.AddMatch(id, team1, team2, stage, page, team1Score, team2Score)
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

	err := r.ParseForm()
	if err != nil {
		return
	}
	query := r.Form

	url := query.Get("url")
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

	highlight := query.Get("highlight")
	if highlight == "" {
		http.Error(w, "Missing page", http.StatusBadRequest)
		return
	}

	title := query.Get("title")
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

func addPlayer(w http.ResponseWriter, r *http.Request) {
	if !IsLoggedIn(r) {
		http.Error(w, "You are not logged in!", http.StatusUnauthorized)
		return
	}

	err := r.ParseForm()
	if err != nil {
		return
	}
	query := r.Form

	steamID := strings.TrimSpace(query.Get("steamid"))
	if steamID == "" {
		http.Error(w, "Missing SteamID", http.StatusBadRequest)
		return
	}

	name := strings.TrimSpace(query.Get("name"))
	if name == "" {
		http.Error(w, "Missing Name", http.StatusBadRequest)
		return
	}

	team := strings.TrimSpace(query.Get("team"))
	if team == "" {
		http.Error(w, "Missing Team", http.StatusBadRequest)
		return
	}

	if err := models.SetPlayerInfo(steamID, name, team); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "Added!")
}
