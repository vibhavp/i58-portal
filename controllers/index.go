package controllers

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/vibhavp/i58-portal/config"
	"github.com/vibhavp/i58-portal/controllers/api/admin"
	"github.com/vibhavp/i58-portal/models"
)

type stat struct {
	Stat   interface{}
	Player models.Player
}

type highest struct {
	ClassName           string
	Class               string
	HighestDPM          stat
	HighestUbersPerDrop stat
	HighestAirshots     stat
	HighestKills        stat
}

var classes = []string{"scout", "soldier", "demoman", "sniper", "medic", "spy",
	"pyro", "heavyweapons", "engineer"}

var classMap = map[string]string{
	"scout":        "Scout",
	"soldier":      "Soldier",
	"demoman":      "Demoman",
	"sniper":       "Sniper",
	"medic":        "Medic",
	"spy":          "Spy",
	"pyro":         "Pyro",
	"heavyweapons": "Heavy",
	"engineer":     "Engineer",
}

var landing = false

func newClassMap() []highest {
	var stats []highest

	for _, class := range classes {
		dpm := models.GetHighestStat("dpm", class)
		var ubersPerDrops models.AvgStats
		if class == "medic" {
			ubersPerDrops = *models.GetHighestStat("ubers_per_drop", "medic")
		}

		//var as models.AvgStats

		// if class == "demoman" || class == "soldier" {
		// 	as = *models.GetHighestStat("airshots", class)
		// }

		var as stat
		if class == "demoman" || class == "soldier" {
			s := models.GetHighestStat("airshots", class)
			as.Stat = int(s.Airshots)
			as.Player = s.Player
		}

		var kills stat
		if class == "sniper" || class == "spy" {
			k := models.GetHighestStat("kills", class)
			kills.Stat = int(k.Kills)
			kills.Player = k.Player
		}

		stats = append(stats, highest{
			classMap[class],
			class,
			stat{int(dpm.DPM), dpm.Player},
			stat{int(ubersPerDrops.UbersPerDrop), ubersPerDrops.Player},
			as,
			kills})
	}

	return stats
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	indexPage := template.Must(template.ParseFiles("views/index.html"))
	eventTime, e := time.Parse("2006-01-02T15:04:05-07:00", "2016-08-26T16:00:00+01:00")
	if e != nil {
		log.Println(e)
	}

	err := indexPage.Execute(w, map[string]interface{}{
		"classes":     newClassMap(),
		"matches":     models.GetAllMatches(),
		"loggedIn":    admin.IsLoggedIn(r),
		"teams":       models.GetAllTeams(),
		"beforeEvent": time.Now().Before(eventTime),
		"tweets":      models.GetAllTweets(),
		"tracking":    config.Config.AnalyticsID,
		"highlights":  models.GetAllHighlights(),
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
		"teamName":   models.GetTeamName(uint(teamID)),
		"matches":    models.GetMatches(uint(teamID)),
		"players":    models.GetTeamPlayers(uint(teamID)),
	})
	if err != nil {
		log.Println(err)
	}
}
