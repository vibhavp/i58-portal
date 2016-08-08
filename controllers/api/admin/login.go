package admin

import (
	"net/http"

	"github.com/vibhavp/i58-portal/config"
	"github.com/vibhavp/i58-portal/controllers/api"
)

func isLoggedIn(r *http.Request) bool {
	if cookie, err := r.Cookie("login"); err == nil {
		return cookie.Value ==
			api.LoginHash(config.Config.Username, config.Config.Password)
	}

	return false
}
