package gohydrate

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/rfbatista/gohydrate/internal/bundler"
	"github.com/rfbatista/gohydrate/internal/cache"
	"github.com/rfbatista/gohydrate/internal/renderer"
	"github.com/rfbatista/logger"
)

type EngineConfig struct {
	BasePath   string
	LogTo      io.Writer
	PagesPath  string
	AssetsPath string
	IsProd     bool
}

type PageConfig struct {
	Filename string
	Props    interface{}
}

func New(c EngineConfig) (*Engine, error) {
	var to io.Writer
	if c.LogTo != nil {
		to = c.LogTo
	} else {
		to = os.Stdout
	}
	mng := cache.New()
	log, _ := logger.New(logger.LoggerConfig{WriteTo: to, LogLevel: logger.Error})
	return &Engine{isProd: c.IsProd, PagesPath: c.PagesPath, cacheManager: mng, BasePath: c.BasePath, log: log, errorPage: []byte("failed to render page"), AssetsPath: c.AssetsPath}, nil
}

type Engine struct {
	isProd       bool
	cacheManager cache.Manager
	BasePath     string
	PagesPath    string
	log          *logger.Logger
	errorPage    []byte
	AssetsPath   string
}

func (e *Engine) RenderPage(c PageConfig) (*renderer.Page, error) {
	fp := filepath.Join(e.PagesFullPath(), c.Filename)
	var props json.RawMessage
	var err error
	if c.Props != nil {
		props, err = json.Marshal(c.Props)
		if err != nil {
			e.log.Error(err.Error())
			return nil, err
		}

	}
	severBuildChan := make(chan serverBuild)
	clientBuildChan := make(chan clientBuild)
	go e.buildServer(c, props, fp, &severBuildChan)
	go e.buildClient(c, props, fp, &clientBuildChan)
	serverBuild := <-severBuildChan
	clientBuild := <-clientBuildChan
	return &renderer.Page{HTML: serverBuild.HTML, CSS: serverBuild.CSS, JS: clientBuild.JS}, nil
}

func (e *Engine) PagesFullPath() string {
	return filepath.Join(e.BasePath, e.PagesPath)
}

func (e *Engine) AddProps(b bundler.BuildResult, props interface{}) bundler.BuildResult {
	if props != nil {
		propsJson, _ := json.Marshal(&props)
		b.JS = fmt.Sprintf(`var props = %s; %s`, string(propsJson), b.JS)
	}
	return b

}
