package db

type Config struct {
	Driver string
	Host string
	Port int
	DBName string
	User string
	Password string
}

var configs = make(map[string]map[string]*Config)

func Configs(driver string) map[string]*Config {
	return configs[driver]
}

func Named(key string, conf *Config) {
	if configs[key] == nil {
		configs[key] = make(map[string]*Config)
	}
	configs[key][conf.Driver] = conf
}

func Development(conf *Config) {
	Named("dev", conf)
}