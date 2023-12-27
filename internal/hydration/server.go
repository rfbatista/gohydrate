package hydration

import (
	"fmt"
	"path/filepath"

	"github.com/rfbatista/gohydrate/internal/bundler"
	"github.com/rfbatista/gohydrate/internal/cache"
	"github.com/rfbatista/gohydrate/internal/entrypoint"
	"github.com/rfbatista/gohydrate/internal/renderer"
	"github.com/rfbatista/logger"
)

type ServerBuild struct {
	HTML  string
	CSS   string
	Error error
}

type RenderServerHTMLConfig struct {
	Log           logger.Logger
	CacheManager  cache.Manager
	PagesFullPath string
	AssetsPath    string
	Filename      string
	FileFullPath            string
	Result           *chan ServerBuild
	Props         string
}

func RenderServerHTML(
	c RenderServerHTMLConfig,
) {
	c.Log.Debug(fmt.Sprintf("rendering structure for page %s", c.Filename))
	b, valid := c.CacheManager.GetServerBuild(c.Filename)
	if !valid {
		r, err := entrypoint.GenerateSeverEntryPoint(nil, c.FileFullPath)
		if err != nil {
			c.Log.Error(err.Error())
			*c.Result <- ServerBuild{Error: err}
			return
		}
		b, err = bundler.ServerBundler(r, c.PagesFullPath, c.AssetsPath)
		if err != nil {
			c.Log.Error(err.Error())
			*c.Result <- ServerBuild{Error: err}
			return
		}
		c.CacheManager.SaveServerBuild(c.Filename, b)
	}
	jsPath := filepath.Join(c.PagesFullPath, "build/index.js")
	if c.Props == "" {
		c.Props = "{}"
	}
	b.JS = fmt.Sprintf(`var props = %s; %s`, c.Props, b.JS)
	p, err := renderer.CreatePage(jsPath, b)
	if err != nil {
		c.Log.Error(fmt.Sprintf("failed to render server page from bundler \n %s", err))
		*c.Result <- ServerBuild{Error: err}
		return
	}
	*c.Result <- ServerBuild{HTML: p.HTML, CSS: p.CSS, Error: nil}
}
