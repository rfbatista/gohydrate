package renderer

import (
	"os"
	"os/exec"
	"strings"

	"github.com/rfbatista/gohydrate/internal/bundler"
)

type Page struct {
	CSS  string
	JS   string
	HTML string
}

func CreatePage(jsPath string, b bundler.BuildResult) (*Page, error) {
	err := os.WriteFile(jsPath, []byte(b.JS), 0644)
	if err != nil {
		return nil, err
	}
	cmd := exec.Command("node", jsPath)
	stdOut := new(strings.Builder)
	stdErr := new(strings.Builder)
	cmd.Stdout = stdOut
	cmd.Stderr = stdErr
  err = cmd.Run()
	if err != nil {
		return nil, err
	}
	htmlFile := stdOut.String()
	return &Page{
		HTML: htmlFile,
		JS:   b.JS,
		CSS:  b.CSS,
	}, nil
}
