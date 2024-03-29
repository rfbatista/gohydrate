package bundler

import (
	"fmt"
	"strings"

	esbuild "github.com/evanw/esbuild/pkg/api"
)

func ServerBundler(entryPoint string, pagesDir string, assetsDirPath string) (BuildResult, error) {
	opts := esbuild.BuildOptions{
		Stdin: &esbuild.StdinOptions{
			Contents:   entryPoint,
			Loader:     esbuild.LoaderJSX,
			ResolveDir: pagesDir,
		},
		Platform:   esbuild.PlatformNode,
		Bundle:     true,
		Write:      false,
		Outdir:     "/",
		Metafile:   false,
		AssetNames: fmt.Sprintf("%s/[name]", strings.TrimPrefix(assetsDirPath, "/")),
		Loader: map[string]esbuild.Loader{ // for loading images properly
			".png":   esbuild.LoaderFile,
			".svg":   esbuild.LoaderFile,
			".jpg":   esbuild.LoaderFile,
			".jpeg":  esbuild.LoaderFile,
			".gif":   esbuild.LoaderFile,
			".bmp":   esbuild.LoaderFile,
			".woff2": esbuild.LoaderFile,
			".woff":  esbuild.LoaderFile,
			".ttf":   esbuild.LoaderFile,
			".eot":   esbuild.LoaderFile,
		},
	}
	result := esbuild.Build(opts)
	if len(result.Errors) > 0 {
		fileLocation := "unknown"
		lineNum := "unknown"
		if result.Errors[0].Location != nil {
			fileLocation = result.Errors[0].Location.File
			lineNum = result.Errors[0].Location.LineText
		}
		return BuildResult{}, fmt.Errorf("%s <br>in %s <br>at %s", result.Errors[0].Text, fileLocation, lineNum)
	}

	var br BuildResult
	for _, file := range result.OutputFiles {
		if strings.HasSuffix(file.Path, "stdin.js") {
			br.JS = string(file.Contents)
		} else if strings.HasSuffix(file.Path, "stdin.css") {
			br.CSS = string(file.Contents)
		}
	}
	return br, nil
}
