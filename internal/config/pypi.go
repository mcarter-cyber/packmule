package config

type pypi struct {
	PypiURL   string `koanf:"pypi-url"`
	Directory string `koanf:"directory"`
}
