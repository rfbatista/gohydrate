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

type ClientBuild struct {
	JS    string
	Error error
}

type HydrateClientConfig struct {
	Log           logger.Logger
	CacheManager  cache.Manager
	PagesFullPath string
	AssetsPath    string
	Filename      string
	Fullpath      string
	Result        *chan ClientBuild
	Props         string
	IsProd        bool
}

func HydrateClient(c HydrateClientConfig) {
	c.Log.Info(fmt.Sprintf("create client hydration for page %s", c.Filename))
	b, valid := c.CacheManager.GetServerBuild(c.Filename)
	if !valid {
		entrypointCode, err := entrypoint.GenerateSeverEntryPoint(nil, c.Fullpath)
		if err != nil {
			c.Log.Error(fmt.Sprintf("failed to generate server entry point %s", err))
			*c.Result <- ClientBuild{Error: err}
			return
		}
		b, err = bundler.ClientBundler(bundler.ClientBundlerConfig{
			EntryPoint: entrypointCode,
			PagesDir:   c.PagesFullPath,
			AssetsPath: c.AssetsPath,
			IsProd:     c.IsProd,
		})
		if err != nil {
			c.Log.Error(fmt.Sprintf("failed to create client bundle \n %s", err))
			*c.Result <- ClientBuild{Error: err}
			return
		}
		c.CacheManager.SaveServerBuild(c.Filename, b)
	}
	jsPath := filepath.Join(c.PagesFullPath, "build/index.js")
	if c.Props != "" {
		b.JS = fmt.Sprintf(`var props = %s; %s`, c.Props, b.JS)
		c.Log.Info("adding props")
	}
	c.Log.Debug("rendering page from JS asset")
	p, err := renderer.CreatePage(jsPath, b)
	if err != nil {
		c.Log.Error(fmt.Sprintf("failed to render client page from bundler \n %s", err))
		*c.Result <- ClientBuild{Error: err}
		return
	}
	*c.Result <- ClientBuild{JS: p.JS, Error: nil}
}
