package config

type Config struct {
	Port      string
	MaxMemory int64
}

func Default() Config {
	return Config{
		Port:      ":6379",
		MaxMemory: 0,
	}
}
