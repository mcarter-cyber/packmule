package config

import "github.com/knadh/koanf/v2"

var K = koanf.New(".")

type Config struct {
	Packmule packmule `koanf:"packmule"`
	Pypi     pypi     `koanf:"pypi"`
}
