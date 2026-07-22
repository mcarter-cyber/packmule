package config

var DefaultConfig = Config{
	Packmule: packmule{
		Directory: "./mirror",
		Workers:   3,
	},
	Pypi: pypi{
		PypiURL:   "https://pypi.org/",
		Directory: "./pypi",
	},
}
