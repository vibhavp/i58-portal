package config

import (
	"log"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/kelseyhightower/envconfig"
)

var Config = struct {
	Username string `envconfig:"USERNAME" default:"i58"`
	Password string `envconfig:"PASSWORD" default:"i58"`
	Address  string `envconfig:"ADDRESS" default:":8080"`
	DBDriver string `envconfig:"DB_DRIVER" default:"sqlite3"`

	DBURL       string `envconfig:"DB_URL" default:"logs.db"`
	AnalyticsID string `envconfig:"ANALYTICS_ID"`
}{}

func init() {
	log.SetFlags(log.Lshortfile)
	err := envconfig.Process("PORTAL", &Config)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("username: %s, password %s", Config.Username, Config.Password)
}
