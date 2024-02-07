package config

type (
	Config struct {
		Apikey string
		Secret string
	}
)

func NewConfig(apiKey, secret string) Config {
	return Config{
		Apikey: apiKey,
		Secret: secret,
	}
}
