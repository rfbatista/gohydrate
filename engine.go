package gohydrate

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/rfbatista/gohydrate/internal/bundler"
	"github.com/rfbatista/gohydrate/internal/cache"
	reactrenderer "github.com/rfbatista/gohydrate/internal/entrypoint"
	"github.com/rfbatista/gohydrate/internal/page"
	"github.com/rfbatista/logger"
)

type EngineConfig struct {
	BasePath  string
	LogTo     io.Writer
	PagesPath string
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
	log, _ := logger.New(logger.LoggerConfig{WriteTo: to, LogLevel: logger.Error})
	return &Engine{PagesPath: c.PagesPath, BasePath: c.BasePath, log: log, errorPage: []byte("failed to render page")}, nil
}

type Engine struct {
	cacheManager cache.Manager
	BasePath     string
	PagesPath    string
	log          *logger.Logger
	errorPage    []byte
}

func (e *Engine) RenderPage(c PageConfig) (*page.Page, error) {
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
	r, valid := e.cacheManager.GetServerBuild(c.Filename)
	if !valid {
		r, err = reactrenderer.GenerateSeverEntryPoint(nil, fp, props)
		if err != nil {
			e.log.Error(err.Error())
			return nil, err
		}
	}
	b, err := bundler.ServerBundler(r, e.PagesFullPath())
	if err != nil {
		e.log.Error(err.Error())
		return nil, err
	}
	jsPath := filepath.Join(e.PagesFullPath(), "build/index.js")
	p, err := page.CreatePage(jsPath, b)
	if err != nil {
		e.log.Error(err.Error())
		return nil, err
	}
	return p, nil
}

func (e *Engine) PagesFullPath() string {
	return filepath.Join(e.BasePath, e.PagesPath)
}
