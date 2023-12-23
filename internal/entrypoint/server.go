package entrypoint

import (
	"encoding/json"
	"strings"
	"text/template"
)

var renderWithProps = `process.stdout.write(renderToString(<App {...props} />))`
var renderWithoutProps = `process.stdout.write(renderToString(<App />))`

func GenerateSeverEntryPoint(imports []string, filePath string, props json.RawMessage) (string, error) {
	imports = append(imports, `import { renderToString } from "react-dom/server";`)
	params := map[string]interface{}{
		"Imports":            imports,
		"FilePath":           filePath,
		"RenderFunction":     renderWithProps,
		"SuppressConsoleLog": true,
	}
	if props != nil {
		propsJson, _ := json.Marshal(&props)
		params["Props"] = string(propsJson)
	}
	templ, err := template.New("buildTemplate").Parse(entryPointTemplate)
	if err != nil {
		return "", err
	}
	var out strings.Builder
	err = templ.Execute(&out, params)
	if err != nil {
		return "", err
	}
	return out.String(), nil
}
