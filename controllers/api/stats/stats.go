package stats

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/vibhavp/i58-portal/models"
)

func getStats(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	class := query.Get("class")
	if class == "" {
		http.Error(w, `{"error":"no class given"}`, http.StatusBadRequest)
		return
	}
	e := json.NewEncoder(w)
	e.Encode(map[string]interface{}{
		"type": "class",
		"data": models.GetClassStats(class),
	})
}

func getPlayerStats(w http.ResponseWriter, r *http.Request) {
	playerID, err := strconv.ParseUint(r.URL.Query().Get("playerid"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	e := json.NewEncoder(w)
	e.Encode(map[string]interface{}{
		"type": "player",
		"data": models.GetPlayerStats(uint(playerID)),
	})
}

func getAllClassStats(w http.ResponseWriter, r *http.Request) {
	e := json.NewEncoder(w)
	e.Encode(map[string]interface{}{
		"type": "allclass",
		"data": models.GetAllClassStats(),
	})
}
