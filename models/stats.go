package models

import (
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
	ID     uint `json:"-" gorm:"primary_key"`
	LogsID uint `json:"-"`

	Class        string  `json:"class" sql:"not null"`
	DPM          float64 `json:"dpm"`
	Kills        int     `json:"kills"`
	Deaths       int     `json:"deaths"`
	KD           float64 `json:"kd"`
	TotalTime    int     `json:"-"`
	Drops        int     `json:"drops"`
	Ubers        int     `json:"ubers"`
	UbersPerDrop float64 `json:"ubers_per_drops"`
	Airshots     int     `json:"airshots,omitempty"`

	PlayerID uint   `json:"-"`
	Player   Player `gorm:"ForeignKey:PlayerID" json:"player" sql:"not null"`
}

type AvgStats struct {
	ID uint `json:"-" gorm:"primary_key"`

	Class        string  `json:"class" sql:"not null"`
	DPM          float64 `json:"dpm"`
	Kills        int     `json:"kills"`
	Deaths       int     `json:"deaths"`
	KD           float64 `json:"kd"`
	Ubers        float64 `json:"ubers"`
	Drops        float64 `json:"drops"`
	UbersPerDrop float64 `json:"ubers_per_drops"`
	Airshots     float64 `json:"airshots,omitempty"`

	PlayerID uint   `json:"-"`
	Player   Player `gorm:"ForeignKey:PlayerID" json:"player" sql:"not null"`
}

func GetHighestStat(stat, class string) *AvgStats {
	s := &AvgStats{}
	db.DB.Preload("Player").Where("class = ?", class).
		Order(stat + " DESC").First(s)
	return s
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

func (stat *Stat) addAvg() {
	avg := &AvgStats{}
	n := 0

	err := db.DB.Where("player_id = ? AND class = ?", stat.PlayerID, stat.Class).
		First(&avg).Error
	if err == nil {
		db.DB.Model(&Stat{}).Where("player_id = ? AND class = ?",
			stat.PlayerID,
			stat.Class).
			Count(&n)
	} else {
		avg.PlayerID = stat.PlayerID
		avg.Class = stat.Class
	}

	setCumAvg(stat.DPM, &avg.DPM, n)
	setCumAvgInt(stat.Kills, &avg.Kills, n)
	setCumAvgInt(stat.Deaths, &avg.Deaths, n)
	setCumAvg(stat.KD, &avg.KD, n)
	setCumAvgIntFloat(stat.Ubers, &avg.Ubers, n)
	setCumAvgIntFloat(stat.Drops, &avg.Drops, n)
	setCumAvg(stat.UbersPerDrop, &avg.UbersPerDrop, n)
	setCumAvgIntFloat(stat.Airshots, &avg.Airshots, n)
	db.DB.Save(avg)
}

func setCumAvg(newPoint float64, avg *float64, n int) {
	*avg = (newPoint + float64(n)*(*avg)) / float64(n+1)
}

func setCumAvgInt(newPoint int, avg *int, n int) {
	*avg = (newPoint + n*(*avg)) / (n + 1)
}

func setCumAvgIntFloat(newPoint int, avg *float64, n int) {
	*avg = (float64(newPoint) + float64(n)*(*avg)) / float64(n+1)
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
			stats[i-1].KD += float64(stats[i-1].Kills+stats[i].Kills) / float64(stats[i-1].Deaths+stats[i].Deaths)
			stats[i-1].Drops += stats[i].Drops
			stats[i-1].UbersPerDrop += stats[i].UbersPerDrop
			stats[i-1].Airshots += stats[i].Airshots

			stats = append(stats[:i], stats[i+1:]...)
		}
	}

	return stats
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
