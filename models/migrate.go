package models

import (
	db "github.com/vibhavp/i58-portal/database"
)

func Migrate() {
	db.DB.AutoMigrate(&Highlight{})
	db.DB.AutoMigrate(&Match{})
	db.DB.AutoMigrate(&Player{})
	db.DB.AutoMigrate(&AllClassStat{})
	db.DB.AutoMigrate(&Stat{})
	db.DB.AutoMigrate(&AvgStats{})
}
