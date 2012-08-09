package db


type DBConfig struct {
	Driver string
	DBName string
	User string
	Password string
}

var configs = make(map[string]*DBConfig)

func Configure(key string, conf *DBConfig) {
	configs[key] = conf
}

func Configs() map[string]*DBConfig {
	return configs
}