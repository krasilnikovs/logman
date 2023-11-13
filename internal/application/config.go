// Package application represent infrastructure code like default logger
// configuration capabilities, server, cli and etc.
package application

// A Configuration contains application config which must contains every application(cli, api)
type Configuration struct {
	AppEnv string `env:"APP_ENV" env-default:"dev"`
}

// A ServerConfiguration contains application config related to api server
type ServerConfiguration struct {
	Configuration

	Port string `env:"PORT" env-default:"8016"`
}
