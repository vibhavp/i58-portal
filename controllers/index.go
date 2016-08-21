package controllers

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/vibhavp/i58-portal/controllers/api/admin"
	"github.com/vibhavp/i58-portal/models"
)

type stat struct {
	Stat   interface{}
	Player models.Player
}

type highest struct {
	Name                string
	HighestDPM          stat
	HighestUbersPerDrop stat
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
		var ubersPerDrops models.AvgStats
		if class == "medic" {
			ubersPerDrops = *models.GetHighestStat("ubers_per_drop", "medic")
		}

		//var as models.AvgStats

		// if class == "demoman" || class == "soldier" {
		// 	as = *models.GetHighestStat("airshots", class)
		// }

		stats[class] = highest{classes[class],
			stat{int(dpm.DPM), dpm.Player},
			stat{int(ubersPerDrops.UbersPerDrop), ubersPerDrops.Player}}
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
		"teams":    models.GetAllTeams(),
	})
	if err != nil {
		log.Println(err)
	}
}

func serveAdmin(w http.ResponseWriter, r *http.Request) {
	loggedIn := admin.IsLoggedIn(r)
	var allPlayers *[]models.Player

	if loggedIn {
		p := models.GetAllPlayers()
		allPlayers = &p
	}

	page := template.Must(template.ParseFiles("views/admin.html"))
	err := page.Execute(w, map[string]interface{}{
		"notLoggedIn": !loggedIn,
		"players":     allPlayers,
	})
	if err != nil {
		log.Println(err)
	}
}

func serveTeam(w http.ResponseWriter, r *http.Request) {
	teamPage := template.Must(template.ParseFiles("views/team.html"))
	query := r.URL.Query()
	teamID, err := strconv.Atoi(query.Get("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = teamPage.Execute(w, map[string]interface{}{
		"selfTeamID": teamID,
		"matches":    models.GetMatches(uint(teamID)),
	})
	if err != nil {
		log.Println(err)
	}
}
