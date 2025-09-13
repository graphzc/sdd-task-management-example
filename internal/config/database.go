package config

type Database struct {
	Driver string `env:"DRIVER"`
	URI    string `env:"URI"`
}
