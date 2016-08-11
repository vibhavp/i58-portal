package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

var Config = struct {
	Username string `envconfig:"USERNAME" required:"true"`
	Password string `envconfig:"PASSWORD" required:"true"`
	Address  string `envconfig:"ADDRESS" default:":8080"`
}{}

func init() {
	envconfig.MustProcess("PORTAL", &Config)
	log.Printf("username: %s, password %s", Config.Username, Config.Password)
}
