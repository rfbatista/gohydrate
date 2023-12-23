package gohydrate

import (
	"bytes"
	"html/template"

	"github.com/rfbatista/gohydrate/internal/page"
	"github.com/rfbatista/gohydrate/internal/templates"
)

type Params struct {
	Title      string
	MetaTags   map[string]string
	OGMetaTags map[string]string
	Links      []struct {
		Href     string
		Rel      string
		Media    string
		Hreflang string
		Type     string
		Title    string
	}
	Page    page.Page
	RouteID string
}

func RenderPageToHTML(params Params) ([]byte, error) {
	t := template.Must(template.New("").Parse(templates.BaseTemplate))
	var output bytes.Buffer
	err := t.Execute(&output, params)
	if err != nil {
		return nil, err
	}
	return output.Bytes(), nil
}
