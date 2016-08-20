package models

import (
	db "github.com/vibhavp/i58-portal/database"
)

type Team struct {
	ID      uint     `gorm:"primary_key" json:"-"`
	Name    string   `sql:"not null;unique"`
	Players []Player `gorm:"ForeignKey:TeamID"`
	Wins    int
	Losses  int
}

func getTeam(name string) uint {
	team := Team{}
	db.DB.Where(&Team{Name: name}).FirstOrCreate(&team)
	return team.ID
}

func GetAllTeams() []Team {
	var teams []Team
	db.DB.Preload("Players").Find(&teams)
	return teams
}
