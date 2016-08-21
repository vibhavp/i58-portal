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

func SetWins(teamID uint, wins int) error {
	return db.DB.Model(&Team{}).
		Where("id = ?", teamID).
		Update("wins", wins).Error
}

func SetLosses(teamID uint, losses int) error {
	return db.DB.Model(&Team{}).
		Where("id = ?", teamID).
		Update("losses", losses).Error
}

func GetMatches(teamID uint) []Match {
	var matches []Match
	db.DB.Preload("Team1").Preload("Team2").
		Where("team1_id = ? OR team2_id = ?", teamID, teamID).
		Find(&matches)
	return matches
}
