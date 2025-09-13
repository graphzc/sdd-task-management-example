package config

type JWT struct {
	AccessTokenSecret     string `env:"ACCESS_TOKEN_SECRET"`
	AccessTokenExpiration string `env:"ACCESS_TOKEN_EXPIRATION_TIME"`
}
