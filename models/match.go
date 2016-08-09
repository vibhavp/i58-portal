package models

import (
	"github.com/TF2Stadium/logstf"
	db "github.com/vibhavp/i58-portal/database"
)

type Match struct {
	ID         uint   `gorm:"primary_key" json:"-"`
	LogsID     int    `sql:"not null;unique"`
	Team1      string `sql:"not null"`
	Team2      string `sql:"not null"`
	Stage      string `sql:"not null"`
	MatchPage  string `sql:"not null"`
	Highlights []Highlight
}

func GetAllMatches() []Match {
	var matches []Match
	db.DB.Preload("Highlights").Find(&matches)
	return matches
}

func exists(id int) bool {
	var count uint
	db.DB.DB().QueryRow("SELECT COUNT(*) FROM matches WHERE logs_id = ?", id).
		Scan(&count)
	return count != 0
}

func AddMatch(logsID int, team1, team2, stage, page string) error {
	if exists(logsID) {
		return db.DB.Model(&Match{}).Where("logs_id = ?", logsID).
			Update(map[string]string{
				"team1":      team1,
				"team2":      team2,
				"stage":      stage,
				"match_page": page,
			}).Error
	}

	db.DB.Save(&Match{
		LogsID:    logsID,
		Team1:     team1,
		Team2:     team2,
		Stage:     stage,
		MatchPage: page,
	})
	return addStats(logsID)
}

func addStats(logsID int) error {
	logs, err := logstf.GetLogs(logsID)
	if err != nil {
		return err
	}

	var updatedIDs []uint

	for steamID, stats := range logs.Players {
		id := getOrCreatePlayerID(steamID, logs.Names)
		updatedIDs = append(updatedIDs, id)

		for _, cstats := range stats.ClassStats {
			if cstats.TotalTime == 0 || cstats.TotalTime < 60*2 {
				continue
			}
			if cstats.Type == "undefined" || cstats.Type == "" {
				continue
			}

			var kd float64
			if cstats.Deaths == 0 {
				kd = float64(cstats.Kills)
			} else {
				kd = float64(cstats.Kills) / float64(cstats.Deaths)
			}

			min := float64(cstats.TotalTime) / 60.0 // minutes played
			stat := &Stat{
				Class:     cstats.Type,
				LogsID:    uint(logsID),
				DPM:       float64(cstats.Damage) / min,
				Kills:     cstats.Kills,
				Deaths:    cstats.Deaths,
				KD:        kd,
				Drops:     stats.Drops,
				TotalTime: cstats.TotalTime,
				PlayerID:  id,
			}
			if cstats.Type == "soldier" || cstats.Type == "demoman" {
				stat.Airshots = stats.Airshots
			}

			db.DB.Save(stat)
		}

	}

	UpdateAvgStats(updatedIDs)
	return nil
}
