package broker

type Config struct {
	Address         string `env:"BROKER_SERVER"`
	GroupID         string `env:"GROUP_ID"`
	AutoOffsetReset string `env:"AUTO_OFFSET_RESET"`
	MaxRetries      int    `env:"MAX_RETRIES" envDefault:"5"`
}
