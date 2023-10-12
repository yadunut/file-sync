package util

type Config struct {
	Port string
	Host string
}

func (c Config) GetUrl() string {
	return c.Host + c.Port
}

func GetConfig() Config {
	return Config{
		Port: ":8080",
		Host: "localhost",
	}
}
