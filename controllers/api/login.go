package api

import (
	"crypto/sha512"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/vibhavp/i58-portal/config"
)

func Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := r.PostForm
	pass := query.Get("password")
	username := query.Get("username")

	if username != config.Config.Username || pass != config.Config.Password {
		http.Error(w, "Invalid username/password", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "login",
		Path:    "/",
		Value:   LoginHash(config.Config.Username, config.Config.Password),
		Expires: time.Now().Add(30 * 24 * time.Hour),
	})

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func LoginHash(uname, pwd string) string {
	hash := sha512.Sum512([]byte(uname + ":" + pwd))
	return hex.EncodeToString(hash[:])
}
