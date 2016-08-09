package models

import (
	"database/sql"

	db "github.com/vibhavp/i58-portal/database"
)

type Highlight struct {
	ID      uint   `gorm:"primary_key" json:"-"`
	MatchID uint   `json:"-"`
	Title   string `sql:"not null" json:"title"`
	URL     string `sql:"not null;unique"`
}

func AddHighlight(logsID int, url, title string) error {
	var matchID uint
	err := db.DB.DB().QueryRow("SELECT id FROM matches WHERE logs_id = $1", logsID).
		Scan(&matchID)
	if err != nil {
		if err == sql.ErrNoRows {
			AddMatch(logsID, "", "", "", "")
		} else {
			return err
		}
	}

	return db.DB.Save(&Highlight{
		MatchID: matchID,
		Title:   title,
		URL:     url,
	}).Error
}
