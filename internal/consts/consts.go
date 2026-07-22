package consts

import "time"

const (
	Unknown  = "unknown"
	FontName = "smslant"

	TimeFormat = time.RFC3339

	SimpleIndexURL   = "https://pypi.org/simple/"
	SimpleJSONAccept = "application/vnd.pypi.simple.v1+json"
	DefaultUserAgent = "pypiindex-collector/1.0 (+https://example.com)"
)
