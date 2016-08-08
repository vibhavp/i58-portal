package models

import (
	"log"

	db "github.com/vibhavp/i58-portal/database"
)

type AllClassStat struct {
	ID     uint `json:"-"`
	LogsID uint `json:"-"`

	DamagePerHeal float64 `json:"damage_per_heal"`

	PlayerID uint   `json:"-"`
	Player   Player `gorm:"ForeignKey:PlayerID" json:"player"`
}
type Stat struct {
	ID     uint `json:"-"`
	LogsID uint `json:"-"`

	Class     string  `json:"class"`
	DPM       float64 `json:"dpm"`
	Kills     int     `json:"kills"`
	Deaths    int     `json:"deaths"`
	KD        float64 `json:"kd"`
	TotalTime int     `json:"-"`
	Drops     int     `json:"drops"`
	Airshots  int     `json:"airshots,omitempty"`

	PlayerID uint   `json:"-"`
	Player   Player `gorm:"ForeignKey:PlayerID" json:"player"`
}

type AvgStats struct {
	ID uint `json:"-"`

	Class    string  `json:"class"`
	DPM      float64 `json:"dpm"`
	Kills    int     `json:"kills"`
	Deaths   int     `json:"deaths"`
	KD       float64 `json:"kd"`
	Drops    float64 `json:"drops"`
	Airshots float64 `json:"airshots,omitempty"`

	PlayerID uint   `json:"-"`
	Player   Player `gorm:"ForeignKey:PlayerID" json:"player"`
}

func addPlayers(names map[string]string) {
	for steamID, name := range names {
		var count uint
		db.DB.Model(&Player{}).Where("steam_id = ?", steamID).Count(&count)
		if count == 0 {
			db.DB.Save(&Player{
				Name:    name,
				SteamID: steamID,
			})
		}
	}
}

var classes = []string{"scout", "soldier", "pyro", "demoman", "heavyweapons",
	"medic", "spy", "sniper", "engineer"}

func UpdateAvgStats(playerIDs []uint) {
	for _, id := range playerIDs {
		for _, class := range classes {
			var final AvgStats
			final.PlayerID = id
			var stats []Stat

			err := db.DB.Model(&Stat{}).Where("player_id = ? AND class = ?", id, class).Find(&stats).Error
			if err != nil {
				log.Println(err)
			}

			if len(stats) == 0 {
				continue
			}

			for _, stat := range stats {
				final.DPM += stat.DPM / float64(len(stats))
				final.Class = stat.Class
				final.Kills += stat.Kills / len(stats)
				final.Deaths += stat.Deaths / len(stats)
				final.KD += stat.KD / float64(len(stats))
				final.Airshots += float64(stat.Airshots) / float64(len(stats))
				final.Drops += float64(stat.Drops) / float64(len(stats))
			}

			db.DB.Where("player_id = ? AND class = ?", id, class).Delete(&AvgStats{})
			db.DB.Save(&final)
		}

	}
}

func GetClassStats(class string) []AvgStats {
	var stats []AvgStats
	db.DB.Preload("Player").Where("class = ?", class).Find(&stats)
	return stats
}

func GetPlayerStats(playerID uint) []Stat {
	var stats []Stat
	db.DB.Preload("Player").Where("player_id = ?", playerID).Order("logs_id").Find(&stats)
	for i := 1; i < len(stats); i++ {
		if stats[i].LogsID == stats[i-1].LogsID {
			stats[i-1].Class += "+" + stats[i].Class
			stats[i-1].DPM = (stats[i-1].DPM + stats[i].DPM) / 2.0
			stats[i-1].Kills += stats[i].Kills
			stats[i-1].Deaths += stats[i].Deaths
			stats[i-1].KD += float64(stats[i-1].Kills) / float64(stats[i-1].Deaths)
			stats[i-1].Drops += stats[i].Drops
			stats[i-1].Airshots += stats[i].Airshots

			stats = append(stats[:i], stats[i+1:]...)
		}
	}

	return stats
}

func GetAllPlayers() []Player {
	var players []Player
	db.DB.Find(&players)
	return players
}

func GetAllClassStats() []AllClassStat {
	var stats []AllClassStat
	players := GetAllPlayers()
	for _, player := range players {
		var stat AllClassStat
		var playerID uint
		db.DB.Model(&AllClassStat{}).
			Select("AVG(damage_per_heal), player_id").
			Where("player_id = ?", player.ID).Row().
			Scan(&stat.DamagePerHeal, &playerID)
		db.DB.First(&stat.Player, playerID)
		stats = append(stats, stat)
	}

	return stats
}
