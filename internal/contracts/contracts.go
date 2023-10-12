package contracts

type Version struct {
	Version string `json:"version"`
}

type Config struct {
	Port string
	Host string
}

func (c Config) GetUrl() string {
	return c.Host + c.Port
}
