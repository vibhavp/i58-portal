package controllers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/vibhavp/i58-portal/controllers/api/admin"
	"github.com/vibhavp/i58-portal/models"
)

type stat struct {
	Stat   interface{}
	Player models.Player
}

type highest struct {
	Name            string
	HighestDPM      stat
	HighestKD       stat
	HighestAirshots stat
}

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

var landing = false

func newClassMap() map[string]highest {
	stats := make(map[string]highest, 9)

	for class, _ := range classes {
		dpm := models.GetHighestStat("dpm", class)
		kd := models.GetHighestStat("kd", class)
		var as models.AvgStats

		if class == "demoman" || class == "soldier" {
			as = *models.GetHighestStat("airshots", class)
		}

		stats[class] = highest{classes[class],
			stat{int(dpm.DPM), dpm.Player},
			stat{kd.KD, kd.Player},
			stat{as.Airshots, as.Player}}
	}
	return stats
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	indexPage := template.Must(template.ParseFiles("views/index.html"))
	err := indexPage.Execute(w, map[string]interface{}{
		"classes":  newClassMap(),
		"classMap": classes,
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
		"players":     models.GetAllPlayers(),
	})
	if err != nil {
		log.Println(err)
	}
}
