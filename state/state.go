package state

import (
	"context"
	"os"

	"Popkat/config"

	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	"github.com/infinitybotlist/eureka/genconfig"
	"github.com/infinitybotlist/eureka/snippets"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

var (
	Logger    *zap.Logger
	Context   = context.Background()
	Validator = validator.New()

	Config *config.Config
)

func Setup() {
	Validator.RegisterValidation("notblank", validators.NotBlank)
	Validator.RegisterValidation("nospaces", snippets.ValidatorNoSpaces)
	Validator.RegisterValidation("https", snippets.ValidatorIsHttps)
	Validator.RegisterValidation("httporhttps", snippets.ValidatorIsHttpOrHttps)

	genconfig.GenConfig(config.Config{})

	cfg, err := os.ReadFile("config.yaml")

	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(cfg, &Config)

	if err != nil {
		panic(err)
	}

	err = Validator.Struct(Config)

	if err != nil {
		panic("configError: " + err.Error())
	}

	Logger = snippets.CreateZap()
}
