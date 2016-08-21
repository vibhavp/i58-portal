package models

import (
	db "github.com/vibhavp/i58-portal/database"
)

type Player struct {
	ID      uint   `gorm:"primary_key" json:"id"`
	Name    string `json:"name"`
	SteamID string `sql:"not null;unique" json:"-"`

	AvgStats *AvgStats `gorm:"ForeignKey:PlayerID" json:"-"`
	TeamID   uint
	Team     Team `gorm:"ForeignKey:TeamID"`
}

func getOrCreatePlayerID(steamID string, name string, teamID uint) uint {
	var id uint
	err := db.DB.Model(&Player{}).Select("id").
		Where("steam_id = ?", steamID).
		Row().Scan(&id)
	if err != nil {
		db.DB.Create(&Player{
			Name:    name,
			TeamID:  teamID,
			SteamID: steamID,
		})
		db.DB.Model(&Player{}).Select("id").
			Where("steam_id = ?", steamID).
			Row().Scan(&id)
	}

	return id
}

func SetPlayerInfo(steamID, name string) error {
	return db.DB.Model(&Player{}).Where("steam_id = ?", steamID).
		Update(map[string]interface{}{
			"name":     name,
			"steam_id": steamID,
		}).Error
}

func GetTeamPlayers(teamID uint) []Player {
	var players []Player
	db.DB.Model(&Player{}).Preload("AvgStats").Where("team_id = ?", teamID).Find(&players)
	return players
}

func GetAllPlayers() []Player {
	var players []Player
	db.DB.Preload("Team").Find(&players)
	return players
}
