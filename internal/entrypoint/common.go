package entrypoint

type EntryPointParams struct {
	Imports      []string
	PageFilePath string
}

var entryPointTemplate = `
import * as React from "react";
{{range $import := .Imports}}{{$import}} {{end}}
import App from "{{ .FilePath }}";
{{if .Props}}const props = {{.Props}}{{end}}
{{ .RenderFunction }}`
