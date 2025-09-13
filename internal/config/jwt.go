package config

type JWT struct {
	AccessTokenSecret      string `env:"ACCESS_TOKEN_SECRET"`
	AccessTokenExpiration  string `env:"ACCESS_TOKEN_EXPIRATION_TIME"`
	RefreshTokenSecret     string `env:"REFRESH_TOKEN_SECRET"`
	RefreshTokenExpiration string `env:"REFRESH_TOKEN_EXPIRATION_TIME"`
}
