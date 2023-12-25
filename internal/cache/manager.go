package cache

import (
	"github.com/rfbatista/gohydrate/internal/bundler"
)

func New() Manager {
	return Manager{
		serverBuilds: make(map[string]bundler.BuildResult),
		clientBuilds: make(map[string]bundler.BuildResult),
	}
}

type Manager struct {
	serverBuilds map[string]bundler.BuildResult
	clientBuilds map[string]bundler.BuildResult
}

func (c *Manager) SaveServerBuild(page string, build bundler.BuildResult) {
	c.serverBuilds[page] = build
}

func (c *Manager) InvalidateServerBuild(page string, build string) {
	delete(c.serverBuilds, page)
}

func (c *Manager) GetServerBuild(page string) (bundler.BuildResult, bool) {
	build, ok := c.clientBuilds[page]
	return build, ok
}

func (c *Manager) SaveClientBuild(page string, build bundler.BuildResult) {
	c.clientBuilds[page] = build
}

func (c *Manager) InvalidateClientBuild(page string, build string) {
	delete(c.clientBuilds, page)
}

func (c *Manager) GetClientBuild(page string) (bundler.BuildResult, bool) {
	build, ok := c.clientBuilds[page]
	return build, ok
}
