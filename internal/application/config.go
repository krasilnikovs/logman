// Package application represent infrastructure code like default logger
// configuration capabilities, server, cli and etc.
package application

// A Configuration contains application config which must contains every application(cli, api)
type Configuration struct {
	// AppEnv contains current environment the value gets from "LOGMAN_ENV" environment variable
	// if the environment variable is not set then by default will use development environment
	AppEnv string `env:"LOGMAN_ENV" env-default:"dev"`
}

// A ApiServerConfiguration contains application config related to api server
type ApiServerConfiguration struct {
	Configuration

	// Port shows on which ports works api server, the value reads from "LOGMAN_PORT" environment variable
	// if the environment variable is not set then will be use "8016" port
	Port string `env:"LOGMAN_PORT" env-default:"8016"`
}
