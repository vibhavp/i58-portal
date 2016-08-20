package models

import (
	db "github.com/vibhavp/i58-portal/database"
)

type Player struct {
	ID      uint   `gorm:"primary_key" json:"id"`
	Name    string `json:"name"`
	SteamID string `sql:"not null;unique" json:"-"`

	TeamID uint
	Team   Team `gorm:"ForeignKey:TeamID"`
}

func getOrCreatePlayerID(steamID string, names map[string]string) uint {
	var id uint
	err := db.DB.Model(&Player{}).Select("id").
		Where("steam_id = ?", steamID).
		Row().Scan(&id)
	if err != nil {
		db.DB.Save(&Player{
			Name:    names[steamID],
			SteamID: steamID,
		})
		db.DB.Model(&Player{}).Select("id").
			Where("steam_id = ?", steamID).
			Row().Scan(&id)
	}

	return id
}

func SetPlayerInfo(steamID, name, team string) error {
	return db.DB.Model(&Player{}).Where("steam_id = ?", steamID).
		Update(map[string]interface{}{
			"name":     name,
			"team_id":  getTeam(team),
			"steam_id": steamID,
		}).Error
}

func GetAllPlayers() []Player {
	var players []Player
	db.DB.Preload("Team").Find(&players)
	return players
}
