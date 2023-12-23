package main

import (
	"os"

	"github.com/rfbatista/gohydrate"
)

func main() {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	b := gohydrate.BuildOptions{
		EntryPoints: []string{path + "/app.jsx"},
		Outdir:      path,
	}
	err = gohydrate.Build(b)
	if err != nil {
		os.Exit(1)
	}
}
