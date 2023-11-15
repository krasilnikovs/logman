package application

// A Configuration contains application config which must contains every application(cli, api)
type Configuration struct {
	// AppEnv contains current environment the value gets from "LOGMAN_ENV" environment variable
	// if the environment variable is not set then by default will use development environment
	AppEnv string `env:"LOGMAN_ENV" env-default:"dev"`

	// DataStoragePath contains path to sqlite database, by default the value is var/data/logman.db
	// to override the path need to set env variable "LOGMAN_DB_PATH"
	DataStoragePath string `env:"LOGMAN_DB_PATH" env-default:"var/data/logman.db"`
}

// A ApiServerConfiguration contains application config related to api server
type ApiServerConfiguration struct {
	Configuration

	// Port shows on which ports works api server, the value reads from "LOGMAN_PORT" environment variable
	// if the environment variable is not set then will be use "8016" port
	Port string `env:"LOGMAN_PORT" env-default:"8016"`
}
