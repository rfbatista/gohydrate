package bundler

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
	"go.uber.org/zap"
)

func NewBundler(logger *zap.Logger) (*Bundler, error) {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	projectPath, err := filepath.EvalSymlinks(ex)
	if err != nil {
		panic(err)
	}
	f := path.Join(path.Dir(projectPath), "/dist")
	return &Bundler{AssetsPath: f, logger: logger}, nil
}

type BundlerPageInfo struct {
	Name        string
	JSPath      string
	CSSPath     string
	JSName      string
	CSSName     string
	EntryPoints []string
}

type buildResult struct {
	JS           string
	CSS          string
	JSPath       string
	CSSPath      string
	JSName       string
	CSSName      string
	Dependencies []string
}

type BundlerPageParams struct {
	Name        string
	EntryPoints []string
}

type Bundler struct {
	logger     *zap.Logger
	AssetsPath string
	pages      []BundlerPageInfo
}

func (b *Bundler) AddPage(pa BundlerPageParams) (BundlerPageInfo, error) {
	hashName := strconv.Itoa(rand.Int())
	cssName := hashName + ".css"
	jsName := hashName + ".js"
	cssPath := b.AssetsPath + "/" + cssName
	jsPath := b.AssetsPath + "/" + jsName
	p := BundlerPageInfo{
		EntryPoints: pa.EntryPoints,
		Name:        pa.Name,
		JSName:      jsName,
		CSSName:     cssName,
		CSSPath:     cssPath,
		JSPath:      jsPath,
	}
	b.pages = append(b.pages, p)
	return p, nil
}

func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *Bundler) BuildCSS(path string) (*api.BuildResult, error) {
	result := api.Build(api.BuildOptions{
		EntryPoints: []string{path},
		Bundle:      true,
		Write:       false,
		Loader: map[string]api.Loader{
			".css": api.LoaderGlobalCSS,
		},
	})
	if len(result.Errors) > 0 {
		return nil, errors.New(result.Errors[0].Text)
	}
	return &result, nil
}

func (b *Bundler) Build() error {
	_, err := os.Stat(b.AssetsPath)
	if os.IsNotExist(err) {
		err := os.Mkdir(b.AssetsPath, os.ModePerm)
		if err != nil {
			return err
		}
	} else {
		if b.AssetsPath == "" {
			return errors.New("empty assets directory path")
		}
		b.logger.Debug("cleaning assets folder")
		err := RemoveContents(b.AssetsPath)
		if err != nil {
			return err
		}
	}
	for _, page := range b.pages {
		b.logger.Debug("building page", zap.String("name", page.Name))
		_, err := serverBundler(b.AssetsPath, page)
		if err != nil {
			return err
		}
	}
	return nil
}

func serverBundler(assetsDirPath string, page BundlerPageInfo) (BuildResult, error) {
	opts := api.BuildOptions{
		Format:            api.FormatIIFE,
		GlobalName:        "bundle",
		MinifySyntax:      true,
		MinifyWhitespace:  true,
		MinifyIdentifiers: false,
		EntryPoints:       page.EntryPoints,
		Platform:          api.PlatformBrowser,
		Bundle:            true,
		Write:             false,
		Outdir:            assetsDirPath,
		Metafile:          false,
		AssetNames:        fmt.Sprintf("%s/[name]", strings.TrimPrefix(assetsDirPath, "/")),
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
	}
	result := api.Build(opts)
	if len(result.Errors) > 0 {
		fileLocation := "unknown"
		lineNum := "unknown"
		if result.Errors[0].Location != nil {
			fileLocation = result.Errors[0].Location.File
			lineNum = result.Errors[0].Location.LineText
		}
		return BuildResult{}, fmt.Errorf(
			"%s <br>in %s <br>at %s",
			result.Errors[0].Text,
			fileLocation,
			lineNum,
		)
	}

	var br BuildResult
	for _, file := range result.OutputFiles {
		if strings.HasSuffix(file.Path, ".js") {
			br.JS = string(file.Contents)
		} else if strings.HasSuffix(file.Path, ".css") {
			br.CSS = string(file.Contents)
		}
	}
	br.CSSName = page.CSSName
	err := os.WriteFile(page.CSSPath, []byte(br.CSS), 0o644)
	if err != nil {
		return BuildResult{}, err
	}
	err = os.WriteFile(page.JSPath, []byte(br.JS), 0o644)
	if err != nil {
		return BuildResult{}, err
	}
	br.JS = page.JSPath
	br.CSS = page.CSSPath
	return br, nil
}
