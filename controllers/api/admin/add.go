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
	if !isLoggedIn(r) {
		http.Error(w, "You are not logged in!", http.StatusUnauthorized)
		return
	}

	url := r.URL.Query().Get("url")
	if url == "" || !reValidURL.MatchString(url) {
		http.Error(w, "Missing/Invalid logs URL", http.StatusBadRequest)
		return
	}

	title := r.URL.Query().Get("title")
	if title == "" {
		http.Error(w, "Missing title", http.StatusBadRequest)
		return
	}

	page := r.URL.Query().Get("page")
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

	err = models.AddMatch(id, title, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Added!")
}

func addHighlight(w http.ResponseWriter, r *http.Request) {
	if !isLoggedIn(r) {
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

	err = models.AddHighlight(id, highlight)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
