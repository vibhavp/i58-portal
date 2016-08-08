package models

import (
	db "github.com/vibhavp/i58-portal/database"
)

type Player struct {
	ID      uint   `gorm:"primary_key" json:"id"`
	Name    string `json:"name"`
	SteamID string `sql:"not null;unique" json:"-"`
	Team    string
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
