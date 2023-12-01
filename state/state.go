package state

import (
	"context"
	"os"

	"Popkat/config"

	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	"github.com/infinitybotlist/eureka/genconfig"
	"github.com/infinitybotlist/eureka/snippets"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

var (
	Pool      *pgxpool.Pool
	Redis     *redis.Client
	Spaces    *minio.Client
	Logger    *zap.Logger
	Context   = context.Background()
	Validator = validator.New()

	Config *config.Config
)

func Setup() {
	// Validator
	Validator.RegisterValidation("notblank", validators.NotBlank)
	Validator.RegisterValidation("nospaces", snippets.ValidatorNoSpaces)
	Validator.RegisterValidation("https", snippets.ValidatorIsHttps)
	Validator.RegisterValidation("httporhttps", snippets.ValidatorIsHttpOrHttps)

	// Config
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

	// Logger
	Logger = snippets.CreateZap()

	// PSQL
	Pool, err = pgxpool.New(Context, Config.Meta.PostgresURL)

	if err != nil {
		panic(err)
	}

	// Redis
	rOptions, err := redis.ParseURL(Config.Meta.RedisURL.Parse())

	if err != nil {
		panic(err)
	}

	Redis = redis.NewClient(rOptions)

	// DO Spaces
	DO, err := minio.New(Config.Meta.Spaces.Endpoint, Config.Meta.Spaces.AccessKey, Config.Meta.Spaces.AccessSecret, Config.Meta.Spaces.SSL)
	if err != nil {
		panic(err)
	}

	Spaces = DO
}
