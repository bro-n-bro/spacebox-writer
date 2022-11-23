package broker

import "time"

type Config struct {
	Address         string        `yaml:"address"`
	GroupID         string        `yaml:"group_id"`
	AutoOffsetReset string        `yaml:"auto_offset_reset"`
	Topics          []string      `yaml:"topics"`
	StartTimeout    time.Duration `yaml:"start_timeout"`
}
