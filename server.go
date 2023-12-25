package gohydrate

import (
	"encoding/json"
	"path/filepath"

	"github.com/rfbatista/gohydrate/internal/bundler"
	"github.com/rfbatista/gohydrate/internal/entrypoint"
	"github.com/rfbatista/gohydrate/internal/renderer"
)

type serverBuild struct {
	HTML  string
	CSS   string
	Error error
}

func (e *Engine) buildServer(c PageConfig, props json.RawMessage, fp string, res *chan serverBuild) {
	b, valid := e.cacheManager.GetServerBuild(c.Filename)
	if !valid {
		r, err := entrypoint.GenerateSeverEntryPoint(nil, fp)
		if err != nil {
			e.log.Error(err.Error())
			*res <- serverBuild{Error: err}
			return
		}
		b, err = bundler.ServerBundler(r, e.PagesFullPath(), e.AssetsPath)
		if err != nil {
			e.log.Error(err.Error())
			*res <- serverBuild{Error: err}
			return
		}
		e.cacheManager.SaveServerBuild(c.Filename, b)
	}
	jsPath := filepath.Join(e.PagesFullPath(), "build/index.js")
	b = e.AddProps(b, props)
	p, err := renderer.CreatePage(jsPath, b)
	if err != nil {
		e.log.Error(err.Error())
		*res <- serverBuild{Error: err}
		return
	}
	*res <- serverBuild{HTML: p.HTML, CSS: p.CSS, Error: nil}
}
