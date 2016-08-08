package controllers

import (
	"html/template"
	"log"
	"net/http"

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
		"classes": classes,
		"matches": models.GetAllMatches(),
	})
	if err != nil {
		log.Println(err)
	}
}
