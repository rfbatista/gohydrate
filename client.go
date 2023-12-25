package gohydrate

import (
	"encoding/json"
	"path/filepath"

	"github.com/rfbatista/gohydrate/internal/bundler"
	"github.com/rfbatista/gohydrate/internal/entrypoint"
	"github.com/rfbatista/gohydrate/internal/renderer"
)

type clientBuild struct {
	JS    string
	Error error
}

func (e *Engine) buildClient(c PageConfig, props json.RawMessage, fp string, res *chan clientBuild) {
	b, valid := e.cacheManager.GetServerBuild(c.Filename)
	if !valid {
		r, err := entrypoint.GenerateSeverEntryPoint(nil, fp)
		if err != nil {
			e.log.Error(err.Error())
			*res <- clientBuild{Error: err}
			return
		}
		b, err = bundler.ClientBundler(r, e.PagesFullPath(), e.AssetsPath, e.isProd)
		if err != nil {
			e.log.Error(err.Error())
			*res <- clientBuild{Error: err}
			return
		}
		e.cacheManager.SaveServerBuild(c.Filename, b)
	}
	jsPath := filepath.Join(e.PagesFullPath(), "build/index.js")
	b = e.AddProps(b, props)
	p, err := renderer.CreatePage(jsPath, b)
	if err != nil {
		e.log.Error(err.Error())
		*res <- clientBuild{Error: err}
		return
	}
  *res <- clientBuild{JS: p.JS, Error: nil}
}
