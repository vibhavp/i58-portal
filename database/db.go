package database

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/vibhavp/i58-portal/config"
)

var DB *gorm.DB

func Open() {
	var err error
	DB, err = gorm.Open(config.Config.DBDriver, config.Config.DBURL)
	if err != nil {
		log.Fatal(err)
	}

	DB.DB().SetMaxOpenConns(1)
}
