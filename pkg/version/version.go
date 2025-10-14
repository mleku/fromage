package version

import (
	_ "embed"
)

//go:embed version
var V string

var URL = "https://github.com/mleku/fromage"
