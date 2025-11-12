package main

import (
	"os"

	"github.com/MeteorsLiu/llarmvp/cmd/internal"
)

func main() {
	if err := internal.Execute(); err != nil {
		os.Exit(1)
	}
}
