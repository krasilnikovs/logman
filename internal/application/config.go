package application

type Configuration struct {
	AppEnv string `env:"APP_ENV" env-default:"dev"`
}

type ServerConfiguration struct {
	Configuration

	Port string `env:"PORT" env-default:"8016"`
}
