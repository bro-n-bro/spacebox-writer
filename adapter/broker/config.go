package broker

import "time"

type Config struct {
	Address         string        `yaml:"address"`
	GroupID         string        `yaml:"group_id"`
	AutoOffsetReset string        `yaml:"auto_offset_reset"`
	StartTimeout    time.Duration `yaml:"start_timeout"`
}
