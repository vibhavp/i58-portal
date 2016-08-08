package stats

import (
	"net/http"
)

func init() {
	http.HandleFunc("/api/stats/get_stats", getStats)
	http.HandleFunc("/api/stats/get_player_stats", getPlayerStats)
	http.HandleFunc("/api/stats/get_all_class_stats", getAllClassStats)
}
