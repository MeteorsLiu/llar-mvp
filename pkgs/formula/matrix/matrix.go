package matrix

import "runtime"

type Matrix struct {
	Require map[string][]string
	Options map[string][]string
}

func Current() Matrix {
	return Matrix{
		Require: map[string][]string{
			"os":   {runtime.GOOS},
			"arch": {runtime.GOARCH},
		},
		Options: map[string][]string{},
	}
}
