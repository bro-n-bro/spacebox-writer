package graphql

import "time"

type Config struct {
	Address      string        `yaml:"address"`
	StartTimeout time.Duration `yaml:"start_timeout"`
}
