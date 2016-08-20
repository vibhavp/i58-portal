package admin

import (
	"net/http"
)

func init() {
	http.HandleFunc("/api/admin/add_match", addMatch)
	http.HandleFunc("/api/admin/add_highlight", addHighlight)
	http.HandleFunc("/api/admin/add_player", addPlayer)

	http.HandleFunc("/api/admin/set_live", setLive)
	http.HandleFunc("/api/admin/unset_live", unsetLive)

	http.HandleFunc("/api/admin/set_wins", setWins)
	http.HandleFunc("/api/admin/set_losses", setLosses)
}
