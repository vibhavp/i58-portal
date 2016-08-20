package models

import (
	"github.com/TF2Stadium/logstf"
	db "github.com/vibhavp/i58-portal/database"
)

type Match struct {
	ID     uint `gorm:"primary_key" json:"-"`
	LogsID int  `sql:"not null;unique"`

	Team1ID uint `sql:"not null"`
	Team1   Team `gorm:"ForeignKey:Team1ID"`
	Team2ID uint `sql:"not null"`
	Team2   Team `gorm:"ForeignKey:Team2ID"`

	Stage      string `sql:"not null"`
	MatchPage  string `sql:"not null"`
	Live       bool
	Highlights []Highlight
}

func GetAllMatches() []Match {
	var matches []Match
	db.DB.Preload("Highlights").Preload("Team1").Preload("Team2").
		Find(&matches)
	return matches
}

func exists(id int) bool {
	var count uint
	db.DB.DB().QueryRow("SELECT COUNT(*) FROM matches WHERE logs_id = ?", id).
		Scan(&count)
	return count != 0
}

func SetMatchLive(matchID int) error {
	return db.DB.Table("matches").Where("id = ?", matchID).
		Update("live", true).Error
}

func UnsetMatchLive(matchID int) error {
	return db.DB.Table("matches").Where("id = ?", matchID).
		Update("live", false).Error
}

func AddMatch(logsID int, team1, team2, stage, page string) error {
	team1ID, team2ID := getTeam(team1), getTeam(team2)

	if exists(logsID) {
		return db.DB.Model(&Match{}).Where("logs_id = ?", logsID).
			Update(map[string]interface{}{
				"team1_id":   team1ID,
				"team2_id":   team2ID,
				"stage":      stage,
				"match_page": page,
			}).Error
	}

	err := db.DB.Create(&Match{
		LogsID:    logsID,
		Team1ID:   team1ID,
		Team2ID:   team2ID,
		Stage:     stage,
		MatchPage: page,
	}).Error
	if err != nil {
		return err
	}

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

			var kd, ud float64
			if cstats.Deaths == 0 {
				kd = float64(cstats.Kills)
			} else {
				kd = float64(cstats.Kills) / float64(cstats.Deaths)
			}

			if cstats.Type == "medic" && cstats.TotalTime > 60*10 {
				if stats.Drops == 0 {
					ud = float64(stats.Ubers)
				} else {
					ud = float64(stats.Ubers) / float64(stats.Drops)
				}
			}

			min := float64(cstats.TotalTime) / 60.0 // minutes played
			stat := &Stat{
				Class:        cstats.Type,
				LogsID:       uint(logsID),
				DPM:          float64(cstats.Damage) / min,
				Kills:        cstats.Kills,
				Deaths:       cstats.Deaths,
				KD:           kd,
				UbersPerDrop: ud,
				Drops:        stats.Drops,
				TotalTime:    cstats.TotalTime,
				PlayerID:     id,
			}
			if cstats.Type == "soldier" || cstats.Type == "demoman" {
				stat.Airshots = stats.Airshots
			}

			db.DB.Save(stat)
			stat.addAvg()
		}

	}

	return nil
}
