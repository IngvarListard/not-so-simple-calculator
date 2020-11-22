package calcapp

type Config struct {
	DBPath     string  `env:"SQLITE_PATH,required"`
	LogsToFile *string `env:"LOGS_TO_FILE"`
	LogLevel   string  `env:"LOG_LEVEL" envDefault:"debug"`
	LogsToJSON bool    `env:"LOGS_TO_JSON" envDefault:"false"`
	AccessLog  bool    `env:"ACCESS_LOG" envDefault:"true"`
}
