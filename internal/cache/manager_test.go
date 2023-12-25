package cache

import "testing"

func TestSaveServerBuild(t *testing.T) {
	manager := Manager{serverBuilds: make(map[string]string)}

	page := "home"
	build := "v1.0"

	manager.SaveServerBuild(page, build)

	if manager.serverBuilds[page] != build {
		t.Errorf("SaveServerBuild failed, expected: %s, got: %s", build, manager.serverBuilds[page])
	}
}

func TestInvalidateServerBuild(t *testing.T) {
	manager := Manager{serverBuilds: make(map[string]string)}

	page := "home"
	build := "v1.0"

	manager.SaveServerBuild(page, build)
	manager.InvalidateServerBuild(page, build)

	if manager.serverBuilds[page] != "" {
		t.Errorf("InvalidateServerBuild failed, expected: '', got: %s", manager.serverBuilds[page])
	}
}

func TestGetServerBuild(t *testing.T) {
	manager := Manager{serverBuilds: make(map[string]string)}

	page := "home"
	build := "v1.0"

	manager.SaveServerBuild(page, build)
	result, exists := manager.GetServerBuild(page)
	if !exists || result != build {
		t.Errorf("GetServerBuild failed for existing build, expected: %s, got: %s", build, result)
	}

	result, exists = manager.GetServerBuild("nonexistent")
	if exists || result != "" {
		t.Errorf("GetServerBuild failed for non-existing build, expected: '', got: %s", result)
	}
}
