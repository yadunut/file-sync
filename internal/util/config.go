package util

type Config struct {
	Port string
	Host string
}

func (c Config) GetUrl() string {
	return c.Host + c.Port
}

// GetConfig should be building the config from the default file / flags, but I can figure that out later.
func GetConfig() Config {
	return Config{
		Port: ":8080",
		Host: "localhost",
	}
}
