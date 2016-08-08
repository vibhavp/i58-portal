package controllers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/vibhavp/i58-portal/controllers/api/admin"
	"github.com/vibhavp/i58-portal/models"
)

var classes = map[string]string{
	"scout":        "Scout",
	"soldier":      "Soldier",
	"demoman":      "Demoman",
	"sniper":       "Sniper",
	"medic":        "Medic",
	"spy":          "Spy",
	"heavyweapons": "Heavy",
	"engineer":     "Engineer",
	"pyro":         "Pyro",
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	indexPage := template.Must(template.ParseFiles("views/index.html"))
	err := indexPage.Execute(w, map[string]interface{}{
		"classes":  classes,
		"matches":  models.GetAllMatches(),
		"loggedIn": admin.IsLoggedIn(r),
	})
	if err != nil {
		log.Println(err)
	}
}

func serveAdmin(w http.ResponseWriter, r *http.Request) {
	page := template.Must(template.ParseFiles("views/admin.html"))
	err := page.Execute(w, map[string]interface{}{
		"notLoggedIn": !admin.IsLoggedIn(r),
	})
	if err != nil {
		log.Println(err)
	}
}
