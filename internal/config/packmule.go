package config

type packmule struct {
	Directory string `koanf:"directory"`
	Workers   int    `koanf:"workers"`
}
