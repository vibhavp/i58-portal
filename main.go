package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/vibhavp/i58-portal/config"
	_ "github.com/vibhavp/i58-portal/controllers"
	"github.com/vibhavp/i58-portal/database"
	"github.com/vibhavp/i58-portal/models"
)

func init() {
	flag.Parse()
}

func main() {
	database.Open()
	models.Migrate()

	log.Printf("serving on %s", config.Config.Address)
	log.Fatal(http.ListenAndServe(config.Config.Address, nil))
}
