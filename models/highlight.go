package models

import (
	db "github.com/vibhavp/i58-portal/database"
)

type Highlight struct {
	ID      uint   `gorm:"primary_key" json:"-"`
	MatchID uint   `json:"-"`
	URL     string `sql:"not null;unique" json:"-"`
}

func AddHighlight(logsID int, url string) error {
	var matchID uint
	err := db.DB.DB().QueryRow("SELECT id FROM matches WHERE logs_id = $1", logsID).
		Scan(&matchID)
	if err != nil {
		return err
	}

	return db.DB.Save(&Highlight{
		MatchID: matchID,
		URL:     url,
	}).Error
}
