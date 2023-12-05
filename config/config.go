package config

import (
	_ "embed"
	"strings"
)

const (
	CurrentEnvProd    = "prod"
	CurrentEnvStaging = "staging"
)

//go:embed current-env
var CurrentEnv string

func init() {
	CurrentEnv = strings.TrimSpace(CurrentEnv)

	if CurrentEnv != CurrentEnvProd && CurrentEnv != CurrentEnvStaging {
		panic("invalid environment")
	}
}

type Differs[T any] struct {
	Staging T `yaml:"staging" comment:"Staging value" validate:"required"`
	Prod    T `yaml:"prod" comment:"Production value" validate:"required"`
}

func (d *Differs[T]) Parse() T {
	if CurrentEnv == CurrentEnvProd {
		return d.Prod
	} else if CurrentEnv == CurrentEnvStaging {
		return d.Staging
	} else {
		panic("invalid environment")
	}
}

func (d *Differs[T]) Production() T {
	return d.Prod
}

type Config struct {
	Meta  Meta  `yaml:"meta" validate:"required"`
	Sites Sites `yaml:"sites" validate:"required"`
}

type Sites struct {
	Frontend Differs[string] `yaml:"frontend" default:"https://popkat.select-list.xyz" comment:"Frontend URL" validate:"required"`
	API      Differs[string] `yaml:"api" default:"https://popkatapi.select-list.xyz" comment:"API URL" validate:"required"`
	Instatus string          `yaml:"instatus" default:"https://status.select-list.xyz" comment:"Instatus Status Page URL" validate:"required"`
}

type Meta struct {
	PostgresURL string          `yaml:"postgresql" default:"postgresql://postgres@0.0.0.0" comment:"PostgreSQL URL" validate:"required"`
	RedisURL    Differs[string] `yaml:"redis" default:"redis://0.0.0.0/1" comment:"Redis URL" validate:"required"`
	Port        Differs[string] `yaml:"port" default:":8081" comment:"Port to run the server on" validate:"required"`
	Spaces      Spaces          `yaml:"spaces" validate:"required"`
}

type Spaces struct {
	Endpoint     string `yaml:"endpoint" comment:"Digital Ocean - Origin Endpoint URL" validate:"required"`
	AccessKey    string `yaml:"access_key" comment:"Digital Ocean - Access Key" validate:"required"`
	AccessSecret string `yaml:"access_secret" comment:"Digital Ocean - Access Secret" validate:"required"`
	SSL          bool   `yaml:"ssl" comment:"Enable SSL?" validate:"required"`
}
