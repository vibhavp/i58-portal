package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/vibhavp/i58-portal/models"
)

func getMatches(w http.ResponseWriter, r *http.Request) {
	matches := models.GetAllMatches()
	enc := json.NewEncoder(w)
	enc.Encode(matches)
}
