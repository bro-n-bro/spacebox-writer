package modules

// Config is a configuration for modules.
type Config struct {
	Modules []string `env:"MODULES"`
}
