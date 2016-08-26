package models

import (
	"time"

	db "github.com/vibhavp/i58-portal/database"
)

type Tweet struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time

	URL string `sql:"not null;unique"`
}

func AddTweet(url string) error {
	return db.DB.Create(&Tweet{URL: url}).Error
}

func GetAllTweets() []Tweet {
	var tweets []Tweet
	db.DB.Order("id desc").Limit(5).Find(&tweets)
	return tweets
}
