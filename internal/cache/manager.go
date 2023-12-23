package cache

type Manager struct {
	serverBuilds map[string]string
}

func (c *Manager) SaveServerBuild(page string, build string) {
	c.serverBuilds[page] = build
}

func (c *Manager) InvalidateServerBuild(page string, build string) {
	c.serverBuilds[page] = ""
}

func (c *Manager) GetServerBuild(page string) (string, bool) {
	build := c.serverBuilds[page]
	if build == "" {
		return "", false
	}
  return build, true
}
