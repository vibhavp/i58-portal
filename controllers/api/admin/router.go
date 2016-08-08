package admin

import (
	"net/http"
)

func init() {
	http.HandleFunc("/api/admin/add_match", addMatch)
	http.HandleFunc("/api/admin/add_highlight", addHighlight)
}
