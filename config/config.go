package config

import (
	"github.com/kelseyhightower/envconfig"
)

var configuration Configuration

type Redis struct {
	Dsn        string   `envconfig:"dsn"`
	Topics     []string `envconfig:"topics"`
	MaxIdle    int      `envconfig:"max_idle" default:"10"`
	MaxTimeout int      `envconfig:"max_timeout" default:"240"`
}

type Db struct {
	Dsn         string `envconfig:"dsn"`
	MaxIdle     int    `envconfig:"max_idle" default:"5"`
	MaxOpen     int    `envconfig:"max_open" default:"10"`
	MaxLifetime int    `envconfig:"max_lifetime" default:"3"`
}

type Configuration struct {
	Env   string `envconfig:"env"`
	Redis Redis  `envconfig:"redis"`
	Db    Db     `envconfig:"db"`
}

func Get() Configuration {
	return configuration
}

func Init() {
	envconfig.MustProcess("MONSTURN", &configuration)
}
