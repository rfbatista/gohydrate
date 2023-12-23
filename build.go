package gohydrate

import (
	"github.com/evanw/esbuild/pkg/api"
)

type BuildOptions struct {
	PagePath   string
	EntryPoints []string
	Outdir      string
}

func BuildPages(o BuildOptions) error {
	result := api.Build(api.BuildOptions{
		EntryPoints: o.EntryPoints,
		Outdir:      o.Outdir,
		Bundle:      true,
		Write:       true,
		LogLevel:    api.LogLevelInfo,
		Loader: map[string]api.Loader{ // for loading images properly
			".png":   api.LoaderFile,
			".svg":   api.LoaderFile,
			".jpg":   api.LoaderFile,
			".jpeg":  api.LoaderFile,
			".gif":   api.LoaderFile,
			".bmp":   api.LoaderFile,
			".woff2": api.LoaderFile,
			".woff":  api.LoaderFile,
			".ttf":   api.LoaderFile,
			".eot":   api.LoaderFile,
		},
	})
	if len(result.Errors) > 0 {
		return buildFailed
	}
	return nil
}

func Build(o BuildOptions) error {
	result := api.Build(api.BuildOptions{
		EntryPoints: o.EntryPoints,
		Outdir:      o.Outdir,
		Bundle:      true,
		Write:       true,
		LogLevel:    api.LogLevelInfo,
		Loader: map[string]api.Loader{ // for loading images properly
			".png":   api.LoaderFile,
			".svg":   api.LoaderFile,
			".jpg":   api.LoaderFile,
			".jpeg":  api.LoaderFile,
			".gif":   api.LoaderFile,
			".bmp":   api.LoaderFile,
			".woff2": api.LoaderFile,
			".woff":  api.LoaderFile,
			".ttf":   api.LoaderFile,
			".eot":   api.LoaderFile,
		},
	})
	if len(result.Errors) > 0 {
		return buildFailed
	}
	return nil
}
