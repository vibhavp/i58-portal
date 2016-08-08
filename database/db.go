package database

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var DB *gorm.DB

func Open() {
	var err error
	DB, err = gorm.Open("sqlite3", "logs.db")
	if err != nil {
		log.Fatal(err)
	}

	DB.DB().SetMaxOpenConns(1)
}
