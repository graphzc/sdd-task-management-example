package config

type CORS struct {
	AllowOrigins []string `env:"ALLOW_ORIGINS" envSeparator:","`
}
