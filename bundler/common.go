package bundler

type BuildResult struct {
	JS           string
	CSSName      string
	CSS          string
	Dependencies []string
}
