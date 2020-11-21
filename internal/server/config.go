package server

type Config struct {
	DBPath string `env:"SQLITE_PATH,required"`
}
