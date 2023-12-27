package gohydrate

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/rfbatista/gohydrate/internal/cache"
	"github.com/rfbatista/gohydrate/internal/hydration"
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
	log, _ := logger.New(logger.LoggerConfig{WriteTo: to, LogLevel: logger.Debug, WithDateTime: true})
	e := &Engine{isProd: c.IsProd, PagesPath: c.PagesPath, cacheManager: mng, BasePath: c.BasePath, log: log, errorPage: []byte("failed to render page"), AssetsPath: c.AssetsPath}
	log.Info(fmt.Sprintf("loading pages from %s", e.PagesFullPath()))
	return e, nil
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
	fullPath := filepath.Join(e.PagesFullPath(), c.Filename)
	var props string
	if c.Props != nil {
		rawProps, err := json.Marshal(c.Props)
		if err != nil {
			e.log.Error(fmt.Sprintf("failed to marshal props for page %s : %s", c.Filename, err))
			return nil, err
		}
		props = string(rawProps)
	}
	severBuildChan := make(chan hydration.ServerBuild)
	clientBuildChan := make(chan hydration.ClientBuild)
	go hydration.HydrateClient(
		hydration.HydrateClientConfig{
			Log:           *e.log,
			CacheManager:  e.cacheManager,
			PagesFullPath: e.PagesFullPath(),
			AssetsPath:    e.AssetsPath,
			Filename:      c.Filename,
			Fullpath:      fullPath,
			Result:        &clientBuildChan,
			Props:         props,
			IsProd:        e.isProd,
		},
	)
	go hydration.RenderServerHTML(
		hydration.RenderServerHTMLConfig{
			Log:           *e.log,
			CacheManager:  e.cacheManager,
			PagesFullPath: e.PagesFullPath(),
			AssetsPath:    e.AssetsPath,
			Filename:      c.Filename,
			FileFullPath:  fullPath,
			Result:        &severBuildChan,
			Props:         props,
		},
	)
	serverBuild := <-severBuildChan
	clientBuild := <-clientBuildChan
	if serverBuild.Error != nil {
		e.log.Error(fmt.Sprintf("failed to build server \n %s", serverBuild.Error))
		return nil, serverBuild.Error
	}
	if clientBuild.Error != nil {
		e.log.Error(fmt.Sprintf("failed to build client \n %s", clientBuild.Error))
		return nil, clientBuild.Error
	}
	return &renderer.Page{HTML: serverBuild.HTML, CSS: serverBuild.CSS, JS: clientBuild.JS}, nil
}

func (e *Engine) PagesFullPath() string {
	return filepath.Join(e.BasePath, e.PagesPath)
}
