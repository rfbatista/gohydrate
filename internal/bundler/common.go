package bundler

type BuildResult struct {
	JS           string
	CSS          string
	Dependencies []string
}
