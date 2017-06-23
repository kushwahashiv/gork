package main

import (
	"github.com/urfave/cli"
)

var cliFlags = []cli.Flag{
	cli.StringFlag{
		Name:   "db-redis-hostname",
		Usage:  "Redis hostname.",
		EnvVar: envPrefix("DB_REDIS_HOSTNAME"),
		Value:  "db-redis-server.default.svc.cluster.local",
	},
	cli.StringFlag{
		Name:   "db-redis-port",
		Usage:  "Redis port number.",
		EnvVar: envPrefix("DB_REDIS_PORT"),
		Value:  "6379",
	},
	cli.IntFlag{
		Name:   "db-redis-database",
		Usage:  "Redis database number.",
		EnvVar: envPrefix("DB_REDIS_DATABASE"),
		Value:  0,
	},
	cli.StringFlag{
		Name:   "db-redis-password",
		Usage:  "Redis password.",
		EnvVar: envPrefix("DB_REDIS_PASSWORD"),
	},
	cli.StringFlag{
		Name:   "gtw-rest-hostname",
		Usage:  "REST gateway hostname.",
		EnvVar: envPrefix("GTW_REST_HOSTNAME"),
	},
	cli.StringFlag{
		Name:   "gtw-rest-port",
		Usage:  "REST gateway port.",
		EnvVar: envPrefix("GTW_REST_PORT"),
		Value:  "8080",
	},
	cli.BoolFlag{
		Name:   "debug, d",
		Usage:  "Enable debug mode.",
		EnvVar: envPrefix("MISC_DEBUG_MODE"),
	},
	cli.BoolFlag{
		Name:   "text, t",
		Usage:  "Output logs as text instead of JSON.",
		EnvVar: envPrefix("MISC_LOG_FORMAT_TEXT"),
	},
}

// envPrefix makes environment var name prefixed.
func envPrefix(name string) (result string) {
	result = "GORK_" + name
	return
}

// NewConfigFromCtx creates a new instance of config.
func NewConfigFromCtx(ctx *cli.Context) (cfg *config, err error) {

	cfg = &config{
		Db: &configDb{
			Redis: &configDbRedis{
				Hostname: ctx.String("db-redis-hostname"),
				Port:     ctx.String("db-redis-port"),
				Database: ctx.Int("db-redis-database"),
				Password: ctx.String("db-redis-password"),
			},
		},
		Gtw: &configGtw{
			Rest: &configGtwRest{
				Hostname: ctx.String("gtw-rest-hostname"),
				Port:     ctx.String("gtw-rest-port"),
			},
		},
		Misc: &configMisc{
			DebugMode:     ctx.Bool("debug"),
			LogFormatText: ctx.Bool("text"),
		},
	}

	err = cfg.Validate()
	if err != nil {
		return nil, err
	}

	return
}

// config represents application configuration store.
type config struct {
	Db   *configDb
	Gtw  *configGtw
	Misc *configMisc
}

// Validate is responsible for data validation.
func (c *config) Validate() (err error) {
	return validation.ValidateStruct(c,
		validation.Field(&c.Db, validation.Required),
		validation.Field(&c.Gtw, validation.Required),
		validation.Field(&c.Misc, validation.Required),
	)
}

// configDb represents databases configuration.
type configDb struct {
	Redis *configDbRedis
}

// Validate is responsible for data validation.
func (c *configDb) Validate() (err error) {
	return validation.ValidateStruct(c,
		validation.Field(&c.Redis, validation.Required),
	)
}

// configDbRedis represents Redis DB configuration.
type configDbRedis struct {
	Hostname string
	Port     string
	Database int
	Password string
}

// Validate is responsible for data validation.
func (c *configDbRedis) Validate() (err error) {
	return validation.ValidateStruct(c,
		validation.Field(&c.Hostname, validation.Required, is.Host),
		validation.Field(&c.Port, validation.Required, is.Port),
		validation.Field(&c.Database, validation.Min(0)),
	)
}

// configGtw represents gateways configuration.
type configGtw struct {
	Rest *configGtwRest
}

// Validate is responsible for data validation.
func (c *configGtw) Validate() (err error) {
	return validation.ValidateStruct(c,
		validation.Field(&c.Rest, validation.Required),
	)
}

// configGtwRest represents REST gateway configuration.
type configGtwRest struct {
	Hostname string
	Port     string
}

// Validate is responsible for data validation.
func (c *configGtwRest) Validate() (err error) {
	return validation.ValidateStruct(c,
		validation.Field(&c.Hostname, is.Host),
		validation.Field(&c.Port, validation.Required, is.Port),
	)
}

// configMisc represents other configuration options.
type configMisc struct {
	DebugMode     bool
	LogFormatText bool
}
